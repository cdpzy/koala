package koala

import (
    "fmt"
  //  "bytes"
    "net/http"
)

type HTTPResponse struct {
    BaseResponse
    ContentType string
    W http.ResponseWriter
}

func (httpResponse *HTTPResponse) NotFound() {}
func (httpResponse *HTTPResponse) NotSupported( allowedMethod string ) {}
func (httpResponse *HTTPResponse) BadRequest( allowedMethod string ) {}

func (httpResponse *HTTPResponse) GetHeader() http.Header {
    return httpResponse.W.Header()   
}

func (httpResponse *HTTPResponse) Write( b string ) error {
    fmt.Println(httpResponse.GetHeader())
    httpResponse.W.Write([]byte(b))
    return nil
} 


func NewHTTPResponse( w http.ResponseWriter ) *HTTPResponse {
    resp := new(HTTPResponse)
    resp.W = w
    resp.Header = make(map[string][]string)
    resp.Proto  = "HTTP"
    resp.ProtoMajor = 1
    resp.ProtoMinor = 0
    resp.Status = "OK"
    resp.StatusCode = RESPONSE_STATUS_CODE_OK
    return resp
}