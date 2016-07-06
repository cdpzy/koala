# koala
Koala Media Server  - 演示阶段

```go
package main

import(
    "github.com/doublemo/koala"
    "log"
    "reflect"
)

func main() {
    // s := koala.NewRTSPServer()
    // log.Println(s.GetRequest().String())
    // log.Println(s.Serve(koala.SERVER_NET_PROTO_TCP, ":554"))
    rtsp := koala.NewRTSPServer()
    rtsp.HandlerFunc(Handler2)
    go rtsp.Serve(":554")


    http := koala.NewHTTPServer()
    http.HandlerFunc(Handler)
    log.Println(http.Serve(":9106"))
}

func Handler( req koala.Request, resp koala.Response ) {
    method := koala.NewHandleMethod( req, resp )
    v := reflect.ValueOf(method).MethodByName(req.GetMethod())
    if !v.IsValid() {
        resp.NotSupported( koala.AllowedMethod )
        return
    }

    v.Call([]reflect.Value{})
}
   

func Handler2( req koala.Request, resp koala.Response ) {
    method := koala.NewHandleMethod( req, resp )
    ch := req.GetInputChan()
    for {
        select {
            case b,ok := <- ch:
                 if !ok || b != nil {
                     break
                 }

                 v := reflect.ValueOf(method).MethodByName(req.GetMethod())
                 if !v.IsValid() {
                     resp.NotSupported( koala.AllowedMethod )
                    return
                 }

                v.Call([]reflect.Value{})
        }
    }
}
```