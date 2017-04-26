package http

import (
	"io"
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/websocket"
)

// ServerOptions http服务器配置选项
type ServerOptions struct {
	Addr           string // 监听地址
	ReadTimeout    int    // 读取超时
	WriteTimeout   int    // 写入超时
	Ssl            bool   // ssl 支持
	SslKey         string // ssl 证书key
	SslCert        string // ssl 证书
	MaxRequestSize int64  // http最大请求限制
}

// Server http服务器
type Server struct {
	s              *http.Server
	Ssl            bool               // ssl 支持
	SslKey         string             // ssl 证书key
	SslCert        string             // ssl 证书
	MaxRequestSize int64              // http最大请求限制
	Router         *Router            // 路由器
	Controllers    *ControllerManager // 控制管理器
}

// Serve 启动服务
func (server *Server) Serve() error {
	server.s.Handler = http.HandlerFunc(server.handle)
	if server.Ssl {
		return server.s.ListenAndServeTLS(server.SslCert, server.SslKey)
	}

	listener, err := net.Listen("tcp", server.s.Addr)
	if err != nil {
		return err
	}

	return server.s.Serve(listener)
}

// handle http处理
func (server *Server) handle(w http.ResponseWriter, r *http.Request) {
	if server.MaxRequestSize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, server.MaxRequestSize)
	}

	upgrade := r.Header.Get("Upgrade")
	if upgrade == "websocket" || upgrade == "Websocket" {

		websocket.Handler(func(ws *websocket.Conn) {
			ws.SetDeadline(time.Now().Add(time.Hour * 24))
			r.Method = "WS"
			server.handleRequest(w, r, ws)
		}).ServeHTTP(w, r)

	} else {
		server.handleRequest(w, r, nil)
	}
}

// handleRequest 处理http请求
func (server *Server) handleRequest(w http.ResponseWriter, r *http.Request, ws *websocket.Conn) {
	startTime := time.Now()
	req := NewRequest(r)
	req.Websocket = ws
	resp := NewResponse(w)
	c := NewKoalaController(req, resp)
	c.Controllers = server.Controllers

	server.Router.Apply(c)

	if c.Result() != nil {
		c.Result().Apply(req, resp)
	} else if c.Response().Status != 0 {
		c.Response().Out.WriteHeader(c.Response().Status)
	}

	if w, ok := resp.Out.(io.Closer); ok {
		w.Close()
	}

	log.Infoln("end:", time.Since(startTime))
}

// NewServer 创建Http服务
func NewServer(op *ServerOptions) *Server {
	s := &Server{
		s: &http.Server{
			Addr:         op.Addr,
			ReadTimeout:  time.Duration(op.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(op.WriteTimeout) * time.Second,
		},

		Ssl:            op.Ssl,
		SslKey:         op.SslKey,
		SslCert:        op.SslCert,
		MaxRequestSize: op.MaxRequestSize,
		Router:         NewRouter(),
		Controllers:    NewControllerManager(),
	}
	return s
}
