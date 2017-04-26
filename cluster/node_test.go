package cluster

import (
	"fmt"
	"net"
	"testing"
	"time"
	// log "github.com/Sirupsen/logrus"
)

func TestNew(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	fmt.Println(New(&Options{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 30,
		ETCDUrl:     "/chkbackends",
		Services:    []string{"test"},
	}))
	Nodes.Local.Name = "test"
	Nodes.Local.Type = "test"
	Nodes.Local.Priority = 1000
	Nodes.Local.Addr = net.ParseIP("192.168.18.152")
	Nodes.Local.Status = 3
	Nodes.Local.Params = map[string]string{"dd": "dxoxox"}
	Nodes.Local.Port = 19024
	Nodes.Join(Nodes.Local)

	n := &Node{
		Name:     "test2",
		Type:     "test",
		Priority: 1000,
		Addr:     net.ParseIP("192.168.18.152"),
		Status:   1,
		Params:   map[string]string{"dd": "dxoxox"},
		Port:     19024,
	}
	Nodes.Join(n)

	t.Log("--------------nodes:")
	for _, r := range Nodes.nodes {
		t.Log(r)
	}

	time.Sleep(20 * time.Second)
	fmt.Println(GetServiceBySort(""))
	go Shutdown()
	go Shutdown()
	go Shutdown()
	go Shutdown()
	time.Sleep(10 * time.Second)
}
