package cluster

import (
	"fmt"
	"strings"
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

func (nm *NodeManager) Count() (count int) {
	nm.RLock()
	count = len(nm.nodes)
	nm.RUnlock()

	return
}

func NewNodeManager() *NodeManager {
	return &NodeManager{nodes: make(map[string]*Node)}
}

// Preload 预加载已经存在节点
func preload() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcdClient.Get(ctx, etcdURL, client.WithPrefix())
	cancel()

	if err != nil {
		log.Warnln("Node preload:", err)
		return
	}

	for _, kvs := range resp.Kvs {
		key := string(kvs.Key)
		log.Debug("Preload node:", key, string(kvs.Value))
		if !strings.HasPrefix(key, etcdURL) || len(key) < 1 {
			continue
		}

		path := strings.TrimPrefix(key, etcdURL+"/")
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 3 {
			continue
		}

		nodeAttr := pathSplit[2]
		nodeName := pathSplit[1]
		nodeType := pathSplit[0]

		// 过滤本地节点
		if nodeName == Local.Name {
			continue
		}

		node := Nodes.FindNodeByName(nodeName)
		if node == nil {
			node = NewNode(nodeName)
			node.SetType(nodeType)
			node.SetStatus(NodeStatusClosed)
		}

		node.Params.Set(nodeAttr, string(kvs.Value))
		Nodes.Register(node)
	}

	// connect
	Nodes.Iterator(func(k string, v *Node) bool {
		if err := serviceInit(v); err != nil {
			port, _ := v.GetPort()
			log.Warningf("Connect to %s[%s:%d] failed", v.Name, v.GetAddr().String(), port)
		}

		return true
	})
}

func reload(name, ntype string) (*Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcdClient.Get(ctx, etcdURL+"/"+ntype+"/"+name, client.WithPrefix())
	cancel()

	if err != nil {
		return nil, err
	}

	node := NewNode(name)
	node.SetType(ntype)
	for _, kvs := range resp.Kvs {
		key := string(kvs.Key)
		if !strings.HasPrefix(key, etcdURL) || len(key) < 1 {
			continue
		}

		path := strings.TrimPrefix(key, etcdURL+"/")
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 3 {
			continue
		}

		nodeAttr := pathSplit[2]
		nodeName := pathSplit[1]
		nodeType := pathSplit[0]

		// 过滤本地节点
		if nodeName == Local.Name || nodeName != name || nodeType != ntype {
			continue
		}

		node.Params.Set(nodeAttr, string(kvs.Value))
	}

	return node, nil
}
