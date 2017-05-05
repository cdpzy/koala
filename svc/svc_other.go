// +build !windows

package svc

import (
	"os"
	"syscall"

	cli "gopkg.in/urfave/cli.v2"
)

func WaitExit(start, stop func(c *cli.Context) error, c *cli.Context, sig ...os.Signal) error {
	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	err := start(c)
	if err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signalNotify(signalChan, sig...)

	<-signalChan

	return stop(c)
}
