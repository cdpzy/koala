package cluster

import "container/heap"

// Queue 节点队列
type Queue []*Node

func (q Queue) Len() int { return len(q) }

func (q Queue) Less(i, j int) bool {
	return q[i].Priority > q[j].Priority
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Idx = i
	q[j].Idx = j
}

func (q *Queue) Push(x interface{}) {
	n := len(*q)
	node := x.(*Node)
	node.Idx = n
	*q = append(*q, node)
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	node := old[n-1]
	*q = old[0 : n-1]
	return node
}

func (q *Queue) Change(node *Node) {
	heap.Fix(q, node.Idx)
}
