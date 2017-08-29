package cluster

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// Node 节点
type Node struct {
	Name   string
	Params *Params
	idx    int
}

// GetType 获取节点类型
func (n *Node) GetType() string {
	return n.Params.String("Type")
}

func (n *Node) SetType(v string) {
	n.Params.Set("Type", v)
}

// GetType 获取节点地址
func (n *Node) GetAddr() net.IP {
	return net.ParseIP(n.Params.String("Addr"))
}

func (n *Node) SetAddr(v net.IP) {
	n.Params.Set("Addr", v.String())
}

// GetType 获取节点端口
func (n *Node) GetPort() (int, error) {
	return n.Params.Int("Port")
}

func (n *Node) SetPort(v int) {
	n.Params.Set("Port", fmt.Sprint(v))
}

// GetType 获取节点策略
func (n *Node) GetPriority() (int, error) {
	return n.Params.Int("Priority")
}

func (n *Node) SetPriority(v int) {
	n.Params.Set("Priority", fmt.Sprint(v))
}

// GetGRPCConn GPRC节点连接
func (n *Node) GetGRPCConn() *grpc.ClientConn {
	conn := n.Params.Get("GRPCConn")
	if conn == nil {
		return nil
	}

	if m, ok := conn.(*grpc.ClientConn); ok {
		return m
	}

	return nil
}

func (n *Node) SetGRPCConn(conn *grpc.ClientConn) {
	n.Params.Set("GRPCConn", conn)
}

func (n *Node) RemoveGRPCConn() {
	n.Params.Remove("GRPCConn")
}

// SetStatus 状态
func (n *Node) GetStatus() NodeStatus {
	v, err := n.Params.Int("Status")
	if err != nil {
		return NodeStatusClosed
	}

	return NodeStatus(v)
}

func (n *Node) SetStatus(v NodeStatus) {
	n.Params.Set("Status", fmt.Sprint(v))
}

func NewNode(name string) *Node {
	return &Node{Name: name, Params: NewParams()}
}
