package koala

import (
    "net"
    "net/http"
    "io"
    "bytes"
    "fmt"
    "errors"
    "strings"
    "time"
)

type RTSPResponse struct {
    BaseResponse
    Socket net.Conn
    out chan []byte
}

func (rtspResponse *RTSPResponse) BadRequest( allowedMethod string ) {
    rtspResponse.Reset()
    rtspResponse.Status     = "Bad Request"
    rtspResponse.StatusCode =  RESPONSE_STATUS_CODE_BADREQUEST
    rtspResponse.Header.Set("Date", time.Now().String())
    rtspResponse.Header.Set("Allow", allowedMethod)
    rtspResponse.Write("")
}

func (rtspResponse *RTSPResponse) NotFound() {
    rtspResponse.Reset()
    rtspResponse.Status     = "Not Found"
    rtspResponse.StatusCode =  RESPONSE_STATUS_CODE_NOTFOUND
    rtspResponse.Header.Set("Date", time.Now().String())
    rtspResponse.Write("")
}

func (rtspResponse *RTSPResponse) NotSupported( allowedMethod string ) {
    rtspResponse.Reset()
    rtspResponse.Status     = "Method Not Allowed"
    rtspResponse.StatusCode =  RESPONSE_STATUS_CODE_NOTALLOWED
    rtspResponse.Header.Set("Date", time.Now().String())
    rtspResponse.Header.Set("Allow", allowedMethod)
    rtspResponse.Write("")
}

func (rtspResponse *RTSPResponse) GetHeader() http.Header {
    return rtspResponse.Header
}

func (rtspResponse *RTSPResponse) Recv() {
    for {
        select {
            case raw, ok := <-rtspResponse.out:
                 if !ok {
                     break
                 }

                fmt.Println(string(raw))
                rtspResponse.Socket.Write(raw)
        }
    }
}

func (rtspResponse *RTSPResponse) Write( b string ) error {
    s      := fmt.Sprintf("%s/%d.%d %d %s\r\n", rtspResponse.Proto, rtspResponse.ProtoMajor, rtspResponse.ProtoMinor, rtspResponse.StatusCode, rtspResponse.Status)
    w      := bytes.NewBuffer([]byte(s))
    err    := rtspResponse.writeHeader( w )
    if err != nil {
        return err
    }

    _, err  = w.WriteString(b)
    if err != nil {
        return err
    }

    rtspResponse.out <- w.Bytes()
    rtspResponse.Reset()
    return nil
} 

func (rtspResponse *RTSPResponse) Reset(){
    rtspResponse.Header = make(map[string][]string)
}

func (rtspResponse *RTSPResponse) writeHeader( w io.Writer ) error {
    if rtspResponse.Header == nil {
        return errors.New("Header is nil")
    }

    for k, v := range rtspResponse.Header {
        _, err := io.WriteString(w, fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ":")))
        if err != nil {
            return err
        }
    }

    if len(rtspResponse.Header) > 0 {
        io.WriteString(w, "\r\n")
    }

    return nil
}

func NewRTSPResponse( socket net.Conn ) *RTSPResponse {
    resp := new(RTSPResponse)
    resp.Socket = socket
    resp.out    = make(chan []byte, 0)
    resp.Header = make(map[string][]string)
    resp.Proto  = "RTSP"
    resp.ProtoMajor = 1
    resp.ProtoMinor = 0
    resp.Status = "OK"
    resp.StatusCode = RESPONSE_STATUS_CODE_OK
    return resp
}