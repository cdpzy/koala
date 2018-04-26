package cluster

// NodeStatus 节点状态类型
type NodeStatus int

const (
	// NodeStatusOK node providing service
	NodeStatusOK NodeStatus = iota
	// NodeStatusClosed node is down
	NodeStatusClosed
)
