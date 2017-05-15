package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEManager(t *testing.T) {
	n := &EventManager{records: make(map[string][]Event)}
	a := &DefaultEvent{name: "dv", callback: func(e Event) {
		fmt.Println(e.Name(), time.Now().UnixNano())
	}}

	for i := 0; i < 100; i++ {
		go func() {
			for {
				time.Sleep(time.Nanosecond * 10)
				n.Register(a)
			}
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			for {
				time.Sleep(time.Nanosecond * 20)
				n.Unregister(a)
			}
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			for {
				time.Sleep(time.Nanosecond * 15)
				n.Trigger(a)
			}
		}()
	}

	time.Sleep(time.Second * 10)
}
