package kafka

import (
	"sync"

	"github.com/Shopify/sarama"
)

var _default_sync_producer *SyncProducerManager

// SyncProducerManager 同步供应管理
type SyncProducerManager struct {
	records map[string]sarama.SyncProducer
	sync.RWMutex
}

func (smp *SyncProducerManager) Register(k string, p sarama.SyncProducer) {
	smp.Lock()
	smp.records[k] = p
	smp.Unlock()
}

func (smp *SyncProducerManager) Unregister(k string) {
	smp.Lock()
	defer smp.Unlock()
	if m, ok := smp.records[k]; ok {
		m.Close()
		delete(smp.records, k)
	}
}

func (smp *SyncProducerManager) Get(k string) (p sarama.SyncProducer) {
	smp.RLock()
	p = smp.records[k]
	smp.RUnlock()
	return
}

func (smp *SyncProducerManager) Iterator(f func(k string, v sarama.SyncProducer) bool) {
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

func NewSyncProducerManager() *SyncProducerManager {
	return &SyncProducerManager{records: make(map[string]sarama.SyncProducer)}
}

func init() {
	_default_sync_producer = NewSyncProducerManager()
}

func GetSyncProducer(k string) sarama.SyncProducer {
	return _default_sync_producer.Get(k)
}

func RegisterSyncProducer(k string, p sarama.SyncProducer) {
	_default_sync_producer.Register(k, p)
}

func UnregisterSyncProducer(k string) {
	_default_sync_producer.Unregister(k)
}

func IteratorSyncProducer(f func(k string, v sarama.SyncProducer) bool) {
	_default_sync_producer.Iterator(f)
}

func CreateSyncProducer(k string, addrs []string, config *sarama.Config) error {
	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return err
	}

	_default_sync_producer.Register(k, producer)
	return nil
}

func CloseSyncProducerAll() {
	_default_sync_producer.Iterator(func(k string, v sarama.SyncProducer) bool {
		_default_async_producer.Unregister(k)
		return true
	})
}
