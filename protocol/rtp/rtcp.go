package rtp

import (
	"errors"
	"fmt"
	"net"
	"sync"
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

type RTCPMember struct {
	SSRC uint32
	Type string
}

type RTCPMemberManager struct {
	mux     sync.RWMutex
	members map[uint32]*RTCPMember
}

type RTCPAdapter struct {
	Member        *RTCPMemberManager
	CName         string
	StreamBitrate int
}

func NewRTCPMemberManager() *RTCPMemberManager {
	return &RTCPMemberManager{members: make(map[uint32]*RTCPMember)}
}

func NewRTCPAdapter(StreamBitrate int, cname string) *RTCPAdapter {
	return &RTCPAdapter{
		Member:        NewRTCPMemberManager(),
		CName:         cname,
		StreamBitrate: StreamBitrate,
	}
}

func NewRTCPMember(ssrc uint32) *RTCPMember {
	return &RTCPMember{SSRC: ssrc}
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

func (rtcpMemberManager *RTCPMemberManager) NumMember() int {
	return len(rtcpMemberManager.members)
}

func (rtcpAdapter *RTCPAdapter) Run(clientPort, serverPort string) {
	addr, _ := net.ResolveUDPAddr("udp", serverPort)
	udpConn, err := net.ListenUDP(addr.Network(), addr)
	if err != nil {
		panic(err)
	}

	defer udpConn.Close()
	fmt.Println("UDP:--", serverPort)
	for {
		data := make([]byte, MAX_RTCP_PACKET_SIZE)
		n, remoteAddr, err := udpConn.ReadFromUDP(data[0:])
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err.Error())
			return
		}

		fmt.Println("UDPREAD:", string(data), remoteAddr, n)
	}
}
