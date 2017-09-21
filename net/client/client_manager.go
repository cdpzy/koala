package client

import (
	"fmt"
	"sync"
)

var (
	_default_clients *ClientManager
)

// ClientManager 客户端管理
type ClientManager struct {
	records map[string]*Client
	autoInc int64
	sync.RWMutex
}

func (cm *ClientManager) Register(id string, c *Client) {
	cm.Lock()
	cm.autoInc++
	cm.records[id] = c
	cm.Unlock()
}

func (cm *ClientManager) Unregister(id string) {
	cm.Lock()
	delete(cm.records, id)
	cm.Unlock()
}

func (cm *ClientManager) Get(id string) (c *Client) {
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

func (cm *ClientManager) NewAutoID() (id string) {
	cm.Lock()
	id = fmt.Sprintf("KC-%d", cm.autoInc+1)
	cm.Unlock()
	return
}

func (cm *ClientManager) Iterator(f func(string, *Client) bool) {
	cm.RLock()
	defer cm.RUnlock()
	for k, v := range cm.records {
		cm.RUnlock()
		b := f(k, v)
		cm.RLock()
		if !b {
			break
		}
	}
}

func NewClientManager() *ClientManager {
	return &ClientManager{records: make(map[string]*Client)}
}

func init() {
	_default_clients = NewClientManager()
}

func Get(id string) *Client {
	return _default_clients.Get(id)
}

func Count() int {
	return _default_clients.Count()
}

func Register(c *Client) {
	c.ID = _default_clients.NewAutoID()
	c.uniqueID = c.ID
	_default_clients.Register(c.ID, c)
}

func RegisterID(id string, c *Client) {
	_default_clients.Register(id, c)
}

func Unregister(id string) {
	_default_clients.Unregister(id)
}

func Iterator(f func(string, *Client) bool) {
	_default_clients.Iterator(f)
}
