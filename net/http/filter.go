package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/websocket"
)

// Filter 过滤操作
type Filter func(route *Router, c *KoalaController, filterChain []Filter)

// PanicFilter 默认Panic 错误处理
func PanicFilter(r *Router, c *KoalaController, filterChain []Filter) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Runtime Error:", err)
			log.Errorf("%s\n", debug.Stack())
			c.SetResult(c.RenderError(500, fmt.Errorf("%v", err)))
		}
	}()

	filterChain[0](r, c, filterChain[1:])
}

// RouterFilter 路由
func RouterFilter(r *Router, c *KoalaController, filterChain []Filter) {
	req := c.Request()
	if method := req.Header.Get("X-HTTP-Method-Override"); method != "" && req.Method == "POST" {
		req.Method = method
	}

	route := r.Find(req.Method, req.URL.Path)
	if route == nil {
		// todo 404 error
		c.SetResult(c.RenderError(http.StatusNotFound, fmt.Errorf("Match: %s %s [Failed]", req.Method, req.URL.Path)))
		return
	}

	if err := c.SetAction(route.Controller, route.Action); err != nil {
		c.SetResult(c.RenderError(http.StatusInternalServerError, fmt.Errorf("500 SERVER ERROR: %s", err)))
		return
	}

	c.Params.Values = make(url.Values)
	c.Params.Route = route.Params

	filterChain[0](r, c, filterChain[1:])
}

// ParamsFilter 参数处理
func ParamsFilter(r *Router, c *KoalaController, filterChain []Filter) {
	c.Params.Parse(c.Request())
	c.Params.Values = c.Params.calcValues()

	filterChain[0](r, c, filterChain[1:])
}

// InvokerFilter 调用Controller
func InvokerFilter(r *Router, c *KoalaController, filterChain []Filter) {
	var (
		methodArgs   []reflect.Value
		boundArg     reflect.Value
		resultValue  reflect.Value
		resultValues []reflect.Value
	)

	methodValue := reflect.ValueOf(c.AppController).MethodByName(c.MethodType.Name)
	websocketType := reflect.TypeOf((*websocket.Conn)(nil))
	for _, arg := range c.MethodType.Args {
		if arg.Type == websocketType {
			boundArg = reflect.ValueOf(c.Request().Websocket)
		} else {
			boundArg = Bind(c.Params, arg.Name, arg.Type)
			if closer, ok := boundArg.Interface().(io.Closer); ok {
				defer closer.Close()
			}
		}
		methodArgs = append(methodArgs, boundArg)
	}

	if methodValue.Type().IsVariadic() {
		resultValues = methodValue.CallSlice(methodArgs)
	} else {
		resultValues = methodValue.Call(methodArgs)
	}

	if len(resultValues) < 1 {
		return
	}

	resultValue = resultValues[0]
	if resultValue.Kind() == reflect.Interface && !resultValue.IsNil() {
		c.SetResult(resultValue.Interface().(Result))
	}
}
