package media

import (
	"fmt"
)

// MediaSubSession 针对不同的媒体格式，指定特定的session
type MediaSubSession interface {
	TrackId() string
	AbsoluteTimeRange() (string, string)
	Duration() float64
	GenerateSDP() string
	SetID(int)
	GetPort() int
	SetInitialPort(int)
	SetMultiplexRTCPWithRTP(bool)
	SetParentDuration(float64)
	GetBitrate() int
	GetClockRate() int
}

// BaseMediaSubSession 基础session
type BaseMediaSubSession struct {
	ID                    int
	FileName              string
	FileSize              int64
	MultiplexRTCPWithRTP  bool // 是否用RTCP端口为RTP
	InitialPort           int  // 起始端口
	rtpPayload            int
	numChannels           int     // 通道
	rtpTimestampFrequency int     // clock rate
	rtpPayloadFormatName  string  // ncoding name
	ParentDuration        float64 // 来到低级计算Duration
	EstBitrate            int
}

func (baseMediaSubSession *BaseMediaSubSession) TrackId() string {
	return fmt.Sprintf("track%d", baseMediaSubSession.ID)
}

func (baseMediaSubSession *BaseMediaSubSession) SetID(id int) {
	baseMediaSubSession.ID = id
}

func (baseMediaSubSession *BaseMediaSubSession) GetPort() int {
	if baseMediaSubSession.MultiplexRTCPWithRTP {
		return baseMediaSubSession.InitialPort
	}

	return (baseMediaSubSession.InitialPort + 1) &^ 1
}

// SetInitialPort 设置起始端口
func (baseMediaSubSession *BaseMediaSubSession) SetInitialPort(p int) {
	baseMediaSubSession.InitialPort = p
}

// SetMultiplexRTCPWithRTP 设置端口复用
func (baseMediaSubSession *BaseMediaSubSession) SetMultiplexRTCPWithRTP(b bool) {
	baseMediaSubSession.MultiplexRTCPWithRTP = b
}

// GenerateRTPMapLine a=rtpmap
func (baseMediaSubSession *BaseMediaSubSession) GenerateRTPMapLine() string {
	playload := baseMediaSubSession.rtpPayload
	if playload >= 96 {
		encodingParamsPart := ""
		if baseMediaSubSession.numChannels != 1 {
			encodingParamsPart = fmt.Sprintf("/%d", baseMediaSubSession.numChannels)
		}

		return fmt.Sprintf("a=rtpmap:%d %s/%d%s\r\n", playload, baseMediaSubSession.rtpPayloadFormatName, baseMediaSubSession.rtpTimestampFrequency, encodingParamsPart)
	}
	return ""
}

// SetParentDuration 来到低级计算Duration
func (baseMediaSubSession *BaseMediaSubSession) SetParentDuration(d float64) {
	baseMediaSubSession.ParentDuration = d
}

// GenerateRangeLine a=range:clock=
func (baseMediaSubSession *BaseMediaSubSession) GenerateRangeLine(absStartTime, absEndTime string, duration float64) string {
	// 支持绝对时间
	if absStartTime != "" {
		if absEndTime != "" {
			return fmt.Sprintf("a=range:clock=%s-%s\r\n", absStartTime, absEndTime)
		}

		return fmt.Sprintf("a=range:clock=%s-\r\n", absStartTime)
	}

	if baseMediaSubSession.ParentDuration >= 0.0 {
		return ""
	}

	if duration == 0.0 {
		return "a=range:npt=0-\r\n"
	}

	return fmt.Sprintf("a=range:npt=0-%.3f\r\n", duration)
}

func (baseMediaSubSession *BaseMediaSubSession) GetBitrate() int {
	return baseMediaSubSession.EstBitrate
}

func (baseMediaSubSession *BaseMediaSubSession) GetClockRate() int {
	return baseMediaSubSession.rtpTimestampFrequency
}
