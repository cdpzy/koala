package koala

import (
    "time"
    "net"
    "log"
    "strings"
    "net/http"
    "net/url"
    "io"
    "io/ioutil"
    "bytes"
    "bufio"
    "errors"
    "fmt"
    "strconv"
)

type RTSPRequest struct {
    BaseRequest
    Method      string
    URL         *url.URL
    Proto        string
    ProtoMajor    int
    ProtoMinor    int
    Header        http.Header
    ContentLength int
    Body          io.ReadCloser     
}

func (rtspRequest *RTSPRequest) GetMethod() string {
    return rtspRequest.Method
}

func (rtspRequest *RTSPRequest) GetHeader() http.Header {
	return rtspRequest.Header
}

func (rtspRequest *RTSPRequest) GetBody() io.ReadCloser {
	return rtspRequest.Body
}

func (rtspRequest *RTSPRequest) String() string {
    s := fmt.Sprintf("%s %s %s/%d.%d\r\n", rtspRequest.Method, rtspRequest.URL, rtspRequest.Proto, rtspRequest.ProtoMajor, rtspRequest.ProtoMinor)
    for k, v := range rtspRequest.Header {
        for _, v := range v {
            s += fmt.Sprintf("%s: %s\r\n", k, v)
        }
    }
    s += "\r\n"
    if rtspRequest.Body != nil {
        str, _ := ioutil.ReadAll(rtspRequest.Body)
        s += string(str)
    }
    return s
}

func (rtspRequest *RTSPRequest) Recv() {
     ip := net.ParseIP(strings.Split( rtspRequest.RemoteAddr.String(), ":")[0])
     log.Printf("new connected from:%v\n", ip)

     for {
         p      := make([]byte, 4096)
         n, err := rtspRequest.Socket.Read(p)
         if err != nil {
             log.Printf("error receiving, bytes:%d reason:%v\n", n, err)
             break
         }

         select {
             case rtspRequest.in <- rtspRequest.parseRequest(bytes.NewBuffer(p)):
             case <-time.After(30 * time.Second):
                  log.Printf("server busy or listen closed.")
         }
     }

     log.Printf("Client shutdown:%v", ip)
}

func (rtspRequest *RTSPRequest) parseRequest( r io.Reader ) error {
    rtspRequest.Header = make(map[string][]string)
    b       := bufio.NewReader( r )
    s, err  := b.ReadString('\n')

    if err != nil {
        return err
    }

    parts  := strings.SplitN(s, " ", 3)
    if len( parts )  != 3 {
        return errors.New(fmt.Sprintf("Invalid Request:%v", s))
    }

    rtspRequest.Method   = parts[0]
    rtspRequest.URL, err = url.Parse(parts[1])
    if err != nil {
        return err
    }

    rtspRequest.Proto, rtspRequest.ProtoMajor, rtspRequest.ProtoMinor, err = rtspRequest.parseVersion(parts[2])
    if err != nil {
        return err
    }

    // parse header
    for {
        if s, err = b.ReadString('\n'); err != nil {
            return err
        } else if s = strings.TrimRight(s, "\r\n"); s == "" {
            break
        }

        head := strings.SplitN(s, ":", 2)
        if len(head) != 2 {
            continue
        }

        rtspRequest.Header.Add(strings.TrimSpace(head[0]), strings.TrimSpace(head[1]))
    }

    rtspRequest.ContentLength, _ = strconv.Atoi(rtspRequest.Header.Get("Content-Length"))
    rtspRequest.Body             = RTSPRequestCloser{b, r}
    return nil
}

func (rtspRequest *RTSPRequest) parseVersion(s string) (proto string, major int, minor int, err error) {
    s = strings.TrimRight(s, "\r\n")
    parts := strings.SplitN(s, "/", 2)
    proto = parts[0]
    parts = strings.SplitN(parts[1], ".", 2)
    if major, err = strconv.Atoi(parts[0]); err != nil {
        return
    }
    if minor, err = strconv.Atoi(parts[1]); err != nil {
        return
    }
    return
}

func NewRTSPRequest( socket net.Conn ) *RTSPRequest {
    req := new(RTSPRequest)
    req.Socket = socket
    req.RemoteAddr = socket.RemoteAddr()
    req.LocalAddr  = socket.LocalAddr()
    req.in     = make(chan error, 0)
    return req
}


type RTSPRequestCloser struct {
    *bufio.Reader
    r io.Reader
}

func (rtspRequestCloser RTSPRequestCloser) Close() error{
    defer func(){
        rtspRequestCloser.Reader = nil
        rtspRequestCloser.r      = nil
    }()

    if rtspRequestCloser.Reader == nil {
        return nil
    }

    if r, ok := rtspRequestCloser.r.(io.ReadCloser); ok {
        return r.Close()
    }

    return nil
}