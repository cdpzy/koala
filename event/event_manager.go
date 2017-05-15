package event

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrFound = errors.New("ErrFound") // 事件已经存在
)

// EventManager 事件管理器
type EventManager struct {
	records map[string][]Event
	sync.RWMutex
}

// Register 事件注册
func (em *EventManager) Register(e Event) error {
	em.Lock()
	defer em.Unlock()

	events, ok := em.records[e.Name()]
	if ok {
		for _, ev := range events {
			a := reflect.ValueOf(ev.CallBack())
			b := reflect.ValueOf(e.CallBack())
			if a.Pointer() == b.Pointer() {
				return ErrFound
			}
		}

		em.records[e.Name()] = append(em.records[e.Name()], e)
	} else {
		em.records[e.Name()] = []Event{e}
	}

	return nil
}

// Unregister 事件注销
func (em *EventManager) Unregister(e Event) {
	em.Lock()
	defer em.Unlock()

	events, ok := em.records[e.Name()]
	if !ok {
		return
	}

	nevents := make([]Event, 0)
	for _, ev := range events {
		a := reflect.ValueOf(ev.CallBack())
		b := reflect.ValueOf(e.CallBack())
		if a.Pointer() == b.Pointer() {
			continue
		}

		nevents = append(nevents, ev)
	}

	em.records[e.Name()] = nevents
}

// Get 获取事件
func (em *EventManager) Get(name string) []Event {
	em.RLock()
	events, ok := em.records[name]
	em.RUnlock()

	if !ok {
		return []Event{}
	}

	return events
}

// Trigger 事件触发
func (em *EventManager) Trigger(e Event) {
	em.RLock()
	events, ok := em.records[e.Name()]
	em.RUnlock()

	if !ok {
		return
	}

	for _, ev := range events {
		ev.CallBack()(e)
		if ev.IsStopPropagation() {
			break
		}
	}
}

func NewEventManager() *EventManager {
	return &EventManager{records: make(map[string][]Event)}
}
