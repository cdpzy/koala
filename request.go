package koala

import (
    "net"
    "golang.org/x/net/websocket"
)

type Request interface {
    GetSocket() net.Conn
    GetWebsocket() *websocket.Conn
    Recv()
    GetInputChan() chan []byte
    GetRemoteAddr() net.Addr
    GetLocalAddr() net.Addr
    GetMethod() string
    String() string
}


type BaseRequest struct {
    Socket    net.Conn
    Websocket *websocket.Conn
    RemoteAddr net.Addr
    LocalAddr  net.Addr
    in        chan []byte
}

func (baseRequest *BaseRequest) GetSocket() net.Conn {
    return baseRequest.Socket
}

func (baseRequest *BaseRequest) GetWebsocket() *websocket.Conn {
    return baseRequest.Websocket
}

func (baseRequest *BaseRequest) Recv(){}


func (baseRequest *BaseRequest) GetInputChan() chan []byte {
    return baseRequest.in
}

func (baseRequest *BaseRequest) GetRemoteAddr() net.Addr {
    return baseRequest.RemoteAddr
}

func (baseRequest *BaseRequest) GetLocalAddr() net.Addr {
    return baseRequest.LocalAddr
}

func (baseRequest *BaseRequest) String() string{
    return ""
}