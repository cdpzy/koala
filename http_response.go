package koala

import (
    "net/http"
)

type HTTPResponse struct {
    Status  int
    ContentType string
    W http.ResponseWriter
}

func NewHTTPResponse( w http.ResponseWriter ) *HTTPResponse {
    resp := new(HTTPResponse)
    resp.W = w

    return resp
}