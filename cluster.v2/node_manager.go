package cluster

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
)

type NodeManager struct {
	nodes map[string]*Node
	sync.RWMutex
}

// Register 注册节点
func (nm *NodeManager) Register(n *Node) *Node {
	nm.Lock()
	node, ok := nm.nodes[n.Name]
	if !ok {
		nm.nodes[n.Name] = n
		node = n
	}
	nm.Unlock()

	if !ok {
		Events.Trigger(&Event{
			Name: EventNodeOnline,
			Node: n,
		})
	}

	return node
}

func (nm *NodeManager) RegisterNoTrigger(n *Node) {
	nm.Lock()
	nm.nodes[n.Name] = n
	nm.Unlock()
}

// UnRegister 移除节点
func (nm *NodeManager) UnRegister(nodeName string) {
	nm.Lock()
	node, ok := nm.nodes[nodeName]
	delete(nm.nodes, nodeName)
	nm.Unlock()

	if !ok {
		return
	}

	Events.Trigger(&Event{
		Name: EventNodeOffline,
		Node: node,
	})
}

// FindNodeByName 查找节点
func (nm *NodeManager) FindNodeByName(nodeName string) *Node {
	nm.RLock()
	defer nm.RUnlock()
	if n, ok := nm.nodes[nodeName]; ok {
		return n
	}

	return nil
}

// Join 节点上线
func (nm *NodeManager) Join(n *Node) error {
	path := nm.Path(n)
	n.Params.Set("Heartbeater", fmt.Sprint(time.Now().Unix()))
	n.Params.Iterator(func(k string, v interface{}) bool {
		if m, ok := v.(string); ok {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := etcdClient.Do(ctx, client.OpPut(path+"/"+k, m)); err != nil {
				log.Errorf("Node online failed: %v", err)
				return false
			}
		}
		return true
	})

	return nil
}

// Leave 退出
func (nm *NodeManager) Leave(n *Node) error {
	path := nm.Path(n) + "/"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := etcdClient.Do(ctx, client.OpDelete(path, client.WithPrefix())); err != nil {
		return err
	}

	return nil
}

// Path 拼合路径
func (nm *NodeManager) Path(n *Node) string {
	return fmt.Sprintf("%s/%s/%s", etcdURL, n.GetType(), n.Name)
}

func (nm *NodeManager) Iterator(f func(k string, v *Node) bool) {
	nm.RLock()
	defer nm.RUnlock()

	for k, v := range nm.nodes {
		nm.RUnlock()
		b := f(k, v)
		nm.RLock()
		if !b {
			break
		}
	}
}

// 获取所有节点
func (nm *NodeManager) All() []*Node {
	nm.RLock()
	defer nm.RUnlock()
	nodes := make([]*Node, 0)

	for _, n := range nm.nodes {
		nm.RUnlock()
		nodes = append(nodes, n)
		nm.RLock()
	}

	return nodes
}

func (nm *NodeManager) One(nodeName string) *Node {
	nm.RLock()
	node, ok := nm.nodes[nodeName]
	nm.RUnlock()

	if !ok {
		return nil
	}
	return node
}

func NewNodeManager() *NodeManager {
	return &NodeManager{nodes: make(map[string]*Node)}
}
