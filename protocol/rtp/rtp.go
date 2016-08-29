package rtp

import (
	"fmt"
	"net"
)

// RTP 传输方式
const (
	RTP_UDP = iota
	RTP_TCP
	RAW_UDP
)

type RTPServer struct{}

func NewRTPServer() *RTPServer {
	return &RTPServer{}
}

func (rtpServer *RTPServer) Serve(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	udpConn, err := net.ListenUDP(addr.Network(), addr)
	if err != nil {
		return err
	}

	defer udpConn.Close()
	fmt.Println("RTP listen to :", address)
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
