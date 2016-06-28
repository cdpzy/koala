package codec

import (
    "fmt"
	"encoding/base64"
    "github.com/doublemo/koala/msic"
)

var OutPacketBufferMaxSize uint = 60000 // default

type H264VideoRTPSink struct {
    sps           string
	pps           string
	spsSize       int
	ppsSize       int
    nextSendTime                    msic.Timeval
	noFramesLeft                    bool
	isFirstPacket                   bool
	currentTimestamp                uint
	ourMaxPacketSize                uint
	timestampPosition               uint
	specialHeaderSize               uint
	numFramesUsedSoFar              uint
	specialHeaderPosition           uint
	curFragmentationOffset          uint
	totalFrameSpecificHeaderSizes   uint
	previousFrameEndedFragmentation bool
	onSendErrorFunc                 interface{}
    ssrc                       uint
	seqNo                      uint
	octetCount                 uint
	packetCount                uint // incl RTP hdr
	timestampBase              uint
	totalOctetCount            uint
	rtpPayloadType             uint
	rtpTimestampFrequency      uint
	timestampFrequency         uint
	rtpPayloadFormatName       string
	enableRTCPReports          bool
	nextTimestampHasBeenPreset bool
    rtpInterface               *RTPInterface
    outBuf                     *OutPacketBuffer
    ourFragmenter              *H264FUAFragmenter
	source                     *H264VideoStreamParser
	source2                    *H264FUAFragmenter
	//transmissionStatsDB        *RTPTransmissionStatsDB
    //source  IFramedSource
	//rtpSink IRTPSink
}

func (h264VideoRTPSink *H264VideoRTPSink) RtpmapLine() string {
    var rtpmapLine string
	if h264VideoRTPSink.rtpPayloadType >= 96 {
		encodingParamsPart := ""
		rtpmapFmt := "a=rtpmap:%d %s/%d%s\r\n"
		rtpmapLine = fmt.Sprintf(rtpmapFmt,
			h264VideoRTPSink.RtpPayloadType(),
			h264VideoRTPSink.RtpPayloadFormatName(),
			h264VideoRTPSink.RtpTimestampFrequency(), encodingParamsPart)
	}

	return rtpmapLine
}

func (h264VideoRTPSink *H264VideoRTPSink) RtpPayloadType() uint {
    return h264VideoRTPSink.rtpPayloadType
}

func (h264VideoRTPSink *H264VideoRTPSink) RtpPayloadFormatName() string {
    return h264VideoRTPSink.rtpPayloadFormatName
}

func (h264VideoRTPSink *H264VideoRTPSink) RtpTimestampFrequency() uint {
    return h264VideoRTPSink.rtpTimestampFrequency
}

func (h264VideoRTPSink *H264VideoRTPSink) SdpMediaType() string {
    return "video"
}

func (h264VideoRTPSink *H264VideoRTPSink) AuxSDPLine() string {
    sps := h264VideoRTPSink.sps
	pps := h264VideoRTPSink.pps
	spsSize := h264VideoRTPSink.spsSize

	if sps == "" || pps == "" {
		if h264VideoRTPSink.ourFragmenter == nil {
			return ""
		}

		framerSource := h264VideoRTPSink.ourFragmenter.InputSource()
		if framerSource == nil {
			return ""
		}

		//framerSource.getSPSandPPS()
	}

	spsBase64 := base64.NewEncoding(sps)
	ppsBase64 := base64.NewEncoding(pps)

	var profileLevelId uint8
	if spsSize >= 4 {
		profileLevelId = (sps[1] << 16) | (sps[2] << 8) | sps[3] // profile_idc|constraint_setN_flag|level_idc
	}

	fmtpFmt := "a=fmtp:%d packetization-mode=1;profile-level-id=%06X;sprop-parameter-sets=%s,%s\r\n"
	return fmt.Sprintf(fmtpFmt, h264VideoRTPSink.RtpPayloadType(), profileLevelId, spsBase64, ppsBase64)
}

func (h264VideoRTPSink *H264VideoRTPSink) continuePlaying() {
	fmt.Println(fmt.Sprintf("H264VideoRTPSink::continuePlaying -> %p", h264VideoRTPSink.source))
	if h264VideoRTPSink.ourFragmenter == nil {
		h264VideoRTPSink.ourFragmenter = NewH264FUAFragmenter(h264VideoRTPSink.source, OutPacketBufferMaxSize)
	} else {
		h264VideoRTPSink.ourFragmenter.reAssignInputSource(h264VideoRTPSink.source)
	}

	h264VideoRTPSink.source2 = h264VideoRTPSink.ourFragmenter
	h264VideoRTPSink.multiFramedPlaying()
}


func (h264VideoRTPSink *H264VideoRTPSink) multiFramedPlaying() {
	fmt.Println("MultiFramedRTPSink::continuePlaying")
	h264VideoRTPSink.buildAndSendPacket(true)
}

func (h264VideoRTPSink *H264VideoRTPSink) buildAndSendPacket(isFirstPacket bool) {
	h264VideoRTPSink.isFirstPacket = isFirstPacket
	var rtpHdr uint = 0x80000000

	rtpHdr |= h264VideoRTPSink.rtpPayloadType << 16
	rtpHdr |= h264VideoRTPSink.seqNo
	h264VideoRTPSink.outBuf.enqueueWord(rtpHdr)

	h264VideoRTPSink.timestampPosition = h264VideoRTPSink.outBuf.curPacketSize()

	h264VideoRTPSink.outBuf.skipBytes(4)

	h264VideoRTPSink.outBuf.enqueueWord(h264VideoRTPSink.ssrc)

	h264VideoRTPSink.specialHeaderPosition = h264VideoRTPSink.outBuf.curPacketSize()
	h264VideoRTPSink.specialHeaderSize     = h264VideoRTPSink.SpecialHeaderSize()
	h264VideoRTPSink.outBuf.skipBytes(h264VideoRTPSink.specialHeaderSize)

	// Begin packing as many (complete) frames into the packet as we can:
	h264VideoRTPSink.noFramesLeft = false
	h264VideoRTPSink.numFramesUsedSoFar = 0
	h264VideoRTPSink.totalFrameSpecificHeaderSizes = 0

	h264VideoRTPSink.packFrame()
}

func (h264VideoRTPSink *H264VideoRTPSink) packFrame() {
	if h264VideoRTPSink.outBuf.haveOverflowData() {
		// Use this frame before reading a new one from the source
		frameSize := h264VideoRTPSink.outBuf.OverflowDataSize()
		presentationTime := h264VideoRTPSink.outBuf.OverflowPresentationTime()
		durationInMicroseconds := h264VideoRTPSink.outBuf.OverflowDurationInMicroseconds()
		h264VideoRTPSink.outBuf.useOverflowData()
		h264VideoRTPSink.afterGettingFrame(frameSize, durationInMicroseconds, presentationTime)
	} else {
		// Normal case: we need to read a new frame from the source
		if h264VideoRTPSink.source2 == nil {
			return
		}
		fmt.Println("packFrame", h264VideoRTPSink.afterGettingFrame)
		//h264VideoRTPSink.source.getNextFrame(h264VideoRTPSink.outBuf.curPtr(), h264VideoRTPSink.outBuf.totalBytesAvailable(), h264VideoRTPSink.afterGettingFrame, h264VideoRTPSink.ourHandlerClosure)
	}
}

func (h264VideoRTPSink *H264VideoRTPSink) afterGettingFrame( frameSize, durationInMicroseconds uint, presentationTime msic.Timeval ) {
	fmt.Println("MultiFramedRTPSink::afterGettingFrame")
}

func (h264VideoRTPSink *H264VideoRTPSink) SpecialHeaderSize( ) uint {
	return 0
}

func NewH264VideoRTPSink(rtpGroupSock *msic.GroupSock, rtpPayloadType uint) *H264VideoRTPSink {
    h264VideoRTPSink := new(H264VideoRTPSink)
    h264VideoRTPSink.rtpPayloadType = rtpPayloadType
    h264VideoRTPSink.rtpTimestampFrequency = 90000
    h264VideoRTPSink.rtpPayloadFormatName  = "H264"
    h264VideoRTPSink.rtpInterface          = NewRTPInterface(h264VideoRTPSink, rtpGroupSock)
    // Default max packet size (1500, minus allowance for IP, UDP, UMTP headers)
	// (Also, make it a multiple of 4 bytes, just in case that matters.)
    h264VideoRTPSink.outBuf                = NewOutPacketBuffer(1000, 1448)


    return h264VideoRTPSink
}
