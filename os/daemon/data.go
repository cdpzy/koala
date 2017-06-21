package daemon

import (
	"encoding/json"
)

// Data 数据
type Data struct {
	Pid      int    // 守护进程ID
	Command  int16  //
	NodeName string // 节点名称
	NodeAddr string // 节点服务监听地址
	NodePort int    // 节点服务监听端口
	NodeType string // 节点类型
	DateUinx int64  // 时间
}

func (d *Data) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode 读取
func (d *Data) Decode(b []byte) error {
	return json.Unmarshal(b, d)
}
