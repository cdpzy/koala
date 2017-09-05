package http

import (
	"fmt"
	"testing"
)

func TestPostJSON(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(PostJSON("http://192.168.2.19:19028", make([]byte, 0)))
	}
}
