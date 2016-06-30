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

import(
    "fmt"
    "github.com/doublemo/koala/msic"
)

type ServerMediaSession struct {
    ssm         bool // 是否为指定播放源
    addr        string 
    description string
    info        string
    miscSDPLines string
    referenceCounter  int // 引用计数
    subsessionCounter int // 子session 计数
    creationTime  msic.Timeval   
    streamName    string
    subSessions   []*ServerSubMediaSession
}

func (serverMediaSession *ServerMediaSession) GenerateSDPDescription() string {
    var (
        sourceFilterLine string
        rangeLine        string
    )

    if serverMediaSession.ssm {
        sourceFilterLine = fmt.Sprintf("a=source-filter: incl IN IP4 * %s\r\n"+
			"a=rtcp-unicast: reflection\r\n", serverMediaSession.addr)
    } else {
        sourceFilterLine = ""
    }

    duration := serverMediaSession.Duration()
    if duration == 0.0 {
        rangeLine = "a=range:npt=0-\r\n"
    } else {
        rangeLine = fmt.Sprintf("a=range:npt=0-%.3f\r\n", duration)
    }

    sdpPrefixFmt := "v=0\r\n" +
		"o=- %d%06d %d IN IP4 %s\r\n" +
		"s=%s\r\n" +
		"i=%s\r\n" +
		"t=0 0\r\n" +
		"a=tool:%s%s\r\n" +
		"a=type:broadcast\r\n" +
		"a=control:*\r\n" +
		"%s" +
		"%s" +
		"a=x-qt-text-nam:%s\r\n" +
		"a=x-qt-text-inf:%s\r\n" +
		"%s"

	sdp := fmt.Sprintf(sdpPrefixFmt,
		serverMediaSession.creationTime.Tv_sec,
		serverMediaSession.creationTime.Tv_usec,
		1,
		serverMediaSession.addr,
		serverMediaSession.description,
		serverMediaSession.info,
		"Koala Media Server V", "1.0",
		sourceFilterLine,
		rangeLine,
		serverMediaSession.description,
		serverMediaSession.info,
		serverMediaSession.miscSDPLines)

    for i := 0; i < serverMediaSession.subsessionCounter; i++ {
		sdpLines := serverMediaSession.subSessions[i].SDPLines()
		sdp += sdpLines
	}

    return sdp
}

func (serverMediaSession *ServerMediaSession) AddSubSession( session  *ServerSubMediaSession) error {
    serverMediaSession.subSessions = append(serverMediaSession.subSessions, session)
    serverMediaSession.subsessionCounter++
    session.IncrTrackNumber()
    return nil
}

func (serverMediaSession *ServerMediaSession) Duration() float32 {
    return 0.0
}

func (serverMediaSession *ServerMediaSession) GetStreamName() string {
    return serverMediaSession.streamName
}

func NewServerMediaSession( description, streamName string ) *ServerMediaSession {
    serverMediaSession := new(ServerMediaSession)
    serverMediaSession.streamName   = streamName
    serverMediaSession.description  = description + ", streamed by the Koala Media Server"
    serverMediaSession.info         = streamName
    serverMediaSession.subSessions  = make([]*ServerSubMediaSession, 0)
    serverMediaSession.addr,_       = msic.OurIPAddress()

    msic.GetTimeOfDay(&serverMediaSession.creationTime)
    return serverMediaSession
}