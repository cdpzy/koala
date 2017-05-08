package lua

import (
	"errors"
	"sync"

	glua "github.com/yuin/gopher-lua"
)

const (
	LStatePoolDefaultNumber int = 1024
)

var (
	ErrorLStatePoolMaxLimit = errors.New("ErrorLStatePoolMaxLimit") // 线程已经到达最大上线
)

// LStatePool LState池
type LStatePool struct {
	mute    sync.Mutex
	saved   []*glua.LState
	limit   int
	counter int
}

// Get pool
func (p *LStatePool) Get() (*glua.LState, error) {
	p.mute.Lock()
	defer p.mute.Unlock()

	n := len(p.saved)
	if n == 0 {
		return p.New()
	}

	x := p.saved[n-1]
	p.saved = p.saved[0 : n-1]
	return x, nil
}

// New Pool
func (p *LStatePool) New() (*glua.LState, error) {
	if p.counter >= p.limit {
		return nil, ErrorLStatePoolMaxLimit
	}

	L := glua.NewState()
	p.counter++
	return L, nil
}

// Put save
func (p *LStatePool) Put(L *glua.LState) {
	p.mute.Lock()
	p.saved = append(p.saved, L)
	p.mute.Unlock()
}

// Shutdown close all
func (p *LStatePool) Shutdown() {
	for _, L := range p.saved {
		L.Close()
	}

	p.counter = 0
	p.saved = make([]*glua.LState, 0)
}

func (p *LStatePool) SetLimit(v int) {
	p.limit = v
}

func NewLStatePool(limit int) *LStatePool {
	return &LStatePool{
		saved: make([]*glua.LState, 0),
		limit: limit,
	}
}
