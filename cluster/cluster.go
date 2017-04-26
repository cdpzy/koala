package cluster

import (
	"math/rand"
	"net"
	"strings"
	"time"

	"container/heap"

	log "github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
)

const (
	defaultPort int = 6901 // 默认端口
)

var (
	Nodes    *NodeManager    = NewNodeManager()      // 节点管理器
	services map[string]bool = make(map[string]bool) // 允许的服务
)

// Options 节点参数
type Options struct {
	Endpoints   []string // etcd集群地址
	DialTimeout int      // etcd连接超时
	ETCDUrl     string   // etcd服务监听地址
	Services    []string // 连接服务
}

// New 创建集群
func New(o *Options) (*NodeManager, error) {
	etcdClient, err := client.New(client.Config{
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

	Nodes.EtcdClient = etcdClient
	Nodes.EtcdURL = o.ETCDUrl
	Nodes.Local = &Node{
		Port:   defaultPort,
		Params: make(map[string]string),
		Status: NodeStatusClosed,
	}

	conn, err := net.Dial("udp", "google.com:80")
	if err == nil {
		defer conn.Close()
		Nodes.Local.Addr = net.ParseIP(strings.Split(conn.LocalAddr().String(), ":")[0])
	} else {
		log.Warnln("Get local addrs:", err)
	}

	Nodes.Watcher()
	Nodes.Preload()

	if err = Nodes.ServiceInit(); err != nil {
		return nil, err
	}

	Nodes.heartbeater()
	return Nodes, nil
}

// Shutdown 节点下线
func Shutdown() {
	Nodes.closeHeartbeater()
	Nodes.Leave(Nodes.Local)
	// waiting
	time.Sleep(1 * time.Second)

	Nodes.CloseWatcher()
}

// GetServiceBySort 获取服务节点
func GetServiceBySort(typeName string) *Node {
	nodes := Nodes.All()
	q := make(Queue, 0)
	for _, n := range nodes {
		if n.Status != NodeStatusOK || n.Type != typeName || n.GRPCConn == nil {
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
		if n.Status != NodeStatusOK || n.Type != typeName || n.GRPCConn == nil {
			continue
		}

		for i := 0; i < n.Priority; i++ {
			q = append(q, n)
		}
	}

	if len(q) < 1 {
		return nil
	}

	return q[rand.Intn(len(q))]
}

// GetServiceAll 获取指定类型的所有服务
func GetServiceAll(typeName string) []*Node {
	nodes := Nodes.All()
	q := make([]*Node, 0)
	for _, n := range nodes {
		if n.Status != NodeStatusOK || n.Type != typeName || n.GRPCConn == nil {
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
		q = append(q, n)
	}

	if len(q) < 1 {
		return nil
	}

	return q
}
