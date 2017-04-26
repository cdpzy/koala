package net

import (
	"net"

	"github.com/doublemo/koala/helper"

	log "github.com/Sirupsen/logrus"
	"github.com/xtaci/kcp-go"
)

// KCPServerOptions KCP服务器
type KCPServerOptions struct {
	Addr            string      //
	ReadBufferSize  int         // per connection tcp socket buffer 32767
	WriteBufferSize int         // UDP listener socket buffer 4194304
	Dscp            int         // set DSCP(6bit) 46
	Sndwnd          int         // per connection UDP send window 32
	Rcvwnd          int         // per connection UDP recv window 32
	Nodelay         int         // ikcp_nodelay()  1
	Interval        int         // ikcp_nodelay()  20
	Resend          int         // ikcp_nodelay() 1
	Nc              int         // ikcp_nodelay() 1
	Mtu             int         // MTU of UDP packets, without IP(20) + UDP(8) 1280
	ClientHandler   HandlerFunc //
}

// KCPServer ..
type KCPServer struct {
	Addr            string       //
	ReadBufferSize  int          //
	WriteBufferSize int          //
	Dscp            int          //
	Sndwnd          int          //
	Rcvwnd          int          //
	Nodelay         int          //
	Interval        int          //
	Resend          int          //
	Nc              int          //
	Mtu             int          //
	ClientHandler   HandlerFunc  //
	listener        net.Listener //
}

// Serve 启动服务
func (s *KCPServer) Serve() error {
	var err error
	s.listener, err = kcp.Listen(s.Addr)
	if err != nil {
		return err
	}

	lis := s.listener.(*kcp.Listener)
	log.Infoln("KCP listening on:", s.listener.Addr())

	if err = lis.SetReadBuffer(s.ReadBufferSize); err != nil {
		return err
	}

	if err = lis.SetWriteBuffer(s.WriteBufferSize); err != nil {
		return err
	}

	if err = lis.SetDSCP(s.Dscp); err != nil {
		return err
	}

	for {
		conn, err := lis.AcceptKCP()
		if err != nil {
			log.Warningln("accept failed:", err)
			continue
		}

		conn.SetWindowSize(s.Sndwnd, s.Rcvwnd)
		conn.SetNoDelay(s.Nodelay, s.Interval, s.Resend, s.Nc)
		conn.SetStreamMode(true)
		conn.SetMtu(s.Mtu)

		go s.handle(conn)
	}
}

// handle 处理
func (s *KCPServer) handle(conn net.Conn) {
	defer func() {
		helper.RecoverStack()
		conn.Close()
	}()

	s.ClientHandler(conn)
}

// NewKCPServer new
func NewKCPServer(op *KCPServerOptions) *KCPServer {
	return &KCPServer{
		Addr:            op.Addr,
		ReadBufferSize:  op.ReadBufferSize,
		WriteBufferSize: op.WriteBufferSize,
		Dscp:            op.Dscp,
		Sndwnd:          op.Sndwnd,
		Rcvwnd:          op.Rcvwnd,
		Nodelay:         op.Nodelay,
		Interval:        op.Interval,
		Resend:          op.Resend,
		Nc:              op.Nc,
		Mtu:             op.Mtu,
		ClientHandler:   op.ClientHandler,
	}
}
