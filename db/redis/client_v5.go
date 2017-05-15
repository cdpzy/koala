package redis

import (
	"sync"

	redisv5 "gopkg.in/redis.v5"
)

var (
	_defautls_cv5m *ClientV5Manager
)

type ClientV5 struct {
	client    *redisv5.Client
	cclient   *redisv5.ClusterClient
	isCluster bool
}

func (c *ClientV5) IsCluster() bool {
	return c.isCluster
}

func (c *ClientV5) Client() *redisv5.Client {
	return c.client
}

func (c *ClientV5) CClient() *redisv5.ClusterClient {
	return c.cclient
}

type ClientV5Manager struct {
	sync.RWMutex
	records map[string]*ClientV5
}

func (cm *ClientV5Manager) Register(k string, c *ClientV5) {
	cm.Lock()
	cm.records[k] = c
	cm.Unlock()
}

func (cm *ClientV5Manager) Unregister(k string) {
	cm.Lock()
	defer cm.Unlock()

	if m, ok := cm.records[k]; ok {
		cm.Unlock()
		if m.IsCluster() {
			m.CClient().Close()
		} else {
			m.Client().Close()
		}
		cm.Lock()
	}

	delete(cm.records, k)
}

func (cm *ClientV5Manager) Get(k string) (c *ClientV5) {
	cm.RLock()
	c = cm.records[k]
	cm.RUnlock()
	return
}

func NewClientV5Manager() *ClientV5Manager {
	return &ClientV5Manager{records: make(map[string]*ClientV5)}
}

func init() {
	_defautls_cv5m = NewClientV5Manager()
}

func CreateV5Cluster(key string, op *redisv5.ClusterOptions) (*ClientV5, error) {
	c := &ClientV5{}
	c.isCluster = true
	c.cclient = redisv5.NewClusterClient(op)

	if err := c.cclient.Ping().Err(); err != nil {
		return nil, err
	}

	_defautls_cv5m.Register(key, c)
	return c, nil
}

func CreateV5(key string, op *redisv5.Options) (*ClientV5, error) {
	c := &ClientV5{}
	c.isCluster = false
	c.client = redisv5.NewClient(op)

	if err := c.client.Ping().Err(); err != nil {
		return nil, err
	}

	_defautls_cv5m.Register(key, c)
	return c, nil
}

func GetV5(key string) *ClientV5 {
	return _defautls_cv5m.Get(key)
}

func CloseV5(key string) {
	_defautls_cv5m.Unregister(key)
}
