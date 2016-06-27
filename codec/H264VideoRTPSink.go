package codec

import (
    "fmt"
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

}

func (h264VideoRTPSink *H264VideoRTPSink) continuePlaying() {
	fmt.Println(fmt.Sprintf("H264VideoRTPSink::continuePlaying -> %p", h264VideoRTPSink.source))
	if this.ourFragmenter == nil {
		h264VideoRTPSink.ourFragmenter = NewH264FUAFragmenter(h264VideoRTPSink.source, OutPacketBufferMaxSize)
	} else {
		h264VideoRTPSink.ourFragmenter.reAssignInputSource(h264VideoRTPSink.source)
	}

	h264VideoRTPSink.source = h264VideoRTPSink.ourFragmenter
	h264VideoRTPSink.multiFramedPlaying()
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
