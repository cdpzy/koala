package kafka

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

// C 数据结构定义
type C struct {
	Type       string      `json:"type"`
	InstanceID string      `json:"instance_id"`
	Table      string      `json:"table"`
	Host       string      `json:"host"`
	Key        string      `json:"key"`
	CreatedAt  int64       `json:"created_at"`
	Data       interface{} `json:"data"`
}

func (c *C) Encode() []byte {
	str, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return sarama.ByteEncoder(str)
}

func (c *C) Decode(b []byte) error {
	return json.Unmarshal(b, c)
}

func NewC() *C {
	host, _ := os.Hostname()
	return &C{
		Host:      host,
		CreatedAt: time.Now().Unix(),
	}
}
