package cluster

// NodeStatus 节点状态类型
type NodeStatus int

const (
	NodeStatusOK     NodeStatus = iota // 正常节点
	NodeStatusClosed                   // 节点关闭
)
