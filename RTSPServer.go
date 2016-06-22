/**
* The MIT License (MIT)
*
* Copyright (c) 2016 doublemo<435420057@qq.com>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
*/

package koala

import(
    "net"
    "errors"
)

const (
    SERVER_NET_PROTO_HTTP = "http"
    SERVER_NET_PROTO_TCP  = "tcp"
    SERVER_NET_PROTO_UDP  = "udp"
)

type RTSPServer struct {
    tcpListen    net.Listener
}

func (RTSP *RTSPServer) supportTCP( addr string ) error {
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        return err
    }

    lis, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        return err
    }

    RTSP.tcpListen = lis
    defer RTSP.tcpListen.Close()

    for{
        conn, err := lis.Accept()
        if err != nil {
            return err
        }

        go func(socket net.Conn) {
             NewRTSPTCPConnection( socket ).Recv()
        }(conn)
    }
    return nil
}

func (RTSP *RTSPServer) supportUDP( addr string ) error {
    return nil
}

func (RTSP *RTSPServer) supportHTTP( addr string ) error {
   return nil
}

func (RTSP *RTSPServer) Serve( proto, addr string ) error {
    switch proto {
    case SERVER_NET_PROTO_HTTP:
         return RTSP.supportHTTP( addr )

    case SERVER_NET_PROTO_TCP:
         return RTSP.supportTCP( addr )

    case SERVER_NET_PROTO_UDP:
         return RTSP.supportUDP( addr )

    default:
        return errors.New("Not support:" + proto)
    }
    return nil
}

func (RTSP *RTSPServer) Stop() {

}

func NewRTSPServer() *RTSPServer {
    return &RTSPServer{}
}