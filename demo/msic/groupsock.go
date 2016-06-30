package msic

import (
    "net"
    "fmt"
)

type Socket struct {
}

type OutputSocket struct {
	Socket
	sourcePort  uint
	lastSentTTL uint
}

func (this *OutputSocket) write(destAddr string, portNum uint, buffer []byte, bufferSize uint) bool {
	udpConn := SetupDatagramSocket(destAddr, portNum)
	return writeSocket(udpConn, buffer)
}

func (this *OutputSocket) sourcePortNum() uint {
	return this.sourcePort
}

type GroupSock struct {
	OutputSocket
	dests   []*destRecord
	portNum uint
	ttl     uint
}

func NewGroupSock(addrStr string, portNum uint) *GroupSock {
	gs := new(GroupSock)
	gs.ttl = 255
	gs.portNum = portNum
	gs.AddDestination(addrStr, portNum)
	return gs
}

func (this *GroupSock) Output(buffer []byte, bufferSize, ttlToSend uint) bool {
	var writeSuccess bool
	for i := 0; i < len(this.dests); i++ {
		dest := this.dests[i]
		if this.write(dest.addrStr, dest.portNum, buffer, bufferSize) {
			writeSuccess = true
		}
	}
	return writeSuccess
}

func (this *GroupSock) handleRead() {
}

func (this *GroupSock) TTL() uint {
	return this.ttl
}

func (this *GroupSock) AddDestination(addr string, port uint) {
	this.dests = append(this.dests, NewDestRecord(addr, port))
}

func (this *GroupSock) delDestination(addr string, port uint) {
}

type destRecord struct {
	addrStr string
	portNum uint
}

func NewDestRecord(addr string, port uint) *destRecord {
	dest := new(destRecord)
	dest.addrStr = addr
	dest.portNum = port
	return dest
}


func SetupDatagramSocket(address string, port uint) *net.UDPConn {
	addr := fmt.Sprintf("%s:%d", address, port)
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil
	}
	return udpConn
}

func writeSocket(conn net.Conn, buffer []byte) bool {
	_, err := conn.Write(buffer)
	if err != nil {
		//fmt.Println(writeBytes)
		return false
	}

	return true
}
