package kafka

import (
	"sync"

	"github.com/Shopify/sarama"
)

var (
	_default_clients *ClientManager
)

// ClientManager 连接管理
type ClientManager struct {
	records map[string]sarama.Client
	sync.RWMutex
}

func (cm *ClientManager) Register(k string, c sarama.Client) {
	cm.Lock()
	cm.records[k] = c
	cm.Unlock()
}

func (cm *ClientManager) Unregister(k string) {
	cm.Lock()
	defer cm.Unlock()
	if m, ok := cm.records[k]; ok {
		m.Close()
		delete(cm.records, k)
	}
}

func (cm *ClientManager) Get(k string) (c sarama.Client) {
	cm.RLock()
	c = cm.records[k]
	cm.RUnlock()
	return
}

func (cm *ClientManager) Iterator(f func(k string, v sarama.Client) bool) {
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
	return &ClientManager{records: make(map[string]sarama.Client)}
}

func init() {
	_default_clients = NewClientManager()
}

func GetClient(k string) sarama.Client {
	return _default_clients.Get(k)
}

func RegisterClient(k string, c sarama.Client) {
	_default_clients.Register(k, c)
}

func UnregisterClient(k string) {
	_default_clients.Unregister(k)
}

func IteratorClient(f func(k string, v sarama.Client) bool) {
	_default_clients.Iterator(f)
}

func CreateClient(k string, addrs []string, config *sarama.Config) error {
	cli, err := sarama.NewClient(addrs, config)
	if err != nil {
		return err
	}

	_default_clients.Register(k, cli)
	return nil
}

func CloseClientAll() {
	_default_clients.Iterator(func(k string, v sarama.Client) bool {
		if !v.Closed() {
			_default_clients.Unregister(k)
		}
		return true
	})
}
