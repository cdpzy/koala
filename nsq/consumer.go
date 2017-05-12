package nsq

import (
	"net"
	"sync"

	gnsq "github.com/go-nsq"
)

var _default_consumer *ConsumerManager

type ConsumerManager struct {
	records map[string]*gnsq.Consumer
	sync.RWMutex
}

func (cm *ConsumerManager) Register(k string, c *gnsq.Consumer) {
	cm.Lock()
	cm.records[k] = c
	cm.Unlock()
}

func (cm *ConsumerManager) Unregister(k string) {
	cm.Lock()
	delete(cm.records, k)
	cm.Unlock()
}

func (cm *ConsumerManager) Get(k string) (c *gnsq.Consumer) {
	cm.RLock()
	c = cm.records[k]
	cm.RUnlock()
	return
}

func NewConsumerManager() *ConsumerManager {
	return &ConsumerManager{records: make(map[string]*gnsq.Consumer)}
}

func init() {
	_default_consumer = NewConsumerManager()
}

func CreateConsumer(key, addr, topic, channel string) (*gnsq.Consumer, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	conf := gnsq.NewConfig()
	conf.LocalAddr = laddr
	cons, err := gnsq.NewConsumer(topic, channel, conf)
	if err != nil {
		return nil, err
	}

	_default_consumer.Register(key, cons)
	return nil, nil
}

func GetConsumer(key string) *gnsq.Consumer {
	return _default_consumer.Get(key)
}
