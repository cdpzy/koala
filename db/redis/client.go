package redis

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	_defautls_cm *ClientManager
)

type Client struct {
	client    *redis.Client
	cclient   *redis.ClusterClient
	isCluster bool
}

func (c *Client) IsCluster() bool {
	return c.isCluster
}

func (c *Client) Client() *redis.Client {
	return c.client
}

func (c *Client) CClient() *redis.ClusterClient {
	return c.cclient
}

type ClientManager struct {
	sync.RWMutex
	records map[string]*Client
}

func (cm *ClientManager) Register(k string, c *Client) {
	cm.Lock()
	cm.records[k] = c
	cm.Unlock()
}

func (cm *ClientManager) Unregister(k string) {
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

func (cm *ClientManager) Get(k string) (c *Client) {
	cm.RLock()
	c = cm.records[k]
	cm.RUnlock()
	return
}

func NewClientManager() *ClientManager {
	return &ClientManager{records: make(map[string]*Client)}
}

func init() {
	_defautls_cm = NewClientManager()
}

func CreateCluster(key string, op *redis.ClusterOptions) (*Client, error) {
	c := &Client{}
	c.isCluster = true
	c.cclient = redis.NewClusterClient(op)

	if err := c.cclient.Ping().Err(); err != nil {
		return nil, err
	}

	_defautls_cm.Register(key, c)
	return c, nil
}

func Create(key string, op *redis.Options) (*Client, error) {
	c := &Client{}
	c.isCluster = false
	c.client = redis.NewClient(op)

	if err := c.client.Ping().Err(); err != nil {
		return nil, err
	}

	_defautls_cm.Register(key, c)
	return c, nil
}

func Get(key string) *Client {
	return _defautls_cm.Get(key)
}

func Close(key string) {
	_defautls_cm.Unregister(key)
}
