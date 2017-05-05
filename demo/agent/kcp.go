package main

import (
	"net"

	log "github.com/Sirupsen/logrus"
	knet "github.com/doublemo/koala/net"
	cli "gopkg.in/urfave/cli.v2"
)

var kcpService *knet.KCPServer

// kcpServe KCP服务器
func kcpServe(c *cli.Context) {
	addr := c.String("server.kcp.addr")
	if len(addr) < 1 {
		return
	}

	kcpService = knet.NewKCPServer(&knet.KCPServerOptions{
		Addr:            addr,
		ReadBufferSize:  c.Int("server.kcp.readbuffersize"),
		WriteBufferSize: c.Int("server.kcp.writebuffersize"),
		Dscp:            c.Int("server.kcp.dscp"),
		Sndwnd:          c.Int("server.kcp.sndwnd"),
		Rcvwnd:          c.Int("server.kcp.rcvwnd"),
		Nodelay:         c.Int("server.kcp.nodelay"),
		Interval:        c.Int("server.kcp.interval"),
		Resend:          c.Int("server.kcp.resend"),
		Nc:              c.Int("server.kcp.nc"),
		Mtu:             c.Int("server.kcp.mtu"),
		ClientHandler: func(conn net.Conn) {
			knet.HandleClient(conn, c.Int("server.kcp.readdeadline"), c.Int("server.kcp.writedeadline"), func(b []byte) ([]byte, error) {
				return nil, nil
			})
		},
	})

	err := kcpService.Serve()
	if err != nil {
		log.Errorln(err)
	}
}
