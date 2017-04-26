package net

import (
	"fmt"
	"koala/net"
	snet "net"

	"github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("koala.net", Loader)
}

func Loader(L *lua.LState) int {
	t := L.NewTypeMetatable("koala.net")
	L.SetGlobal("koala.net", t)
	L.SetField(t, "new", L.NewFunction(newTCPServer))
	L.SetFuncs(t, map[string]lua.LGFunction{"newTCPServer": newTCPServer})
	L.SetField(t, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{"Serve": Serve}))
	L.Push(t)
	return 1
}

func newTCPServer(L *lua.LState) int {
	addr := L.CheckString(1)
	readBufferSize := L.CheckInt(2)
	writeBufferSize := L.CheckInt(3)
	//handle := L.CheckFunction(4)
	//handler := *handle.(net.HandlerFunc)

	op := &net.TCPServerOptions{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		Addr:            addr,
		ClientHandler: func(conn snet.Conn) {
			net.HandleClient(conn, 30, 30, func(b []byte) ([]byte, error) {
				fmt.Println("b:", b)
				return b, nil
			})
		},
	}

	s := net.NewTCPServer(op)
	ud := L.NewUserData()
	ud.Value = s
	L.SetMetatable(ud, L.GetTypeMetatable("koala.net"))
	L.Push(ud)
	return 1
}

func Serve(L *lua.LState) int {
	ud := L.CheckUserData(1)
	v, ok := ud.Value.(*net.TCPServer)
	if !ok {
		return 0
	}

	v.Serve()
	return 1
}
