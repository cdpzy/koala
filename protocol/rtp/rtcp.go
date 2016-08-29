package rtp

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

const (
	MAX_RTCP_PACKET_SIZE int = 1456

	// bytes (1500, minus some allowance for IP, UDP, UMTP headers)
	PREFERRED_RTCP_PACKET_SIZE = 1000
)

// RTCP packet types:
const (
	RTCP_PT_SR int = 200 + iota
	RTCP_PT_RR
	RTCP_PT_SDES
	RTCP_PT_BYE
	RTCP_PT_APP
	RTCP_PT_RTPFB // Generic RTP Feedback [RFC4585]
	RTCP_PT_PSFB  // Payload-specific [RFC4585]
	RTCP_PT_XR    // extended report [RFC3611]
	RTCP_PT_AVB   // AVB RTCP packet ["Standard for Layer 3 Transport Protocol for Time Sensitive Applications in Local Area Networks." Work in progress.]
	RTCP_PT_RSI   // Receiver Summary Information [RFC5760]
	RTCP_PT_TOKEN // Port Mapping [RFC6284]
	RTCP_PT_IDMS  // IDMS Settings [RFC7272]
)

// SDES tags:
const (
	RTCP_SDES_END int = iota
	RTCP_SDES_CNAME
	RTCP_SDES_NAME
	RTCP_SDES_EMAIL
	RTCP_SDES_PHONE
	RTCP_SDES_LOC
	RTCP_SDES_TOOL
	RTCP_SDES_NOTE
	RTCP_SDES_PRIV
)

// RTCP member type
const (
	RTCP_MEMBER_TYPE_SR int = iota // 发送者
	RTCP_MEMBER_TYPE_RR            // 接收者
)

// RTPC 包设置 、RFC 3550, chapters 6.3.1 and A.7
const (
	RTCP_MIN_TIME     float64 = 5.0  // 最小发包时间
	RTCP_SR_FR        float64 = 0.25 // 发送者分值
	RTCP_RR_FR        float64 = 0.75 // 接者分值
	RTCP_COMPENSATION float64 = 2.71828 - 1.5
)

type RTCPMember struct {
	SSRC           uint32
	Type           string
	Active         bool
	LastPacketTime int64
}

type RTCPMemberManager struct {
	mux     sync.RWMutex
	members map[uint32]*RTCPMember
}

type RTCPAdapter struct {
	SRMember      *RTCPMemberManager
	RRMember      *RTCPMemberManager
	CName         string
	StreamBitrate int
	WeSent        bool
	AvrgSize      float64
	Initial       bool
	ClockRate     int
}

func NewRTCPMemberManager() *RTCPMemberManager {
	return &RTCPMemberManager{members: make(map[uint32]*RTCPMember)}
}

func NewRTCPAdapter(StreamBitrate int, cname string, ClockRate int) *RTCPAdapter {
	return &RTCPAdapter{
		SRMember:      NewRTCPMemberManager(),
		RRMember:      NewRTCPMemberManager(),
		CName:         cname,
		StreamBitrate: StreamBitrate,
		ClockRate:     ClockRate,
	}
}

func NewRTCPMember(ssrc uint32) *RTCPMember {
	return &RTCPMember{SSRC: ssrc, Active: true}
}

func (rtcpMemberManager *RTCPMemberManager) Add(member *RTCPMember) error {
	rtcpMemberManager.mux.Lock()
	defer rtcpMemberManager.mux.Unlock()

	if rtcpMemberManager.IsMember(member.SSRC) {
		return errors.New("SSRCFound")
	}

	rtcpMemberManager.members[member.SSRC] = member
	return nil
}

func (rtcpMemberManager *RTCPMemberManager) Remove(ssrc uint32) {
	rtcpMemberManager.mux.Lock()
	defer rtcpMemberManager.mux.Unlock()

	delete(rtcpMemberManager.members, ssrc)
}

func (rtcpMemberManager *RTCPMemberManager) IsMember(ssrc uint32) bool {
	if _, ok := rtcpMemberManager.members[ssrc]; ok {
		return true
	}

	return false
}

func (rtcpMemberManager *RTCPMemberManager) GetAll() map[uint32]*RTCPMember {
	return rtcpMemberManager.members
}

func (rtcpMemberManager *RTCPMemberManager) NumMember() int {
	return len(rtcpMemberManager.members)
}

func (rtcpAdapter *RTCPAdapter) Run(clientPort, serverPort string) {
	go func() {
		err := rtcpAdapter.Serve(serverPort)
		if err != nil {
			fmt.Println("failed to RTCP listen:", serverPort, "(", err, ")")
		}
	}()

	bw := float64(rtcpAdapter.ClockRate) * 8.0 / 20.0
	if bw < 1 {
		bw = 64000.0 / 20.0
	}

	// 	const (
	//     senderInfoLen  = 20
	//     reportBlockLen = 24
	// )
	avg := rtcpAdapter.SRMember.NumMember()*20 + 24 + 20
	ti, td := rtcpAdapter.CalNextPacket(bw, float64(avg), rtcpAdapter.WeSent, rtcpAdapter.Initial)
	nextTime := ti + time.Now().UnixNano()
	go rtcpAdapter.Loop(ti, td, nextTime)
}

func (rtcpAdapter *RTCPAdapter) Loop(ti, td, next int64) {
	granularity := time.Duration(250e6) // 250 ms
	//ssrcTimeout := 5 * td
	// dataTimeout := 2 * ti
	ticker := time.NewTicker(granularity)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			if now < next {
				continue
			}

			rrMembers := rtcpAdapter.RRMember.GetAll()
			for ssrc, rr := range rrMembers {
				switch rr.Active {
				case true:

				}

				fmt.Println(ssrc)
			}
		}
	}
}

func (rtcpAdapter *RTCPAdapter) Serve(serverPort string) error {
	addr, err := net.ResolveUDPAddr("udp", serverPort)
	if err != nil {
		return err
	}

	udpConn, err := net.ListenUDP(addr.Network(), addr)
	if err != nil {
		return err
	}

	defer udpConn.Close()
	fmt.Println("RTCP listen to :", serverPort)
	for {
		data := make([]byte, MAX_RTCP_PACKET_SIZE)
		n, remoteAddr, err := udpConn.ReadFromUDP(data[0:])
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err.Error())
			return err
		}

		fmt.Println("UDPREAD:", string(data), remoteAddr, n)
	}
	return nil
}

// CalNextPacket 计算RTCP包发送时间
// RFC 3550
// rtcpBW 目标 RTCP 带宽。例如用于会话中所有成员的 RTCP 带宽。单位 bit/s。这将是程序开始时,指定给“会话带宽”参数的一部分。
// weSent 自当前第二个前面的 RTCP 发送后,应用程序又发送了数据,则此项为 true。
// avg_rtcp_size: 此参与者收到的和发送的 RTCP 复合包的平均大小。单位:bit。按 6.2 节,此大小包括底层传输层和网络层协议头。
// initial: 如果应用程序还未发送 RTCP 包,则标记为 true。
func (rtcpAdapter *RTCPAdapter) CalNextPacket(rtcpBW, avrgSize float64, weSent, initial bool) (int64, int64) {
	min := RTCP_MIN_TIME
	if initial {
		min /= 2
	}

	srNum := rtcpAdapter.SRMember.NumMember()
	rrNum := rtcpAdapter.RRMember.NumMember()
	members := srNum + rrNum
	n := members
	if srNum <= int(float64(members)*RTCP_SR_FR) {
		if weSent {
			rtcpBW *= RTCP_SR_FR
			n = srNum
		} else {
			rtcpBW *= RTCP_RR_FR
			n = members - srNum
		}
	}

	t := math.Max(float64(n)*(avrgSize/rtcpBW), min)
	td := int64(t * 1e9)

	buffer := make([]byte, 2)
	rand.Read(buffer[:])

	rndNo := uint16(buffer[0])
	rndNo |= uint16(buffer[1]) << 8
	rndF := float64(rndNo)/65536.0 + 0.5

	t *= rndF
	t /= RTCP_COMPENSATION

	// return as nanoseconds
	return int64(t * 1e9), td
}
