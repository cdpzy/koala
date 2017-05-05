package main

import (
	"net"

	log "github.com/Sirupsen/logrus"
	knet "github.com/doublemo/koala/net"
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
				return nil, nil
			})
		},
	})

	err := tcpService.Serve()
	if err != nil {
		log.Errorln(err)
	}
}
