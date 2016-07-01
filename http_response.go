package koala

import (
    "net/http"
)

type HTTPResponse struct {
    BaseResponse
    Status  int
    ContentType string
    W http.ResponseWriter
}

func (httpResponse *HTTPResponse) Write( b []byte ) error {
    return nil
} 

func NewHTTPResponse( w http.ResponseWriter ) *HTTPResponse {
    resp := new(HTTPResponse)
    resp.W = w

    return resp
}