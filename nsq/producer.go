package nsq

import (
	"sync"

	gnsq "github.com/go-nsq"
)

var _default_producer *ProducerManager

type Producer struct {
	Topic    string
	Producer *gnsq.Producer
}

type ProducerManager struct {
	records map[string]*Producer
	sync.RWMutex
}

func (pm *ProducerManager) Register(k string, p *Producer) {
	pm.Lock()
	pm.records[k] = p
	pm.Unlock()
}

func (pm *ProducerManager) Unregister(k string) {
	pm.Lock()
	delete(pm.records, k)
	pm.Unlock()
}

func (pm *ProducerManager) Get(k string) (p *Producer) {
	pm.RLock()
	p = pm.records[k]
	pm.RUnlock()
	return
}

func NewProducerManager() *ProducerManager {
	return &ProducerManager{records: make(map[string]*Producer)}
}

func init() {
	_default_producer = NewProducerManager()
}

func CreateProducer(key, addr, topic string) (*gnsq.Producer, error) {
	conf := gnsq.NewConfig()
	pro, err := gnsq.NewProducer(addr, conf)
	if err != nil {
		return nil, err
	}

	_default_producer.Register(key, &Producer{Topic: topic, Producer: pro})
	return nil, nil
}

func GetProducer(key string) *Producer {
	return _default_producer.Get(key)
}

func Send(p *Producer, b []byte) error {
	return p.Producer.Publish(p.Topic, b)
}
