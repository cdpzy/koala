package cluster

import (
	"fmt"
	"net"
	"testing"
)

func TestNodeSet(t *testing.T) {
	n := &Node{
		Name:   "test1",
		Params: NewParams(),
	}

	n.SetType("ddddd")
	fmt.Println(n.GetType())
	n.SetAddr(net.ParseIP("127.0.0.1"))
	fmt.Println(n.GetAddr())

}

func TestReload(t *testing.T) {

}
