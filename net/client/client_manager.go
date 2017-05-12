package client

import (
	"sync"
)

var (
	_anonymous_clients *ClientManager
)

// ClientManager 客户端管理
type ClientManager struct {
	records map[int64]*Client
	autoInc int64
	sync.RWMutex
}

func (cm *ClientManager) Register(id int64, c *Client) {
	cm.Lock()
	cm.autoInc++
	cm.records[id] = c
	cm.Unlock()
}

func (cm *ClientManager) Unregister(id int64) {
	cm.Lock()
	delete(cm.records, id)
	cm.Unlock()
}

func (cm *ClientManager) Get(id int64) (c *Client) {
	cm.RLock()
	c = cm.records[id]
	cm.RUnlock()

	return
}

func (cm *ClientManager) Count() (count int) {
	cm.RLock()
	count = len(cm.records)
	cm.RUnlock()
	return
}

func (cm *ClientManager) NewAutoID() (id int64) {
	cm.Lock()
	id = cm.autoInc + 1
	cm.Unlock()
	return
}

func (cm *ClientManager) Iterator(f func(int64, *Client) bool) {
	cm.RLock()
	defer cm.RUnlock()
	for k, v := range cm.records {
		cm.RUnlock()
		if b := f(k, v); !b {
			break
		}
		cm.RLock()
	}
}

func NewClientManager() *ClientManager {
	return &ClientManager{records: make(map[int64]*Client)}
}

func init() {
	_anonymous_clients = NewClientManager()
}

func GetAnonymous(id int64) *Client {
	return _anonymous_clients.Get(id)
}

func RegisterAnonymous(c *Client) {
	c.ID = _anonymous_clients.NewAutoID()
	_anonymous_clients.Register(c.ID, c)
}

func UnregisterAnonymous(id int64) {
	_anonymous_clients.Unregister(id)
}

func AnonymousIterator(f func(int64, *Client) bool) {
	_anonymous_clients.Iterator(f)
}
