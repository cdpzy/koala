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

package session

import (
    "fmt"
    "os"
    "github.com/doublemo/koala/codec"
    "github.com/doublemo/koala/msic"
)

type ServerSubMediaSession struct {
    trackNumber int
    trackId    string
    streamName string
    initialPortNum int
    cname      string
    sdpLines   string
    fileSize   int64
    portNumForSDP int
}

func (serverSubMediaSession *ServerSubMediaSession) SDPLines() string {
    if serverSubMediaSession.sdpLines != "" {
        return serverSubMediaSession.sdpLines
    }

    inputSource    := serverSubMediaSession.createNewStreamSource()
    rtpPayloadType := 96 + serverSubMediaSession.GetTrackNumber() - 1
    var dummyAddr string
	dummyGroupSock := msic.NewGroupSock(dummyAddr, 0)
	dummyRTPSink := serverSubMediaSession.createNewRTPSink(dummyGroupSock, uint(rtpPayloadType))
    fmt.Println("inputSource:", inputSource, rtpPayloadType, dummyRTPSink)
    serverSubMediaSession.setSDPLinesFromRTPSink(dummyRTPSink, inputSource, 500)
    return serverSubMediaSession.sdpLines
}

func (serverSubMediaSession *ServerSubMediaSession) setSDPLinesFromRTPSink( rtpSink *codec.H264VideoRTPSink, inputSource *codec.H264VideoStreamParser, estBitrate int){
    mediaType      := rtpSink.SdpMediaType()
    rtpmapLine     := rtpSink.RtpmapLine()
    rtpPayloadType := rtpSink.RtpPayloadType()

    rangeLine      := serverSubMediaSession.rangeSDPLine()
    auxSDPLine     := serverSubMediaSession.getAuxSDPLine(rtpSink)

    ipAddr := "0.0.0.0"
    sdpFmt := "m=%s %d RTP/AVP %d\r\n" +
		"c=IN IP4 %s\r\n" +
		"b=AS:%d\r\n" +
		"%s" +
		"%s" +
		"%s" +
		"a=control:%s\r\n"

	serverSubMediaSession.sdpLines = fmt.Sprintf(sdpFmt,
		mediaType,
		serverSubMediaSession.portNumForSDP,
		rtpPayloadType,
		ipAddr,
		estBitrate,
		rtpmapLine,
		rangeLine,
		auxSDPLine,
		serverSubMediaSession.GetTrackId())
}




func (serverSubMediaSession *ServerSubMediaSession) GetTrackId() string {
    if serverSubMediaSession.trackId == "" {
        serverSubMediaSession.trackId = fmt.Sprintf("track%d", serverSubMediaSession.trackNumber)
    }

    return serverSubMediaSession.trackId
}

func (serverSubMediaSession *ServerSubMediaSession) IncrTrackNumber() {
    serverSubMediaSession.trackNumber++
}

func (serverSubMediaSession *ServerSubMediaSession) GetTrackNumber() int {
    return serverSubMediaSession.trackNumber
}

func (serverSubMediaSession *ServerSubMediaSession) rangeSDPLine() string {
	return "a=range:npt=0-\r\n"
}

func (serverSubMediaSession *ServerSubMediaSession) getAuxSDPLine(rtpSink *codec.H264VideoRTPSink) string {
	return rtpSink.AuxSDPLine()
}


func (serverSubMediaSession *ServerSubMediaSession) createNewStreamSource() *codec.H264VideoStreamParser{
    fileSource := codec.NewByteStreamFileSource(serverSubMediaSession.streamName)
    if fileSource == nil {
        return nil
    }

    serverSubMediaSession.fileSize = fileSource.FileSize()

    return codec.NewH264VideoStreamParser(fileSource)
}

func (serverSubMediaSession *ServerSubMediaSession) createNewRTPSink( dummyGroupSock *msic.GroupSock, rtpPayloadType uint ) *codec.H264VideoRTPSink {
    return codec.NewH264VideoRTPSink(dummyGroupSock, rtpPayloadType)
}

func NewServerSubMediaSession( streamName string ) *ServerSubMediaSession {
    serverSubMediaSession := new(ServerSubMediaSession)
    serverSubMediaSession.streamName     = streamName
    serverSubMediaSession.initialPortNum = 6970
    serverSubMediaSession.cname, _       = os.Hostname()
    serverSubMediaSession.sdpLines       = ""
    return serverSubMediaSession
}