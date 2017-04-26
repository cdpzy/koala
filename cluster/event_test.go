package cluster

import "testing"
import "fmt"

func TestEvent(t *testing.T) {
	e1 := &Event{
		Name: EventNodeOffline,
		Data: make(map[string]interface{}),
		CallBack: func(e Event) {
			fmt.Println("e1:", e)
		},
	}

	e2 := &Event{
		Name: EventNodeOnline,
		Data: make(map[string]interface{}),
		CallBack: func(e Event) {
			fmt.Println("e2:", e)
		},
	}

	e3 := &Event{
		Name: EventNodeAttributeChanged,
		Data: make(map[string]interface{}),
		CallBack: func(e Event) {
			fmt.Println("e3:", e)
		},
	}

	e4 := &Event{
		Name: EventNodeServiceInit,
		Data: map[string]interface{}{"dddd": "xx"},
		CallBack: func(e Event) {
			fmt.Println("e4:", e)
		},
	}

	Nodes.events.Register(e1)
	Nodes.events.Register(e2)
	Nodes.events.Register(e3)
	Nodes.events.Register(e4)

	t.Log(Nodes.events)
}
