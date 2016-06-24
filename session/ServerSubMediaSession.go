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
)

type ServerSubMediaSession struct {
    trackNumber int
    trackId    string
    streamName string
    initialPortNum int
    cname      string
    sdpLines   string
}

func (serverSubMediaSession *ServerSubMediaSession) SDPLines() string {
    return ""
}

func (serverSubMediaSession *ServerSubMediaSession) GetTrackId() string {
    if serverSubMediaSession.trackId == "" {
        serverSubMediaSession.trackId = fmt.Sprintf("track%d", serverSubMediaSession.trackId)
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

func (serverSubMediaSession *ServerSubMediaSession) createNewStreamSource(){
    
}

func NewServerSubMediaSession( streamName string ) *ServerSubMediaSession {
    serverSubMediaSession := new(ServerSubMediaSession)
    serverSubMediaSession.streamName     = streamName
    serverSubMediaSession.initialPortNum = 6970
    serverSubMediaSession.cname, _       = os.Hostname()
    serverSubMediaSession.sdpLines       = ""
    return serverSubMediaSession
}