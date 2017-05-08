package lua

import (
	glua "github.com/yuin/gopher-lua"
)

var pool *LStatePool

func init() {
	pool = NewLStatePool(LStatePoolDefaultNumber)
}

func LStateGet() (*glua.LState, error) {
	return pool.Get()
}

func LStatePut(L *glua.LState) {
	pool.Put(L)
}

func LStateLimit(limit int) {
	pool.SetLimit(limit)
}

func LStateShutdown() {
	pool.Shutdown()
}
