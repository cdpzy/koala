package media

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/doublemo/koala/media/h264"
)

// H264MediaSubSession H264
type H264MediaSubSession struct {
	BaseMediaSubSession
	SDPLines string
	pps      []byte
	sps      []byte
	sei      []byte
}

// NewH264MediaSubSession H264
func NewH264MediaSubSession(FileName string) *H264MediaSubSession {
	sess := new(H264MediaSubSession)
	sess.FileName = FileName
	sess.ID = 1
	sess.EstBitrate = 500
	sess.numChannels = 0
	sess.rtpPayloadFormatName = "H264"
	sess.rtpTimestampFrequency = 90000

	return sess
}

// AbsoluteTimeRange 计算绝对时间
func (h264MediaSubSession *H264MediaSubSession) AbsoluteTimeRange() (string, string) {
	return "", ""
}

// Duration 获取媒体时长
func (h264MediaSubSession *H264MediaSubSession) Duration() float64 {
	return 0.0
}

// GenerateSDP 创建媒体相关SDP信息
func (h264MediaSubSession *H264MediaSubSession) GenerateSDP() string {
	if h264MediaSubSession.SDPLines != "" {
		return h264MediaSubSession.SDPLines
	}
	h264MediaSubSession.rtpPayload = 96 + h264MediaSubSession.ID - 1

	rtpmapLine := h264MediaSubSession.GenerateRTPMapLine()
	rtcpmuxLine := ""
	if h264MediaSubSession.MultiplexRTCPWithRTP {
		rtcpmuxLine = "a=rtcp-mux\r\n"
	}

	s, e := h264MediaSubSession.AbsoluteTimeRange()
	rangeLine := h264MediaSubSession.GenerateRangeLine(s, e, h264MediaSubSession.Duration())
	auxSDPLine := h264MediaSubSession.GenerateAuxSDPLine()
	sdpFmt := "m=%s %d RTP/AVP %d\r\n" +
		"c=IN IP4 %s\r\n" +
		"b=AS:%d\r\n" +
		"%s%s%s%sa=control:%s\r\n"

	sdp := fmt.Sprintf(sdpFmt,
		"video",
		h264MediaSubSession.GetPort(),
		h264MediaSubSession.rtpPayload,
		"192.168.18.150",
		h264MediaSubSession.EstBitrate,
		rtpmapLine,  // a=rtpmap:... (if present)
		rtcpmuxLine, // a=rtcp-mux:... (if present)
		rangeLine,   // a=range:... (if present)
		auxSDPLine,  // optional extra SDP line
		h264MediaSubSession.TrackId(),
	)

	h264MediaSubSession.SDPLines = sdp
	return sdp
}

// GenerateAuxSDPLine aux
func (h264MediaSubSession *H264MediaSubSession) GenerateAuxSDPLine() string {
	err := h264MediaSubSession.readMediaParameters()
	if err != nil {
		return ""
	}

	if len(h264MediaSubSession.sps) < 1 {
		return ""
	}

	spsObject := h264.NewSequenceParameterSetRBSP()
	spsObject.ParseBytes(h264MediaSubSession.sps)

	pps := base64.StdEncoding.EncodeToString(h264MediaSubSession.pps)
	sps := base64.StdEncoding.EncodeToString(h264MediaSubSession.sps)
	return fmt.Sprintf("a=fmtp:%d packetization-mode=1;profile-level-id=%06X;sprop-parameter-sets=%s,%s\r\n", h264MediaSubSession.rtpPayload, spsObject.ProfileIdc, sps, pps)
}

//
func (h264MediaSubSession *H264MediaSubSession) readMediaParameters() error {
	fid, err := os.Open(h264MediaSubSession.FileName)
	if err != nil {
		return err
	}

	defer fid.Close()

	fileInfo, err := fid.Stat()
	if err != nil {
		return err
	}

	h264MediaSubSession.FileSize = fileInfo.Size()
	nal := h264.NewNalUnit()

	var (
		bytes      []byte
		preNalUint *h264.NalUnit
	)

	for {
		data := make([]byte, 20000)
		_, err := fid.Read(data)
		if err != nil {
			break
		}

		idx := 0
		oldbytes := bytes
		bytes = make([]byte, len(oldbytes)+20000)
		copy(bytes[0:len(oldbytes)], oldbytes[0:])
		copy(bytes[len(oldbytes):], data[0:])
		for idx < len(bytes) {
			if preNalUint != nil {
				_, byteI := nal.IndexByte(bytes, idx)
				if byteI == -1 {
					preNalUint.ParameterBytes = append(preNalUint.ParameterBytes, bytes...)
					break
				}

				if byteI != 0 {
					preNalUint.ParameterBytes = append(preNalUint.ParameterBytes, bytes[idx:byteI]...)
				}

				switch preNalUint.NalUnitType {
				case h264.NALU_TYPE_PPS:
					h264MediaSubSession.pps = preNalUint.ParameterBytes
				case h264.NALU_TYPE_SPS:
					h264MediaSubSession.sps = preNalUint.ParameterBytes
				case h264.NALU_TYPE_SEI:
					h264MediaSubSession.sei = preNalUint.ParameterBytes
				}
				preNalUint = nil
				idx = byteI
			}

			n, err := nal.ParseBytes(bytes[idx:])
			if err != nil {
				break
			}

			idx += n
			switch nal.NalUnitType {
			case h264.NALU_TYPE_PPS:
				h264MediaSubSession.pps = nal.ParameterBytes
			case h264.NALU_TYPE_SPS:
				h264MediaSubSession.sps = nal.ParameterBytes
			case h264.NALU_TYPE_SEI:
				h264MediaSubSession.sei = nal.ParameterBytes
			}
		}
		bytes = bytes[idx:]
		preNalUint = nal
	}
	return nil
}
