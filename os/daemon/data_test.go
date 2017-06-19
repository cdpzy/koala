package daemon

import "testing"
import "fmt"

func TestData(t *testing.T) {
	d := &Data{
		Pid:      12333,
		Command:  190,
		NodeName: "agent",
		NodeAddr: "127.0.0.1",
		NodePort: 12013,
		NodeType: "agent",
	}

	x := d.Encode()
	fmt.Println("x:", x)
	xd := &Data{}
	xd.Decode(x)
	fmt.Println("x:", xd)

}
