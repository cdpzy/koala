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
    "os"
    "log"
    "github.com/doublemo/koala/msic"
)

type ByteStreamFileSource struct {
    presentationTime      msic.Timeval
	fileSize              int64
	numBytesToStream      int64
	lastPlayTime          uint
	playTimePerFrame      uint
	preferredFrameSize    uint
	haveStartedReading    bool
	limitNumBytesToStream bool
    fid *os.File
    buffTo               []byte  
    isCurrentlyAwaitingData bool
    maxSize                 uint
	frameSize               uint
	numTruncatedBytes       uint
	durationInMicroseconds  uint     
}


func (byteStreamFileSource *ByteStreamFileSource) NextFrame() {
    if byteStreamFileSource.limitNumBytesToStream && byteStreamFileSource.numBytesToStream == 0 {
        byteStreamFileSource.handleClosure()
        return
    }

    byteStreamFileSource.readFromFile()
}

func (byteStreamFileSource *ByteStreamFileSource) handleClosure() {
    defer byteStreamFileSource.fid.Close()

    byteStreamFileSource.isCurrentlyAwaitingData = false
    byteStreamFileSource.haveStartedReading = false
}

func (byteStreamFileSource *ByteStreamFileSource) readFromFile() error {
    _, err := byteStreamFileSource.fid.Read(byteStreamFileSource.buffTo)
    if err != nil {
        log.Println("readFromFile:", err)
        return err
    }

    if byteStreamFileSource.playTimePerFrame > 0 && byteStreamFileSource.preferredFrameSize > 0 {
        if byteStreamFileSource.presentationTime.Tv_sec == 0 && byteStreamFileSource.presentationTime.Tv_usec == 0 {
            msic.GetTimeOfDay( &byteStreamFileSource.presentationTime  )
        } else {
            uSeconds := byteStreamFileSource.presentationTime.Tv_usec + int64(byteStreamFileSource.lastPlayTime)
            byteStreamFileSource.presentationTime.Tv_sec += uSeconds / 1000000
			byteStreamFileSource.presentationTime.Tv_usec = uSeconds % 1000000
        }


        byteStreamFileSource.lastPlayTime = (byteStreamFileSource.playTimePerFrame * byteStreamFileSource.frameSize) / byteStreamFileSource.preferredFrameSize
        byteStreamFileSource.durationInMicroseconds = byteStreamFileSource.lastPlayTime
    } else {
        msic.GetTimeOfDay(&byteStreamFileSource.presentationTime)
    }

    byteStreamFileSource.afterGetting()
    return nil
}

func (byteStreamFileSource *ByteStreamFileSource) afterGetting() {
    byteStreamFileSource.isCurrentlyAwaitingData = false
    
}

func (byteStreamFileSource *ByteStreamFileSource) FileSize() int64 {
    return byteStreamFileSource.fileSize
}



func NewByteStreamFileSource( fileName string ) *ByteStreamFileSource {
    fid, err := os.Open(fileName)
    if err != nil {
        log.Println("fileOpen error:", fileName, err)
        return  nil
    }

    source := new(ByteStreamFileSource)
    source.fid = fid
    source.buffTo = make([]byte, 20000)
    stat, _ := fid.Stat()
    source.fileSize  = stat.Size()

    return source
}