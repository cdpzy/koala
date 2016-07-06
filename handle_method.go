package koala

import (
    "log"
    "strings"
    "time"
    "strconv"
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
    header := handleMethod.w.GetHeader()
    path   := strings.Trim(handleMethod.r.GetURL().Path, "/")
    csep   := header.Get("CSeq")

    session, err   := ServerMediaSessionManager.Create( path )
    if err != nil {
        handleMethod.w.BadRequest(AllowedMethod)
        return
    }

    sdpDescription     := session.GenerateSDPDescription()
    sdpDescriptionSize := len(sdpDescription)
    if sdpDescriptionSize < 1 {
        handleMethod.w.NotFound()
        return
    }

    header.Set("CSeq", csep)
    header.Set("Date", time.Now().String())
    header.Set("Content-Base", handleMethod.r.GetURL().String())
    header.Set("Content-Type", "application/sdp")
    header.Set("Content-Length", strconv.Itoa(sdpDescriptionSize))
    handleMethod.w.Write(sdpDescription)
}


func NewHandleMethod( r Request, w Response ) *HandleMethod {
    handle := new(HandleMethod)
    handle.r = r
    handle.w = w
    return handle
}

