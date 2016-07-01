package koala

import (
    "time"
    "net"
    "log"
    "strings"
)

type RTSPRequest struct {
    BaseRequest
}

func (rtspRequest *RTSPRequest) GetMethod() string {
    return ""
}

func (rtspRequest *RTSPRequest) Recv() {
     ip := net.ParseIP(strings.Split( rtspRequest.RemoteAddr.String(), ":")[0])
     log.Printf("new connected from:%v\n", ip)

     for {
         p      := make([]byte, 4096)
         n, err := rtspRequest.Socket.Read(p)
         if err != nil {
             log.Printf("error receiving, bytes:%d reason:%v\n", n, err)
             break
         }

         select {
             case rtspRequest.in <- p:
             case <-time.After(30 * time.Second):
                  log.Printf("server busy or listen closed.")
         }
     }

     log.Printf("Client shutdown:%v", ip)
}

func NewRTSPRequest( socket net.Conn ) *RTSPRequest {
    req := new(RTSPRequest)
    req.Socket = socket
    req.RemoteAddr = socket.RemoteAddr()
    req.LocalAddr  = socket.LocalAddr()
    req.in     = make(chan []byte, 0)
    return req
}