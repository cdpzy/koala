package http

import (
	"reflect"
	"testing"
)

func TestServer(t *testing.T) {
	op := &ServerOptions{
		Addr: ":6021",
	}

	s := NewServer(op)
	s.Router.Register(NewRoute("get", "/find/list/:page", "TAction", "koala/net/http/TestController"))
	s.Controllers.Register(&TestController{}, []*MethodType{{Name: "TAction", Args: []*MethodArg{&MethodArg{Name: "page", Type: reflect.TypeOf((*int)(nil))}}}})
	t.Log(s.Serve())
}

type TestController struct {
	*KoalaController
}

func (req *TestController) TAction(page int) Result {
	return req.RenderText("OKKKK %v", page)
}
