package mediasession

import (
    "fmt"
    "github.com/doublemo/koala/helper"
)

const (
    BANK_SIZE uint = 150000
    NO_MORE_BUFFERED_INPUT = 1
)

type StreamParser struct {
    curBankNum            uint
    curParserIndex        uint
    saveParserIndex       uint
    totNumValidBytes      uint
    savedParserIndex      uint
    remainingUnparsedBits uint
    savedRemainingUnparsedBits uint
    haveSeenEOF           bool
    bank                  []byte
    curBank               []byte
    clientContinueFunc     interface{}
    clientOnInputCloseFunc interface{}
    lastSeenPresentationTime helper.Time
    inputSource            FramedSource
}

func (streamParser *StreamParser) SetInputSource( source FramedSource ) {
    streamParser.inputSource = source
}

func (streamParser *StreamParser) BankSize() uint {
    return BANK_SIZE
}

func (streamParser *StreamParser) Get4Bytes() uint{
    rs := streamParser.test4Bytes()
    streamParser.curParserIndex +=  4
    streamParser.remainingUnparsedBits = 0
    return rs
}

func (streamParser *StreamParser) Get2Bytes() uint {
    streamParser.ensureValidBytes(2)

    ptr := streamParser.nextToParse()
    rs  := (ptr[0] << 8) | ptr[1]

    streamParser.curParserIndex += 2
    streamParser.remainingUnparsedBits = 0
    return uint(rs)
}

func (streamParser *StreamParser) Get1Bytes() uint {
    streamParser.ensureValidBytes(1)
    streamParser.curParserIndex ++
    return uint(streamParser.CurBank()[streamParser.curParserIndex])
}


func (streamParser *StreamParser) restoreSavedParserState() {
    streamParser.curParserIndex = streamParser.savedParserIndex
    streamParser.remainingUnparsedBits = streamParser.savedRemainingUnparsedBits
}

func (streamParser *StreamParser) testBytes( to []byte, numBytes uint ) {
    streamParser.ensureValidBytes(numBytes)
    to = streamParser.nextToParse()[:numBytes]
}

func (streamParser *StreamParser) skipBytes( numBytes uint ) {
    streamParser.ensureValidBytes(numBytes)
    streamParser.curParserIndex += numBytes
}

func (streamParser *StreamParser) test4Bytes() uint {
    streamParser.ensureValidBytes(4)

    ptr := streamParser.nextToParse()
    return uint((ptr[0] << 24) | (ptr[1] << 16) | (ptr[2] << 8) | ptr[3])
}

func (streamParser *StreamParser) nextToParse() []byte{
    return streamParser.CurBank()[streamParser.curParserIndex:]
}

func (streamParser *StreamParser) curOffset() uint {
    return streamParser.curParserIndex
}

func (streamParser *StreamParser) HaveSeenEOF() bool {
    return streamParser.haveSeenEOF
}

func (streamParser *StreamParser) saveParserState() {
    streamParser.savedParserIndex = streamParser.curParserIndex
    streamParser.savedRemainingUnparsedBits = streamParser.remainingUnparsedBits
}

func (streamParser *StreamParser) TotNumValidBytes() uint {
    return streamParser.totNumValidBytes
}

func (streamParser *StreamParser) ensureValidBytes( numBytesNeeded uint ) {
    if streamParser.curParserIndex + numBytesNeeded <= streamParser.totNumValidBytes {
        return
    }

    streamParser.ensureValidBytes1(numBytesNeeded)
}

func (streamParser *StreamParser) ensureValidBytes1( numBytesNeeded uint ) uint {
    maxFramedSize := streamParser.inputSource.MaxFrameSize()
    if maxFramedSize > numBytesNeeded {
        numBytesNeeded = maxFramedSize
    }

    if streamParser.curParserIndex + numBytesNeeded > BANK_SIZE {
        numBytesToSave := streamParser.totNumValidBytes + streamParser.savedParserIndex

        streamParser.curBankNum = (streamParser.curBankNum + 1) % 2
        streamParser.curBank    = streamParser.bank[streamParser.curBankNum:]
        streamParser.curBank    = streamParser.curBank[streamParser.saveParserIndex : streamParser.saveParserIndex + numBytesToSave]

        streamParser.curParserIndex   -= streamParser.savedParserIndex
        streamParser.savedParserIndex  = 0
        streamParser.totNumValidBytes  = numBytesToSave
    }

    if streamParser.curParserIndex + numBytesNeeded > BANK_SIZE {
        panic("StreamParser Internal error")
    }

    maxNumBytesToRead := BANK_SIZE - streamParser.totNumValidBytes
    streamParser.inputSource.Next(streamParser.CurBank(), maxNumBytesToRead, streamParser.afterGettingBytes, streamParser.onInputClosure)
    return NO_MORE_BUFFERED_INPUT
}

func (streamParser *StreamParser) afterGettingBytes(numBytesRead uint, presentationTime helper.Time) {
    if streamParser.totNumValidBytes + numBytesRead > BANK_SIZE {
        fmt.Printf("StreamParser::afterGettingBytes() warning: read %d bytes; expected no more than %d\n", numBytesRead, BANK_SIZE-streamParser.totNumValidBytes)
    }

    streamParser.lastSeenPresentationTime = presentationTime
    streamParser.restoreSavedParserState()
    streamParser.clientContinueFunc.(func())()
}

func (streamParser *StreamParser) onInputClosure() {
    if !streamParser.haveSeenEOF {
        streamParser.haveSeenEOF = true
        streamParser.afterGettingBytes(0, streamParser.lastSeenPresentationTime)
    } else {
        streamParser.haveSeenEOF = false
        if streamParser.clientOnInputCloseFunc != nil {
            streamParser.clientOnInputCloseFunc.(func())()
        }
    }
}

func (streamParser *StreamParser) CurBank() []byte {
    return streamParser.curBank
}