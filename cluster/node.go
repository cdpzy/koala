package cluster

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"

	"encoding/json"

	"reflect"
	"strings"

	"github.com/doublemo/koala/helper"

	"strconv"

	"sync"

	log "github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
)

// NodeStatus 节点状态类型
type NodeStatus int

const (
	NodeStatusOK     NodeStatus = iota // 正常节点
	NodeStatusClosed                   // 节点关闭
)

// Node 节点
type Node struct {
	Name        string            // 节点名称
	Type        string            // 节点类型
	Addr        net.IP            // 节点服务地址
	Port        int               // 服务端口
	Priority    int               // 优先选择权
	Status      NodeStatus        // 节点状态
	GRPCConn    *grpc.ClientConn  // 服务GRPC连接
	Params      map[string]string // 节点参数
	connecting  bool              // 服务是否正在连接中
	Heartbeater int64             // 心跳
	Idx         int               //
}

// Set 设置参数
func (n *Node) Set(k, v string) bool {

	t := reflect.TypeOf(n).Elem()
	field, ok := t.FieldByName(k)
	if !ok {
		return false
	}

	vset := reflect.ValueOf(n).Elem().FieldByName(k)
	if !vset.IsValid() {
		return false
	}

	if !vset.CanSet() {
		return false
	}

	switch field.Type.Kind() {
	case reflect.String:
		vset.SetString(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if vint, err := strconv.ParseInt(v, 10, 64); err == nil {
			vset.SetInt(vint)
		} else {
			return false
		}

	case reflect.TypeOf((net.IP)(nil)).Kind():
		vset.Set(reflect.ValueOf(net.ParseIP(v)))

	case reflect.TypeOf((NodeStatus)(0)).Kind():
		if vint, err := strconv.ParseInt(v, 10, 64); err == nil {
			vset.SetInt(vint)
		} else {
			return false
		}

	case reflect.MapOf(reflect.TypeOf(reflect.String), reflect.TypeOf(reflect.String)).Kind():
		data := make(map[string]string)
		err := json.Unmarshal([]byte(v), &data)
		if err != nil {
			return false
		}

		vset.Set(reflect.ValueOf(data))

	default:
		log.Warnln("Node attr parse:", k, v)
		return false
	}

	return true
}

// NodeManager 节点管理器
type NodeManager struct {
	nodes           map[string]*Node // 节点存储器
	Local           *Node            // 本地节点
	EtcdClient      *client.Client   // etcd 连接池
	EtcdURL         string           //
	watcherCtrl     chan struct{}    // 监听控制
	mute            sync.RWMutex     //
	events          *EventManager    //
	heartbeaterCtrl chan struct{}    // 心跳控制
}

// Register 注册节点
func (nm *NodeManager) Register(n *Node) {
	nm.mute.Lock()
	_, ok := nm.nodes[n.Name]
	nm.nodes[n.Name] = n
	nm.mute.Unlock()
	if !ok {
		nm.events.Trigger(&Event{
			Name: EventNodeOnline,
			Node: n,
			Data: make(map[string]interface{}),
		})
	}
}

func (nm *NodeManager) RegisterNoTrigger(n *Node) {
	nm.mute.Lock()
	nm.nodes[n.Name] = n
	nm.mute.Unlock()
}

// UnRegister 移除节点
func (nm *NodeManager) UnRegister(nodeName string) {
	nm.mute.Lock()
	node, ok := nm.nodes[nodeName]
	delete(nm.nodes, nodeName)
	nm.mute.Unlock()

	if !ok {
		return
	}

	nm.events.Trigger(&Event{
		Name: EventNodeOffline,
		Node: node,
		Data: make(map[string]interface{}),
	})
}

// FindNodeByName 查找节点
func (nm *NodeManager) FindNodeByName(nodeName string) *Node {
	nm.mute.RLock()
	defer nm.mute.RUnlock()
	if n, ok := nm.nodes[nodeName]; ok {
		return n
	}

	return nil
}

// Join 节点上线
func (nm *NodeManager) Join(n *Node) error {
	path := nm.Path(n)
	params, err := json.Marshal(n.Params)
	if err != nil {
		return err
	}

	ops := []client.Op{
		client.OpPut(path+"/Addr", n.Addr.String()),
		client.OpPut(path+"/Port", fmt.Sprint(n.Port)),
		client.OpPut(path+"/Priority", fmt.Sprint(n.Priority)),
		client.OpPut(path+"/Type", n.Type),
		client.OpPut(path+"/Status", fmt.Sprint(n.Status)),
		client.OpPut(path+"/Params", string(params)),
		client.OpPut(path+"/Heartbeater", fmt.Sprint(time.Now().Unix())),
	}

	for _, op := range ops {
		if _, err := nm.EtcdClient.Do(context.Background(), op); err != nil {
			return err
		}
	}
	return nil
}

// Leave 退出
func (nm *NodeManager) Leave(n *Node) error {
	path := nm.Path(n) + "/"
	if _, err := nm.EtcdClient.Do(context.Background(), client.OpDelete(path, client.WithPrefix())); err != nil {
		return err
	}

	return nil
}

// Path 拼合路径
func (nm *NodeManager) Path(n *Node) string {
	return fmt.Sprintf("%s/%s/%s", nm.EtcdURL, n.Type, n.Name)
}

// Preload 预加载已经存在节点
func (nm *NodeManager) Preload() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	resp, err := nm.EtcdClient.Get(ctx, nm.EtcdURL, client.WithPrefix())
	cancel()

	if err != nil {
		log.Warnln("Node preload:", err)
		return
	}

	for _, kvs := range resp.Kvs {
		key := string(kvs.Key)
		log.Debug("Preload node:", key, string(kvs.Value))
		if !strings.HasPrefix(key, nm.EtcdURL) || len(key) < 1 {
			continue
		}

		path := strings.TrimPrefix(key, nm.EtcdURL+"/")
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 3 {
			continue
		}

		nodeAttr := pathSplit[2]
		nodeName := pathSplit[1]
		nodeType := pathSplit[0]

		// 过滤本地节点
		if nodeName == nm.Local.Name {
			continue
		}

		node := nm.FindNodeByName(nodeName)
		if node == nil {
			node = &Node{
				Name:   nodeName,
				Type:   nodeType,
				Status: NodeStatusClosed,
				Params: make(map[string]string),
			}
		}

		node.Set(nodeAttr, string(kvs.Value))
		nm.Register(node)
	}
}

// Watcher 节点变化监听
func (nm *NodeManager) Watcher() {
	nm.watcherCtrl = make(chan struct{})
	readyed := make(chan struct{})
	go func() {
		defer helper.RecoverStack()
		defer func() {
			log.Infof("Cluster Watcher:stoped")
		}()

		log.Infof("Cluster Watcher:started")
		readyed <- struct{}{}
		watch := nm.EtcdClient.Watch(context.Background(), nm.EtcdURL, client.WithPrefix())
		for {
			select {
			case resp, ok := <-watch:
				if !ok {
					return
				}

				for _, ev := range resp.Events {
					key := string(ev.Kv.Key)
					log.Debug("Watcher:", key, string(ev.Kv.Value), " - ", ev.Type)
					if !strings.HasPrefix(key, nm.EtcdURL) || len(key) < 1 {
						continue
					}

					path := strings.TrimPrefix(key, nm.EtcdURL+"/")
					pathSplit := strings.Split(path, "/")
					if len(pathSplit) != 3 {
						continue
					}

					nodeAttr := pathSplit[2]
					nodeName := pathSplit[1]
					nodeType := pathSplit[0]

					switch ev.Type {
					case client.EventTypeDelete:
						// 如果删除需要删除莫个节点，那个删除Addr才能成功将节点移除
						if nodeAttr != "Addr" {
							continue
						}

						node := nm.FindNodeByName(nodeName)
						if node == nil {
							continue
						}

						if node.GRPCConn != nil {
							node.GRPCConn.Close()
						}

						nm.UnRegister(node.Name)

					case client.EventTypePut:
						node := nm.FindNodeByName(nodeName)
						if node == nil {
							node = &Node{
								Name:   nodeName,
								Type:   nodeType,
								Status: NodeStatusClosed,
								Params: make(map[string]string),
							}
						}

						node.Set(nodeAttr, string(ev.Kv.Value))

						if nodeAttr == "Heartbeater" {
							nm.RegisterNoTrigger(node)
						} else {
							nm.Register(node)
						}

						if nodeAttr == "Addr" {
							node.GRPCConn = nil
						}

						nm.events.Trigger(&Event{
							Name: EventNodeAttributeChanged,
							Node: node,
							Data: map[string]interface{}{"attribute": nodeAttr, "attributeVal": string(ev.Kv.Value)},
						})

						if node.GRPCConn == nil &&
							!node.connecting &&
							nm.Local != nil &&
							node.Name != nm.Local.Name &&
							node.Status == NodeStatusOK &&
							node.Addr != nil &&
							node.Port > 0 {

							node.connecting = true
							go func() {
								ticker := time.NewTicker(5 * time.Second)
								for {
									if node.GRPCConn != nil || !services[node.Type] {
										return
									}

									err := nm.serviceInit(node)
									if err == nil {
										node.connecting = false
										return
									}

									log.Warningf("Connect to %s[%s:%d] failed", node.Name, node.Addr, node.Port)
									select {
									case <-ticker.C:
									}
								}
							}()
						}
					}
				}

			case <-nm.watcherCtrl:
				return
			}
		}
	}()

	<-readyed
	close(readyed)
}

// CloseWatcher 关闭
func (nm *NodeManager) CloseWatcher() {
	if nm.watcherCtrl != nil {
		nm.watcherCtrl <- struct{}{}
	}
}

// ServiceInit 连接服务
func (nm *NodeManager) ServiceInit() {
	nm.mute.Lock()
	defer nm.mute.Unlock()

	for _, node := range nm.nodes {
		nm.mute.Unlock()
		err := nm.serviceInit(node)
		if err != nil {
			log.Warningf("Connect to %s[%s:%d] failed", node.Name, node.Addr, node.Port)
		}

		nm.mute.Lock()
	}
}

// serviceInit 服务连接
func (nm *NodeManager) serviceInit(node *Node) error {
	if !services[node.Type] {
		return nil
	}

	if node.Addr == nil || node.Port < 1 {
		return fmt.Errorf("Invalid address:%s %s %s %d", node.Name, node.Type, node.Addr.String(), node.Port)
	}

	addr := fmt.Sprintf("%s:%d", node.Addr.String(), node.Port)
	log.Debug("Waiting connect to node:", addr)
	conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		log.Debug("Connect to node:", addr, "[failed]", err)
		return fmt.Errorf("Connect to %s failed, error:%v", addr, err)
	}

	log.Debug("Connect to node:", addr, "[OK]")
	node.GRPCConn = conn
	nm.events.Trigger(&Event{
		Name: EventNodeServiceInit,
		Node: node,
		Data: make(map[string]interface{}),
	})
	return nil
}

// heartbeater 节点心跳
func (nm *NodeManager) heartbeater() {
	nm.heartbeaterCtrl = make(chan struct{})
	go func() {
		defer func() {
			log.Infoln("Node heartbeater stoped.")
		}()

		log.Infoln("Node heartbeater started.")

		ticker := time.After(5 * time.Second)
		for {
			select {
			case <-ticker:
				nodes := nm.All()
				for _, n := range nodes {
					s := time.Now().Sub(time.Unix(n.Heartbeater, 0)).Seconds()
					if s > 10 {
						nm.UnRegister(n.Name)
					}
				}

				if nm.Local != nil {
					nm.EtcdClient.Do(context.Background(), client.OpPut(nm.Path(nm.Local)+"/Heartbeater", fmt.Sprint(time.Now().Unix())))
				}

				ticker = time.After(5 * time.Second)

			case <-nm.heartbeaterCtrl:
				return
			}
		}
	}()
}

// closeHeartbeater节点心跳关闭
func (nm *NodeManager) closeHeartbeater() {
	if nm.heartbeaterCtrl != nil {
		nm.heartbeaterCtrl <- struct{}{}
	}
}

// 获取所有节点
func (nm *NodeManager) All() []*Node {
	nm.mute.RLock()
	defer nm.mute.RUnlock()
	nodes := make([]*Node, 0)

	for _, n := range nm.nodes {
		nm.mute.RUnlock()
		nodes = append(nodes, n)
		nm.mute.RLock()
	}

	return nodes
}

func (nm *NodeManager) One(nodeName string) *Node {
	nm.mute.RLock()
	node, ok := nm.nodes[nodeName]
	nm.mute.RUnlock()

	if !ok {
		return nil
	}

	return node
}

// NewNodeManager 节点管理器
func NewNodeManager() *NodeManager {
	return &NodeManager{nodes: make(map[string]*Node), events: NewEventManager()}
}
