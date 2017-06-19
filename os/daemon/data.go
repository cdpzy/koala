package daemon

import (
	"bytes"
	"encoding/binary"
	"time"
)

// Data 数据
type Data struct {
	Pid      int32  // 守护进程ID
	Command  int16  //
	NodeName string // 节点名称
	NodeAddr string // 节点服务监听地址
	NodePort int16  // 节点服务监听端口
	NodeType string // 节点类型
	DateUinx int64  // 时间
}

func (d *Data) Encode() []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, d.Pid)
	binary.Write(&buffer, binary.LittleEndian, d.Command)
	binary.Write(&buffer, binary.LittleEndian, time.Now().Unix())

	nameBytes := []byte(d.NodeName)
	binary.Write(&buffer, binary.LittleEndian, int16(len(nameBytes)))
	binary.Write(&buffer, binary.LittleEndian, nameBytes)

	addrBytes := []byte(d.NodeAddr)
	binary.Write(&buffer, binary.LittleEndian, int16(len(addrBytes)))
	binary.Write(&buffer, binary.LittleEndian, addrBytes)

	binary.Write(&buffer, binary.LittleEndian, d.NodePort)

	typeBytes := []byte(d.NodeType)
	binary.Write(&buffer, binary.LittleEndian, int16(len(typeBytes)))
	binary.Write(&buffer, binary.LittleEndian, typeBytes)

	return buffer.Bytes()
}

// Decode 读取
func (d *Data) Decode(b []byte) {
	reader := bytes.NewReader(b)
	binary.Read(reader, binary.LittleEndian, &d.Pid)
	binary.Read(reader, binary.LittleEndian, &d.Command)
	binary.Read(reader, binary.LittleEndian, &d.DateUinx)

	var size int16
	binary.Read(reader, binary.LittleEndian, &size)
	nameBytes := make([]byte, size)
	binary.Read(reader, binary.LittleEndian, &nameBytes)
	d.NodeName = string(nameBytes)

	binary.Read(reader, binary.LittleEndian, &size)
	addrBytes := make([]byte, size)
	binary.Read(reader, binary.LittleEndian, &addrBytes)
	d.NodeAddr = string(addrBytes)

	binary.Read(reader, binary.LittleEndian, &d.NodePort)

	binary.Read(reader, binary.LittleEndian, &size)
	typeBytes := make([]byte, size)
	binary.Read(reader, binary.LittleEndian, &typeBytes)
	d.NodeType = string(typeBytes)
}
