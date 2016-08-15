package koala

import (
	"log"
	"net"
	"runtime"
	"time"
)

type RTSPServer struct {
	listener    net.Listener
	handlerFunc HandlerFunc
}

func (rtspServer *RTSPServer) HandlerFunc(handlerFunc HandlerFunc) {
	rtspServer.handlerFunc = handlerFunc
}

func (rtspServer *RTSPServer) Serve(addr string) error {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	rtspServer.listener, err = net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}

	defer rtspServer.Close()

	go func() {
		time.Sleep(100 * time.Millisecond)
		log.Printf("Listening on %s...\n", addr)
	}()

	for {
		conn, err := rtspServer.listener.Accept()
		if err != nil {
			return err
		}

		go func(socket net.Conn) {
			defer func() {
				socket.Close()
				if r := recover(); r != nil {
					if _, ok := r.(runtime.Error); ok {
						log.Printf("RTSPServer runtime error:%v\n", r)
					}
				}
			}()

			req := NewRTSPRequest(socket)
			resp := NewRTSPResponse(socket)

			go req.Recv()
			go resp.Recv()
			rtspServer.handlerFunc(req, resp)
		}(conn)
	}
	return nil
}

func (rtspServer *RTSPServer) Close() {
	rtspServer.listener.Close()
}

func NewRTSPServer() *RTSPServer {
	return new(RTSPServer)
}
