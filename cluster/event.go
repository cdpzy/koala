package cluster

import (
	"sync"
)

const (
	EventNodeAttributeChanged string = "EventNodeAttributeChanged" // 节点属性变化事件
	EventNodeOnline           string = "EventNodeOnline"           // 节点上线
	EventNodeOffline          string = "EventNodeOffline"          // 节点下线
	EventNodeServiceInit      string = "EventNodeServiceInit"      // 节点服务
)

// EventCallBack 事件回调
type EventCallBack func(Event)

// Event 事件
type Event struct {
	Name     string                 // 事件名称
	CallBack EventCallBack          // 回调
	Node     *Node                  // 影响节点
	Data     map[string]interface{} //
}

// EventManager 事件管理器
type EventManager struct {
	mute   sync.RWMutex
	events map[string][]*Event //
}

// Register 注册事件
func (em *EventManager) Register(e *Event) {
	em.mute.Lock()
	defer em.mute.Unlock()

	if m, ok := em.events[e.Name]; ok {
		for _, i := range m {
			if e == i {
				return
			}
		}

		em.events[e.Name] = append(em.events[e.Name], e)
	} else {
		em.events[e.Name] = []*Event{e}
	}
}

// UnRegister 除去
func (em *EventManager) UnRegister(eventName string) {
	em.mute.Lock()
	delete(em.events, eventName)
	em.mute.Unlock()
}

// UnRegisterEvent 根据事件去除
func (em *EventManager) UnRegisterEvent(e *Event) {
	em.mute.Lock()
	defer em.mute.Unlock()

	if m, ok := em.events[e.Name]; ok {
		ev := make([]*Event, 0)
		for _, i := range m {
			if i == e {
				continue
			}

			ev = append(ev, i)
		}

		if len(ev) < 1 {
			delete(em.events, e.Name)
		} else {
			em.events[e.Name] = ev
		}

	}
}

// Trigger 事件触发
func (em *EventManager) Trigger(e *Event) {
	em.mute.RLock()
	ev, ok := em.events[e.Name]
	em.mute.RUnlock()

	if !ok {
		return
	}

	for _, i := range ev {
		if i.CallBack == nil {
			continue
		}

		data := make(map[string]interface{})
		if i.Data != nil {
			for k, v := range i.Data {
				data[k] = v
			}
		}

		if e.Data != nil {
			for k, v := range e.Data {
				data[k] = v
			}
		}

		e.Data = data
		i.CallBack(*e)
	}
}

// NewEventManager ..
func NewEventManager() *EventManager {
	return &EventManager{events: make(map[string][]*Event)}
}

// AddEvent 注册事件
func AddEvent(e *Event) {
	Nodes.events.Register(e)
}

// RemoveEvent 删除事件
func RemoveEvent(eventName string) {
	Nodes.events.UnRegister(eventName)
}

func RemoveEventFunc(e *Event) {
	Nodes.events.UnRegisterEvent(e)
}
