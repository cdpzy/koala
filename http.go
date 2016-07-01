package koala

import (
    "time"
    "net/http"
    "strings"
    "golang.org/x/net/websocket"
    "log"
)

type HTTPServer struct {
    server    *http.Server
    SSL       bool
    SSLCert   string
    SSLKey    string
    Maxrequestsize int64
    handlerFunc HandlerFunc
}

func (httpServer *HTTPServer) HandlerFunc( handlerFunc HandlerFunc ) {
    httpServer.handlerFunc = handlerFunc
}

func (httpServer *HTTPServer) handle(w http.ResponseWriter, r *http.Request) {
    if httpServer.Maxrequestsize > 0 {
        r.Body = http.MaxBytesReader(w, r.Body, httpServer.Maxrequestsize)
    }

    upgrade := strings.ToLower(r.Header.Get("Upgrade"))
    if upgrade == "websocket" {
        websocket.Handler(func( ws *websocket.Conn ){
            ws.SetDeadline(time.Now().Add(time.Hour * 24))
            r.Method = "WS"
            httpServer.handleInternal(w, r, ws)
        }).ServeHTTP(w, r)
    } else {
         httpServer.handleInternal(w, r, nil)
    }
}


func (httpServer *HTTPServer) handleInternal(w http.ResponseWriter, r *http.Request, ws *websocket.Conn) {
    start    := time.Now()
    request  := NewHTTPRequest(r)
    response := NewHTTPResponse(w)
    request.Websocket = ws
    httpServer.handlerFunc( request, response )
    
    //start.Format("2006/01/02 15:04:05.000")
    log.Printf("%v %v %10v %v %v", 
         "",
         "",
         time.Since(start),
         r.Method,
         r.URL.Path,
    )
}

func (httpServer *HTTPServer) Serve( addr string ) error {
    httpServer.server = &http.Server{
        Addr   :      addr,
        Handler:      http.HandlerFunc(httpServer.handle),
		ReadTimeout:  time.Duration(90) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,
    }

    go func() {
		time.Sleep(100 * time.Millisecond)
		log.Printf("Listening on %s...\n", httpServer.server.Addr)
	}()

    if httpServer.SSL {
        return httpServer.server.ListenAndServeTLS(httpServer.SSLCert, httpServer.SSLKey)
    } else {
        return httpServer.server.ListenAndServe()
    }
    return nil
}

func NewHTTPServer() *HTTPServer {
    return new(HTTPServer)
}