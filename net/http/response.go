package http

import (
	"net/http"
)

// Response 输出
type Response struct {
	Out         http.ResponseWriter
	Status      int
	ContentType string
}

// WriteHeader 写入Header
func (response *Response) WriteHeader(statusCode int, contentType string) {
	if response.Status == 0 {
		response.Status = statusCode
	}

	if response.ContentType == "" {
		response.ContentType = contentType
	}

	response.Out.Header().Set("Content-Type", response.ContentType)
	response.Out.WriteHeader(response.Status)
}

// NewResponse 创建输出控制器
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{Out: w}
}
