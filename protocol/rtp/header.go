package rtp

import (
	"fmt"
	"regexp"
	"strings"
)

type TransportHeader struct {
	StreamingMode     byte
	StreamingModeName string
	DestinationTTL    int
	ClientRTPPortNum  int
	ClientRTCPPortNum int
	RTPChannelID      int
	RTCPChannelID     int
	DestinationAddr   string
}

type RangeHeader struct {
	Start        float64
	End          float64
	AbsStartTime string
	AbsEndTime   string
	IsNow        bool
}

// ParseTransportHeader parse
func ParseTransportHeader(s string) *TransportHeader {
	var (
		ttl     int
		num     int
		p1      int
		p2      int
		rtpCid  int
		rtcpCid int
	)

	tp := new(TransportHeader)
	tp.StreamingMode = RTP_UDP
	tp.DestinationTTL = 255
	tp.ClientRTPPortNum = 0
	tp.ClientRTCPPortNum = 1
	tp.RTPChannelID = 255
	tp.RTCPChannelID = 255

	re := regexp.MustCompile("(;|\r\n)")
	data := re.Split(s, -1)

	for _, item := range data {
		item = strings.TrimSpace(item)
		if strings.EqualFold(item, "RTP/AVP/TCP") {
			tp.StreamingMode = RTP_TCP
		} else if strings.EqualFold(item, "RAW/RAW/UDP") || strings.EqualFold(item, "MP2T/H2221/UDP") {
			tp.StreamingMode = RAW_UDP
			tp.StreamingModeName = item
		} else if strings.Index(item, "destination=") != -1 {
			tp.DestinationAddr = item[12:]
		} else if num, _ = fmt.Sscanf(item, "ttl%d", &ttl); num == 1 {
			tp.DestinationTTL = ttl
		} else if num, _ = fmt.Sscanf(item, "client_port=%d-%d", &p1, &p2); num == 2 {
			tp.ClientRTPPortNum = p1
			if tp.StreamingMode == RAW_UDP {
				tp.ClientRTCPPortNum = 0
			} else {
				tp.ClientRTCPPortNum = p2
			}
		} else if num, _ = fmt.Sscanf(item, "client_port=%s", &p1); num == 1 {
			tp.ClientRTPPortNum = p1
			if tp.StreamingMode == RAW_UDP {
				tp.ClientRTCPPortNum = 0
			} else {
				tp.ClientRTCPPortNum = p1
			}
		} else if num, _ = fmt.Sscanf(item, "interleaved=%d-%d", &rtpCid, &rtcpCid); num == 2 {
			tp.RTPChannelID = rtpCid
			tp.RTCPChannelID = rtcpCid
		}

	}
	return tp
}

// ParseRangeHeader npt =
func ParseRangeHeader(s string) (*RangeHeader, error) {
	var (
		start           float64
		end             float64
		numCharsMatched int
	)
	rangeHeader := new(RangeHeader)
	num, err := fmt.Sscanf(s, "npt = %f - %f", &start, &end)
	if err != nil {
		return nil, err
	}

	if num == 2 {
		rangeHeader.Start = start
		rangeHeader.End = end
		return rangeHeader, nil
	}

	num, err = fmt.Sscanf(s, "npt = %f -", &start)
	if err != nil {
		return nil, err
	}

	if num == 1 {
		rangeHeader.Start = start
		return rangeHeader, nil
	}

	if strings.EqualFold(s, "npt = now -") {
		rangeHeader.Start = 0.0
		rangeHeader.End = 0.0
		return rangeHeader, nil
	}

	num, err = fmt.Sscanf(s, "clock = %n", &numCharsMatched)
	if err != nil {
		return nil, err
	}

	if numCharsMatched > 0 {
		as, ae := "", ""
		utcTimes := string(s[numCharsMatched:])
		num, err = fmt.Sscanf(utcTimes, "%[^-]-%s", &as, &ae)
		if err != nil {
			return nil, err
		}

		if num == 2 {
			rangeHeader.AbsStartTime = as
			rangeHeader.AbsEndTime = ae
		} else if num == 1 {
			rangeHeader.AbsStartTime = as
		}
	}

	return rangeHeader, nil
}
