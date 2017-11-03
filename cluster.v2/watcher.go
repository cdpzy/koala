package cluster

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
	"github.com/doublemo/koala/helper"
)

var watcherCtrl chan struct{} // 监听控制

func watcher() {
	watcherCtrl = make(chan struct{})
	readyed := make(chan struct{})
	go func() {
		defer helper.RecoverStackPanic()
		defer func() {
			log.Infof("Cluster Watcher:stoped")
		}()

		log.Infof("Cluster Watcher:started")
		readyed <- struct{}{}
		watch := etcdClient.Watch(context.Background(), etcdURL, client.WithPrefix())
		for {
			select {
			case resp, ok := <-watch:
				if !ok {
					return
				}

				for _, ev := range resp.Events {
					key := string(ev.Kv.Key)
					if !strings.HasSuffix(key, "/Heartbeater") {
						log.Debug("Watcher:", key, string(ev.Kv.Value), " - ", ev.Type)
					}

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

					switch ev.Type {
					case client.EventTypeDelete:
						eventTypeDelete(nodeName, nodeType, nodeAttr)
					case client.EventTypePut:
						eventTypePut(nodeName, nodeType, nodeAttr, string(ev.Kv.Value))
					}

				}

			case <-watcherCtrl:
				return
			}
		}
	}()
	<-readyed
}

func eventTypePut(nodeName, nodeType, nodeAttr, nodeVale string) {
	node := Nodes.FindNodeByName(nodeName)
	if node == nil {
		if nodeAttr == "Heartbeater" {
			n, err := reload(nodeName, nodeType)
			if err != nil {
				return
			}
			node = n
		} else {
			node = NewNode(nodeName)
			node.SetType(nodeType)
			node.SetStatus(NodeStatusClosed)
		}

		node = Nodes.Register(node)
	}

	node.Params.Set(nodeAttr, nodeVale)
	if nodeAttr == "Addr" {
		node.RemoveGRPCConn()
	}

	if nodeAttr != "Heartbeater" {
		Events.Trigger(&Event{
			Name: EventNodeAttributeChanged,
			Node: node,
		})
	}

	//log.Infof("TEST-----------------------%s-%s-%s:[%d][%d]", nodeName, nodeType, nodeAttr, Nodes.Count(), node.Params.Count())
	port, _ := node.GetPort()
	if node.GetGRPCConn() == nil &&
		!node.Params.Exist("connecting") &&
		Local != nil &&
		node.Name != Local.Name &&
		node.GetStatus() == NodeStatusOK &&
		node.GetAddr() != nil &&
		port > 0 &&
		services[node.GetType()] {
		node.Params.Set("connecting", "true")
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				if node.GetGRPCConn() != nil && !services[node.GetType()] {
					return
				}

				if err := serviceInit(node); err == nil {
					node.Params.Remove("connecting")
					return
				}

				mport, _ := node.GetPort()
				log.Warningf("Connect to %s[%s:%d] failed", node.Name, node.GetAddr().String(), mport)
				<-ticker.C
				if node.GetStatus() != NodeStatusOK {
					node.Params.Remove("connecting")
					return
				}
			}
		}()
	}
}

func eventTypeDelete(nodeName, nodeType, nodeAttr string) {
	// 如果删除需要删除莫个节点，那个删除Addr才能成功将节点移除
	if nodeAttr != "Addr" {
		return
	}

	node := Nodes.FindNodeByName(nodeName)
	if node == nil {
		return
	}

	if m := node.GetGRPCConn(); m != nil {
		m.Close()
	}

	node.Params.Remove("connecting")
	Nodes.UnRegister(nodeName)
}

func serviceInit(n *Node) error {
	if !services[n.GetType()] {
		return nil
	}

	addr := n.GetAddr()
	port, _ := n.GetPort()
	if addr == nil || port < 1 {
		return fmt.Errorf("Invalid address:%s %s %s %d", n.Name, n.GetType(), addr.String(), port)
	}

	maddr := fmt.Sprintf("%s:%d", addr.String(), port)
	log.Debug("Waiting connect to node:", maddr)
	conn, err := grpc.Dial(maddr, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Debug("Connect to node:", maddr, "[failed]", err)
		return fmt.Errorf("Connect to %s failed, error:%v", maddr, err)
	}

	log.Debug("Connect to node:", maddr, "[OK]")
	n.SetGRPCConn(conn)
	Events.Trigger(&Event{
		Name: EventNodeServiceInit,
		Node: n,
	})
	return nil
}

func stopWatcher() {
	if watcherCtrl != nil {
		close(watcherCtrl)
	}
}
