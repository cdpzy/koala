package media

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/doublemo/koala/media/h264"
)

// H264MediaSubSession H264
type H264MediaSubSession struct {
	BaseMediaSubSession
	SDPLines   string
	EstBitrate int
	pps        []byte
	sps        []byte
	sei        []byte
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
		"192.168.18.152",
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
	fmt.Println(h264MediaSubSession.readMediaParameters())
	return ""
}

//
func (h264MediaSubSession *H264MediaSubSession) readMediaParameters() error {
	fid, err := os.Open(h264MediaSubSession.FileName)
	if err != nil {
		return err
	}

	fileInfo, err := fid.Stat()
	if err != nil {
		return err
	}

	h264MediaSubSession.FileSize = fileInfo.Size()

	for {
		data := make([]byte, 20000)
		_, err := fid.Read(data)
		if err != nil {
			break
		}

		reader := bytes.NewReader(data)
		nal := h264.NewNalUnit()

		for {
			err = nal.ParseBytes(reader)
			if err == io.EOF {
				break
			}

			//if err.Error() == "NotFound" {
			switch nal.NalUnitType {
			case h264.NALU_TYPE_PPS:
				nal.Read(reader, &h264MediaSubSession.pps)
				fmt.Println("OKPPS")
			case h264.NALU_TYPE_SPS:
				//nal.Read(reader, &h264MediaSubSession.sps)
			case h264.NALU_TYPE_SEI:
				//nal.Read(reader, &h264MediaSubSession.sei)
			}
			//}
		}
	}

	fmt.Println("h264MediaSubSession.pps = ", h264MediaSubSession.pps)
	fmt.Println("h264MediaSubSession.sps = ", h264MediaSubSession.sps)
	fmt.Println("h264MediaSubSession.sei = ", h264MediaSubSession.sei)
	return nil
}
