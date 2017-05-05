package main

import (
	log "github.com/Sirupsen/logrus"
	khttp "github.com/doublemo/koala/net/http"
	cli "gopkg.in/urfave/cli.v2"
)

var httpService *khttp.Server

// httpServe HTTP服务器
func httpServe(c *cli.Context) {
	addr := c.String("server.http.addr")
	if len(addr) < 1 {
		return
	}

	httpService = khttp.NewServer(&khttp.ServerOptions{
		Addr:           addr,
		ReadTimeout:    c.Int("server.http.readtimeout"),
		WriteTimeout:   c.Int("server.http.writetimeout"),
		Ssl:            c.Bool("server.http.ssl"),
		SslKey:         c.String("server.http.sslkey"),
		SslCert:        c.String("server.http.sslcert"),
		MaxRequestSize: c.Int64("server.http.maxrequestsize"),
	})

	err := httpService.Serve()
	if err != nil {
		log.Errorln(err)
	}
}
