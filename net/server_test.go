package net

import "testing"
import "net"
import "fmt"
import log "github.com/Sirupsen/logrus"

func TestTCPServer(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	s := NewTCPServer(&TCPServerOptions{
		Addr:            ":6062",
		ReadBufferSize:  32767,
		WriteBufferSize: 32767,
		ClientHandler: func(conn net.Conn) {
			HandleClient(conn, 30, 30, func(b []byte) ([]byte, error) {
				fmt.Println("b:", b)
				return b, nil
			})
		},
	})

	s.Serve()
}
