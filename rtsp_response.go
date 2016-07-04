package koala

import (
    "net"
    "io"
)

type RTSPResponse struct {
    BaseResponse
    Socket net.Conn
    out chan []byte
}

func (rtspResponse *RTSPResponse) NotFound() {

}

func (rtspResponse *RTSPResponse) Recv() {
    for {
        select {
            case raw, ok := <-rtspResponse.out:
                 if !ok {
                     break
                 }
                rtspResponse.Socket.Write(raw)
        }
    }
}

func (rtspResponse *RTSPResponse) Write( b []byte ) error {
    rtspResponse.out <- b
    return nil
} 

func NewRTSPResponse( socket net.Conn ) *RTSPResponse {
    resp := new(RTSPResponse)
    resp.Socket = socket
    resp.out    = make(chan []byte, 0)
    return resp
}