package client

import (
	"sync"
)

type Params struct {
	records map[string]interface{}
	sync.RWMutex
}

func (p *Params) Set(k string, v interface{}) {
	p.Lock()
	p.records[k] = v
	p.Unlock()
}

func (p *Params) Get(k string) (v interface{}, ok bool) {
	p.RLock()
	v, ok = p.records[k]
	p.RUnlock()

	return
}

func (p *Params) Remove(k string) {
	p.Lock()
	delete(p.records, k)
	p.Unlock()
}

func (p *Params) In(k string) (ok bool) {
	p.RLock()
	_, ok = p.records[k]
	p.RUnlock()

	return
}

func (p *Params) Iterator(f func(string, interface{}) bool) {
	p.RLock()
	defer p.RUnlock()
	for k, v := range p.records {
		p.RUnlock()
		b := f(k, v)
		p.RLock()
		if !b {
			break
		}
	}
}

func NewParams() *Params {
	return &Params{records: make(map[string]interface{})}
}
