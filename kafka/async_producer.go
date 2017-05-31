package kafka

import (
	"sync"

	"github.com/Shopify/sarama"
)

var _default_async_producer *AsyncProducerManager

// ASyncProducerManager 异步供应管理
type AsyncProducerManager struct {
	records map[string]sarama.AsyncProducer
	sync.RWMutex
}

func (smp *AsyncProducerManager) Register(k string, p sarama.AsyncProducer) {
	smp.Lock()
	smp.records[k] = p
	smp.Unlock()
}

func (smp *AsyncProducerManager) Unregister(k string) {
	smp.Lock()
	defer smp.Unlock()
	if m, ok := smp.records[k]; ok {
		m.Close()
		delete(smp.records, k)
	}
}

func (smp *AsyncProducerManager) Get(k string) (p sarama.AsyncProducer) {
	smp.RLock()
	p = smp.records[k]
	smp.RUnlock()
	return
}

func (smp *AsyncProducerManager) Iterator(f func(k string, v sarama.AsyncProducer) bool) {
	smp.RLock()
	defer smp.RUnlock()

	for k, v := range smp.records {
		smp.RUnlock()
		b := f(k, v)
		smp.RLock()
		if !b {
			break
		}
	}
}

func NewAsyncProducerManager() *AsyncProducerManager {
	return &AsyncProducerManager{records: make(map[string]sarama.AsyncProducer)}
}

func init() {
	_default_async_producer = NewAsyncProducerManager()
}

func GetAsyncProducer(k string) sarama.AsyncProducer {
	return _default_async_producer.Get(k)
}

func RegisterAsyncProducer(k string, p sarama.AsyncProducer) {
	_default_async_producer.Register(k, p)
}

func UnregisterAsyncProducer(k string) {
	_default_async_producer.Unregister(k)
}

func IteratorAsyncProducer(f func(k string, v sarama.AsyncProducer) bool) {
	_default_async_producer.Iterator(f)
}

func CreateAsyncProducer(k string, addrs []string, config *sarama.Config) error {
	producer, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		return err
	}

	_default_async_producer.Register(k, producer)
	return nil
}

func CloseAsyncProducerAll() {
	_default_async_producer.Iterator(func(k string, v sarama.AsyncProducer) bool {
		_default_async_producer.Unregister(k)
		return true
	})
}
