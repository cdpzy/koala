package http

import (
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

// Request 请求
type Request struct {
	*http.Request
	ContentType     string          // 内容类型
	Format          string          // 内容格式
	Websocket       *websocket.Conn // web
	AcceptLanguages AcceptLanguages // 客户端接受的语言
	Locale          string          //
}

// NewRequest 创建Request
func NewRequest(r *http.Request) *Request {
	return &Request{
		Request:         r,
		ContentType:     ResolveContentType(r),
		Format:          ResolveFormat(r),
		AcceptLanguages: ResolveAcceptLanguage(r),
	}
}

// ResolveContentType 获取客户端请内容类型
func ResolveContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return "text/html"
	}

	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

// ResolveFormat 获取客户端请求接受的  MIME
func ResolveFormat(req *http.Request) string {
	accept := req.Header.Get("accept")

	switch {
	case accept == "",
		strings.HasPrefix(accept, "*/*"), // */
		strings.Contains(accept, "application/xhtml"),
		strings.Contains(accept, "text/html"):
		return "html"
	case strings.Contains(accept, "application/json"),
		strings.Contains(accept, "text/javascript"),
		strings.Contains(accept, "application/javascript"):
		return "json"
	case strings.Contains(accept, "application/xml"),
		strings.Contains(accept, "text/xml"):
		return "xml"
	case strings.Contains(accept, "text/plain"):
		return "txt"
	}

	return "html"
}
