package cluster

import (
	"container/heap"
	"context"
	"math/rand"
	"net"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
)

var (
	Local      *Node                                   // 本地节点
	etcdClient *client.Client                          // etcd 连接池
	etcdURL    string                                  //
	Nodes      *NodeManager    = NewNodeManager()      // 节点管理器
	services   map[string]bool = make(map[string]bool) // 允许的服务
	Events     *EventManager   = NewEventManager()     // 事件
)

// Options 节点参数
type Options struct {
	LocalName   string   // 当前节点名称
	Endpoints   []string // etcd集群地址
	DialTimeout int      // etcd连接超时
	ETCDUrl     string   // etcd服务监听地址
	Services    []string // 连接服务
}

// New 创建集群
func New(o *Options) (*NodeManager, error) {
	var err error
	etcdClient, err = client.New(client.Config{
		Endpoints:   o.Endpoints,
		DialTimeout: time.Duration(o.DialTimeout) * time.Second,
	})

	if err != nil {
		return nil, err
	}

	if o.Services != nil {
		for _, service := range o.Services {
			services[service] = true
		}
	}

	etcdURL = o.ETCDUrl
	Local = NewNode(o.LocalName)
	Local.SetPort(defaultPort)
	Local.SetStatus(NodeStatusClosed)

	conn, err := net.Dial("udp", "google.com:80")
	if err == nil {
		conn.Close()
		Local.SetAddr(net.ParseIP(strings.Split(conn.LocalAddr().String(), ":")[0]))
	} else {
		log.Warnln("Get local addrs:", err)
	}

	watcher()
	preload()
	heartbeater()
	return Nodes, nil
}

// Shutdown 节点下线
func Shutdown() {
	stopHeartbeater()
	Nodes.Leave(Local)
	// waiting
	time.Sleep(1 * time.Second)
	stopWatcher()
}

// GetServiceBySort 获取服务节点
func GetServiceBySort(typeName string) *Node {
	nodes := Nodes.All()
	q := make(Queue, 0)
	for _, n := range nodes {
		if n.GetStatus() != NodeStatusOK || n.GetType() != typeName || n.GetGRPCConn() == nil {
			continue
		}

		q = append(q, n)
	}

	if len(q) < 1 {
		return nil
	}

	heap.Init(&q)
	node := heap.Pop(&q).(*Node)
	return node
}

// GetServiceByRandom 随机获取服务节点
func GetServiceByRandom(typeName string) *Node {
	nodes := Nodes.All()
	q := make([]*Node, 0)
	for _, n := range nodes {
		priority, _ := n.GetPriority()
		if n.GetStatus() != NodeStatusOK || n.GetType() != typeName || n.GetGRPCConn() == nil || priority < 1 {
			continue
		}

		for i := 0; i < priority; i++ {
			q = append(q, n)
		}
	}

	if len(q) < 1 {
		return nil
	}

	ShuffleNode(q)
	return q[rand.Intn(len(q))]
}

// GetServiceAll 获取指定类型的所有服务
func GetServiceAll(typeName string) []*Node {
	nodes := Nodes.All()
	q := make([]*Node, 0)
	for _, n := range nodes {
		if n.GetStatus() != NodeStatusOK || n.GetType() != typeName || n.GetGRPCConn() == nil {
			continue
		}

		q = append(q, n)
	}

	if len(q) < 1 {
		return nil
	}

	return q
}

// Find 获取某个节点
func Find(nodeName string) *Node {
	return Nodes.One(nodeName)
}

func FindByType(t string) []*Node {
	nodes := Nodes.All()
	q := make([]*Node, 0)
	for _, n := range nodes {
		if n.GetType() != t {
			continue
		}

		q = append(q, n)
	}

	if len(q) < 1 {
		return nil
	}

	return q
}

func ShuffleNode(a []*Node) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
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
