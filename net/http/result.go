package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Result 内容格式输出接口
type Result interface {
	Apply(req *Request, resp *Response)
}

// ErrorResult 错误
type ErrorResult struct {
	Error error
}

// RenderTextResult 文本信息
type RenderTextResult struct {
	text string
}

// Apply 错误输出
func (errorResult ErrorResult) Apply(req *Request, resp *Response) {
	//format := req.Format
	status := resp.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}

	resp.WriteHeader(status, "text/plain; charset=utf-8")
	resp.Out.Write([]byte(errorResult.Error.Error()))
}

// Apply 文本输出
func (renderTextResult RenderTextResult) Apply(req *Request, resp *Response) {
	resp.WriteHeader(http.StatusOK, "text/plain; charset=utf-8")
	resp.Out.Write([]byte(renderTextResult.text))
}

type ContentDisposition string

var (
	Attachment ContentDisposition = "attachment"
	Inline     ContentDisposition = "inline"
)

type BinaryResult struct {
	Reader   io.Reader
	Name     string
	Length   int64
	Delivery ContentDisposition
	ModTime  time.Time
}

func (r *BinaryResult) Apply(req *Request, resp *Response) {
	disposition := string(r.Delivery)
	if r.Name != "" {
		disposition += fmt.Sprintf(`; filename="%s"`, r.Name)
	}
	resp.Out.Header().Set("Content-Disposition", disposition)

	// If we have a ReadSeeker, delegate to http.ServeContent
	if rs, ok := r.Reader.(io.ReadSeeker); ok {
		// http.ServeContent doesn't know about response.ContentType, so we set the respective header.
		if resp.ContentType != "" {
			resp.Out.Header().Set("Content-Type", resp.ContentType)
		} else {
			contentType := ContentTypeByFilename(r.Name)
			resp.Out.Header().Set("Content-Type", contentType)
		}
		http.ServeContent(resp.Out, req.Request, r.Name, r.ModTime, rs)
	} else {
		// Else, do a simple io.Copy.
		if r.Length != -1 {
			resp.Out.Header().Set("Content-Length", strconv.FormatInt(r.Length, 10))
		}
		resp.WriteHeader(http.StatusOK, ContentTypeByFilename(r.Name))
		io.Copy(resp.Out, r.Reader)
	}

	// Close the Reader if we can
	if v, ok := r.Reader.(io.Closer); ok {
		v.Close()
	}
}
