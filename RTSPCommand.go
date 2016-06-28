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
    "log"
    "fmt"
    "errors"
    //"github.com/doublemo/koala/msic"
    //"strings"
)

// interface
type IRTSPCommand interface {
    ParseCommand( *RTSPTCPConnection ) error

    handelOptions( *RTSPTCPConnection ) error

    handelDESCRIBE( *RTSPTCPConnection  ) error

    handelSETUP(*RTSPTCPConnection) error

    handelTEARDOWN() error

    handelPLAY(*RTSPTCPConnection) error

    handelPAUSE() error

    handelGET_PARAMETER() error

    handelSET_PARAMETER() error

    // Allow RTSP Methed
    GetAllowCommand() []string
}

type RTSPCommand struct {
    rtsp *RTSPServer
}

func (rtspCommand *RTSPCommand) ParseCommand( client *RTSPTCPConnection ) error {
    req := rtspCommand.rtsp.GetRequest()
    log.Println("method:", req.GetMethod())
    switch req.GetMethod() {
    case "OPTIONS":
        rtspCommand.handelOptions( client )

    case "DESCRIBE":
        rtspCommand.handelDESCRIBE( client )

    case "SETUP":
        rtspCommand.handelSETUP( client )

    case "PLAY":
        rtspCommand.handelPLAY( client )
    }
    
    return nil
}

func (rtspCommand *RTSPCommand) handelOptions( client *RTSPTCPConnection ) error {
    header  := rtspCommand.rtsp.GetRequest().GetHeader()
    cseq    := header.Get("CSeq")
    buf     := fmt.Sprintf("RTSP/1.0 200 OK\r\n"+
		"CSeq: %s\r\n"+
		"%sPublic: %s\r\n\r\n",
		cseq, DateHeader(), rtspCommand.GetAllowCommand())

    client.Send([]byte(buf))
    return nil
}

func (rtspCommand *RTSPCommand) handelDESCRIBE( client *RTSPTCPConnection ) error {
    header  := rtspCommand.rtsp.GetRequest().GetHeader()
    url     := rtspCommand.rtsp.GetRequest().GetURL()
    cseq    := header.Get("CSeq")

    session := rtspCommand.rtsp.LookupServerMediaSession("test.264")
    //log.Println("session:", session)
    if session == nil {
        rtspCommand.handleCommandNotFound( client )
        return nil
    }

    sdpDescription := session.GenerateSDPDescription()
    sdpDescriptionSize := len(sdpDescription)
    log.Println("sdpDescriptionSize:", sdpDescriptionSize)
    if sdpDescriptionSize < 1 {
        rtspCommand.handleCommandNotFound( client )
        return errors.New("404 File Not Found, Or In Incorrect Format")
    }

    buf :=  fmt.Sprintf("RTSP/1.0 200 OK\r\n"+
		"CSeq: %s\r\n"+
		"%s"+
		"Content-Base: %s/\r\n"+
		"Content-Type: application/sdp\r\n"+
		"Content-Length: %d\r\n\r\n"+
		"%s",
		cseq, DateHeader(), url.String(), sdpDescriptionSize, sdpDescription)
    
    log.Println("Buffer:", buf)
    client.Send([]byte(buf))
    return nil
}

func (rtspCommand *RTSPCommand) handelSETUP(client *RTSPTCPConnection) error {
    header  := rtspCommand.rtsp.GetRequest().GetHeader()
    cseq    := header.Get("CSeq")

    buf :=  fmt.Sprintf("RTSP/1.0 200 OK\r\n"+
				"CSeq: %s\r\n"+
				"%s"+
				"Transport: RTP/AVP;multicast;destination=%s;source=%s;port=%d-%d;ttl=%d\r\n"+
				"Session: %08X\r\n\r\n", cseq,
				DateHeader(),
				client.RemoteAddr,
				client.LocalAddr,
				6970,
				6971,
				220,
				"sdd855522")

    log.Println("Buffer:", buf)
    client.Send([]byte(buf)) 
    return nil
}

func (rtspCommand *RTSPCommand) handelTEARDOWN() error {
    return nil
}

func (rtspCommand *RTSPCommand) handelPLAY(client *RTSPTCPConnection) error {
    // PLAY rtsp://127.0.0.1:554/test.264/ RTSP/1.0
    // Cseq: 5
    // User-Agent: LibVLC/2.2.4 (LIVE555 Streaming Media v2016.02.22)
    // Session: 736464383535353232
    // Range: npt=0.000-
    log.Println(rtspCommand.rtsp.GetRequest().String())
    return nil
}

func (rtspCommand *RTSPCommand) handelPAUSE() error {
    return nil
}

func (rtspCommand *RTSPCommand) handelGET_PARAMETER() error {
    return nil
}

func (rtspCommand *RTSPCommand) handelSET_PARAMETER() error {
    return nil
}

func (rtspCommand RTSPCommand) GetAllowCommand() []string {
    return []string{
        "OPTIONS", 
        "DESCRIBE", 
        "SETUP", 
        "TEARDOWN", 
        "PLAY", 
        "PAUSE", 
        "GET_PARAMETER", 
        "SET_PARAMETER",
    }
}

func (rtspCommand RTSPCommand) handleCommandNotFound( client *RTSPTCPConnection ) {
    buf := fmt.Sprintf("HTTP/1.0 404 Not Found\r\n%s\r\n\r\n", DateHeader())
    client.Send([]byte(buf))
}

func NewRTSPCommand( rtsp *RTSPServer ) *RTSPCommand {
    return &RTSPCommand{rtsp : rtsp}
}