package koala

import (
    "log"
    "strings"
    "time"
)

const AllowedMethod =  "OPTIONS, DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, GET_PARAMETER, SET_PARAMETER"

type HandleMethod struct{
    r Request
    w Response
}

// http 
func (handleMethod *HandleMethod) GET() {
    header := handleMethod.w.GetHeader()
    csep   := handleMethod.r.GetHeader().Get("CSeq")
    header.Set("CSeq", csep)
    header.Set("Date", time.Now().String())
    header.Set("Cache-Control", "no-cache")
    header.Set("Pragma", "no-cache")
    header.Set("Content-Type", "application/x-rtsp-tunnelled")
    handleMethod.w.Write("")
}

func (handleMethod *HandleMethod) POST() {
    log.Println(handleMethod.w.GetHeader())
}

func (handleMethod *HandleMethod) OPTIONS() {
    csep   := handleMethod.r.GetHeader().Get("CSeq")
    header := handleMethod.w.GetHeader()
    header.Set("CSeq", csep)
    header.Set("Date", time.Now().String())
    header.Set("Public", AllowedMethod)
    handleMethod.w.Write("")
}

func (handleMethod *HandleMethod) DESCRIBE() {
    path := strings.Trim(handleMethod.r.GetURL().Path, "/")
    log.Println("path:", path)
}


func NewHandleMethod( r Request, w Response ) *HandleMethod {
    handle := new(HandleMethod)
    handle.r = r
    handle.w = w
    return handle
}


