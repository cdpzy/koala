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

package codec

import (
    "fmt"
    "github.com/doublemo/koala/msic"
)

var TRANSPORT_SYNC_BYTE byte = 0x47
var TRANSPORT_PACKET_SIZE uint = 188
var TRANSPORT_PACKETS_PER_NETWORK_PACKET uint = 7


type H264VideoStreamParser struct {
    afterGettingFunc        interface{}
	onCloseFunc             interface{}
	buffTo                  []byte
	maxSize                 uint
	frameSize               uint
	numTruncatedBytes       uint
	durationInMicroseconds  uint
	isCurrentlyAwaitingData bool
	presentationTime        msic.Timeval
    source                  *ByteStreamFileSource
}


func (h264VideoStreamParser *H264VideoStreamParser) getNextFrame(buffTo []byte, maxSize uint, afterGettingFunc interface{}, onCloseFunc interface{}) {
    if h264VideoStreamParser.isCurrentlyAwaitingData {
		panic("FramedSource::getNextFrame(): attempting to read more than once at the same time!")
	}

    h264VideoStreamParser.buffTo = buffTo
	h264VideoStreamParser.maxSize = maxSize
	h264VideoStreamParser.onCloseFunc = onCloseFunc
	h264VideoStreamParser.afterGettingFunc = afterGettingFunc
	h264VideoStreamParser.isCurrentlyAwaitingData = true    

    h264VideoStreamParser.source.NextFrame()
}


func (h264VideoStreamParser *H264VideoStreamParser) afterGetting() {
	fmt.Println("FramedSource::afterGetting")
	h264VideoStreamParser.isCurrentlyAwaitingData = false

	if h264VideoStreamParser.afterGettingFunc != nil {
		h264VideoStreamParser.afterGettingFunc.(func(frameSize, durationInMicroseconds uint, presentationTime msic.Timeval))(h264VideoStreamParser.frameSize, h264VideoStreamParser.durationInMicroseconds, h264VideoStreamParser.presentationTime)
	}
}

func (h264VideoStreamParser *H264VideoStreamParser) handleClosure() {
	h264VideoStreamParser.isCurrentlyAwaitingData = false

	if h264VideoStreamParser.onCloseFunc != nil {
		h264VideoStreamParser.onCloseFunc.(func())()
	}
}

func (h264VideoStreamParser *H264VideoStreamParser) stopGettingFrames() {
	h264VideoStreamParser.isCurrentlyAwaitingData = false
}

func (h264VideoStreamParser *H264VideoStreamParser) maxFrameSize() uint {
	return 0
}


func NewH264VideoStreamParser( source *ByteStreamFileSource ) *H264VideoStreamParser{
    h264VideoStreamParser := new(H264VideoStreamParser)
    h264VideoStreamParser.source = source
    return h264VideoStreamParser
}