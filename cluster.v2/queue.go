package cluster

import "container/heap"

// Queue 节点队列
type Queue []*Node

func (q Queue) Len() int { return len(q) }

func (q Queue) Less(i, j int) bool {
	priority1, _ := q[i].GetPriority()
	priority2, _ := q[j].GetPriority()
	return priority1 > priority2
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].idx = i
	q[j].idx = j
}

func (q *Queue) Push(x interface{}) {
	n := len(*q)
	node := x.(*Node)
	node.idx = n
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
	heap.Fix(q, node.idx)
}
