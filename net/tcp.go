package net

import (
	"fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// HandlerFunc 处理连接
type HandlerFunc func(net.Conn, chan struct{})

// TCPServerOptions TCP服务器配置
type TCPServerOptions struct {
	Addr            string      // 监听地址
	ReadBufferSize  int         // 读取缓存大小 32767
	WriteBufferSize int         // 写入缓存大小 32767
	ClientHandler   HandlerFunc // 用户接入后处理
}

// TCPServer TCP 服务器
type TCPServer struct {
	Addr            string           //
	ReadBufferSize  int              //
	WriteBufferSize int              //
	ClientHandler   HandlerFunc      //
	wg              sync.WaitGroup   //
	listener        *net.TCPListener //
	closed          chan struct{}    //
}

// Serve 开启服务
func (s *TCPServer) Serve() error {
	tcpaddr, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.listener, err = net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		return err
	}

	s.closed = make(chan struct{})
	defer func() {
		if s.listener != nil {
			s.listener.Close()
		}
		close(s.closed)
	}()

	log.Infoln("TCP listening on:", s.listener.Addr())
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			return err
		}

		conn.SetReadBuffer(s.ReadBufferSize)
		conn.SetWriteBuffer(s.WriteBufferSize)

		s.wg.Add(1)
		go s.handle(conn)
	}
}

func (s *TCPServer) handle(conn net.Conn) {
	defer s.wg.Done()
	defer func() {
		conn.Close()
		fmt.Println("Conn Closed", conn.RemoteAddr())
	}()

	s.ClientHandler(conn, s.closed)
}

func (s *TCPServer) Close() {
	if s.listener != nil {
		lis := s.listener
		s.listener = nil
		lis.Close()
	}

	s.wg.Wait()
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(op *TCPServerOptions) *TCPServer {
	return &TCPServer{
		ReadBufferSize:  op.ReadBufferSize,
		WriteBufferSize: op.WriteBufferSize,
		Addr:            op.Addr,
		ClientHandler:   op.ClientHandler,
	}
}
