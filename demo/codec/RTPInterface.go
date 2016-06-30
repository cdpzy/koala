package codec

import (
     "github.com/doublemo/koala/msic"
     "net"
)

type RTPInterface struct {
	gs                         *msic.GroupSock
	owner                      interface{}
	auxReadHandlerFunc         interface{}
	tcpStreams                 *TCPStreamRecord
	nextTCPReadStreamSocketNum int
}

func NewRTPInterface(owner interface{}, gs *msic.GroupSock) *RTPInterface {
	rtpInterface := new(RTPInterface)
	rtpInterface.gs = gs
	rtpInterface.owner = owner
	return rtpInterface
}

func (this *RTPInterface) startNetworkReading( /*handlerProc interface*/ ) {
}

func (this *RTPInterface) stopNetworkReading() {
}

func (this *RTPInterface) GS() *msic.GroupSock {
	return this.gs
}

func (this *RTPInterface) addStreamSocket(sockNum net.Conn, streamChannelId uint) {
	if sockNum == nil {
		return
	}

	this.tcpStreams = NewTCPStreamRecord(sockNum, streamChannelId)
}

func (this *RTPInterface) delStreamSocket() {
}

func (this *RTPInterface) sendPacket(packet []byte, packetSize uint) bool {
	return this.gs.Output(packet, packetSize, this.gs.TTL())
}

func (this *RTPInterface) handleRead() bool {
	return true
}

type TCPStreamRecord struct {
	streamChannelId uint
	streamSocketNum net.Conn
}

func NewTCPStreamRecord(streamSocketNum net.Conn, streamChannelId uint) *TCPStreamRecord {
	tcpStreamRecord := new(TCPStreamRecord)
	tcpStreamRecord.streamChannelId = streamChannelId
	tcpStreamRecord.streamSocketNum = streamSocketNum
	return tcpStreamRecord
}

///////////// Help Functions ///////////////
func sendRTPOverTCP(socketNum net.Conn, packet []byte, packetSize, streamChannelId uint) {
	dollar := []byte{'$'}
	channelId := []byte{byte(streamChannelId)}
	socketNum.Write(dollar)
	socketNum.Write(channelId)
}