package cluster

import (
	"net"
	"strconv"

	"google.golang.org/grpc"
)

const (
	cNodeType       = "Type"
	cNodeAddr       = "Addr"
	cNodePort       = "Port"
	cNodeStatus     = "Status"
	cNodePriority   = "Priority"
	cNodeConnection = "GRPCConn"
)

// Node each node stands for one grpc server
type Node struct {
	Name   string
	Params *Params
	idx    int
}

// GetType self defined node type, each type stands for one type of service
func (n *Node) GetType() string {
	return n.Params.String(cNodeType)
}

// SetType set service type of this node
func (n *Node) SetType(v string) {
	n.Params.Set(cNodeType, v)
}

// GetAddr address of this service
func (n *Node) GetAddr() net.IP {
	return net.ParseIP(n.Params.String(cNodeAddr))
}

// SetAddr set address of this service
func (n *Node) SetAddr(v net.IP) {
	n.Params.Set(cNodeAddr, v.String())
}

// GetPort port of this service
func (n *Node) GetPort() (int, error) {
	return n.Params.Int(cNodePort)
}

// SetPort port of this service
func (n *Node) SetPort(v int) {
	n.Params.Set(cNodePort, strconv.FormatInt(int64(v), 10))
}

// GetPriority priority of this node, used of load balance
func (n *Node) GetPriority() (int, error) {
	return n.Params.Int(cNodePriority)
}

// SetPriority priority of this node, used of load balance
func (n *Node) SetPriority(v int) {
	n.Params.Set(cNodePriority, strconv.FormatInt(int64(v), 10))
}

// GetGRPCConn connected onece been discovered until node down.
func (n *Node) GetGRPCConn() *grpc.ClientConn {
	if n.Params == nil {
		return nil
	}

	conn := n.Params.Get(cNodeConnection)
	if conn == nil {
		return nil
	}

	if m, ok := conn.(*grpc.ClientConn); ok {
		return m
	}

	return nil
}

// SetGRPCConn connected onece been discovered until node down.
func (n *Node) SetGRPCConn(conn *grpc.ClientConn) {
	n.Params.Set(cNodeConnection, conn)
}

// RemoveGRPCConn connected onece been discovered until node down.
func (n *Node) RemoveGRPCConn() {
	n.Params.Remove(cNodeConnection)
}

// GetStatus node service status
func (n *Node) GetStatus() NodeStatus {
	v, err := n.Params.Int(cNodeStatus)
	if err != nil {
		return NodeStatusClosed
	}

	return NodeStatus(v)
}

// SetStatus node service status
func (n *Node) SetStatus(v NodeStatus) {
	n.Params.Set(cNodeStatus, strconv.FormatInt(int64(v), 10))
}

// NewNode create one new node. NOTE: name should be uniq.
func NewNode(name string) *Node {
	return &Node{Name: name, Params: NewParams()}
}
