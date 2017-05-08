package main

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/doublemo/koala/lua"
	knet "github.com/doublemo/koala/net"
	glua "github.com/yuin/gopher-lua"
	cli "gopkg.in/urfave/cli.v2"
)

var tcpService *knet.TCPServer

// tcpServe TCP服务器
func tcpServe(c *cli.Context) {
	tcpService = knet.NewTCPServer(&knet.TCPServerOptions{
		Addr:            c.String("server.tcp.addr"),
		ReadBufferSize:  c.Int("server.tcp.readbuffersize"),
		WriteBufferSize: c.Int("server.tcp.writebuffersize"),
		ClientHandler: func(conn net.Conn) {
			knet.HandleClient(conn, c.Int("server.tcp.readdeadline"), c.Int("server.tcp.writedeadline"), func(b []byte) ([]byte, error) {
				L, err := lua.LStateGet()
				if err != nil {
					log.Error(err)
					return nil, nil
				}

				defer lua.LStatePut(L)
				if err := L.DoFile("lua/router.lua"); err != nil {
					log.Error(err)
					return nil, nil
				}

				if err := L.CallByParam(glua.P{
					Fn:      L.GetGlobal("route"),
					NRet:    1,
					Protect: true,
				}, glua.LString(b)); err != nil {
					log.Error(err)
					return nil, nil
				}

				return nil, nil
			})
		},
	})

	err := tcpService.Serve()
	if err != nil {
		log.Errorln(err)
	}
}
