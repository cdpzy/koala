package helper

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// InUint 查询V值 是否在L列表中
func InUint(v uint, l []uint) bool {
	for _, b := range l {
		if b == v {
			return true
		}
	}
	return false
}

// RandomU32 生成uint32随机数
func RandomU32() uint32 {
	return rand.Uint32()
}

// RandomU16 生成uint16随机数
func RandomU16() uint16 {
	return uint16(rand.Uint32())
}

// Random64 生成int64随机数
func Random64() int64 {
	return rand.Int63()
}
