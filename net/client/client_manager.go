package net

import (
	"sync"
)

var (
	_default_clients *ClientManager
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

func (cm *ClientManager) AutoInc() (id int64) {
	cm.Lock()
	id = cm.autoInc
	cm.Unlock()
	return
}
