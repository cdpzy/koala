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

import (
    "net"
    "sync"
    "log"
    "runtime"
    "strings"
    "time"
)

type RTSPTCPConnection struct {
    socket net.Conn
    in       chan []byte
	out      chan []byte
    sendLock sync.RWMutex
    closein  chan bool
    closeout chan bool
    RemoteAddr net.Addr
    LocalAddr  net.Addr
}

func (rtspTCPConnection *RTSPTCPConnection) Recv(){
    defer func(){
        rtspTCPConnection.Close()
        if r := recover(); r != nil {
            if _, ok := r.(runtime.Error); ok {
				log.Printf("connection Running fail:%v\n", r)
			}
        }
    }()

    socket := rtspTCPConnection.socket
    rtspTCPConnection.RemoteAddr = socket.RemoteAddr()
    rtspTCPConnection.LocalAddr  = socket.LocalAddr()

	ip := net.ParseIP(strings.Split( socket.RemoteAddr().String(), ":")[0])
	log.Printf("new connected from:%v\n", ip)

    go rtspTCPConnection.request()
    go rtspTCPConnection.response()

    for {
        socket.SetReadDeadline(time.Now().Add(120 * time.Second))
        packet := make([]byte, 4096)
        n, err := socket.Read(packet)
        if err != nil {
            log.Printf("error receiving header, bytes:%d reason:%v\n", n, err)
            break
        }

        select {
            case rtspTCPConnection.in <- packet:
            case <-time.After(120 * time.Second):
                log.Printf("server busy or listen closed.")
        }
    }

    log.Printf("Client shutdown:%v", ip)
}

func (rtspTCPConnection *RTSPTCPConnection) request(){
    defer func(){
        close(rtspTCPConnection.in)
		close(rtspTCPConnection.closein)
		log.Print("Connection request stoped")
    }()

    for {
        select {
        case packet, ok := <-rtspTCPConnection.in :
             if !ok {
                 log.Printf("error receiving packet, packet:%v", packet)
                 break
             }

        case <-rtspTCPConnection.closein:
             return
        }
    }
}

func (rtspTCPConnection *RTSPTCPConnection) response(){
    defer func(){
        close(rtspTCPConnection.out)
		close(rtspTCPConnection.closeout)
		log.Print("Connection response stoped")
    }()

    for {
        select {
        case raw, ok := <-rtspTCPConnection.out:
             if !ok {
                 break
             }

             log.Printf("write:%v", raw)
             rtspTCPConnection.socket.Write(raw)

        case <-rtspTCPConnection.closeout:
             return
        }
    }
}

func (rtspTCPConnection *RTSPTCPConnection) Send( b []byte ) {
    rtspTCPConnection.sendLock.Lock()
    defer rtspTCPConnection.sendLock.Unlock()

    rtspTCPConnection.out <- b
}

func (rtspTCPConnection *RTSPTCPConnection) Close() {
    rtspTCPConnection.closein  <- true
    rtspTCPConnection.closeout <- true
    rtspTCPConnection.socket.Close()
}

func NewRTSPTCPConnection( socket net.Conn ) *RTSPTCPConnection {
    return &RTSPTCPConnection{
        socket  : socket,
        in      : make(chan []byte, 0),
        out     : make(chan []byte, 0),
        closein : make(chan bool),
		closeout: make(chan bool),
    }
}