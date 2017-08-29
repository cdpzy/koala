package cluster

import (
	"fmt"
	"strconv"
	"sync"
)

// Params 线程安全参数
type Params struct {
	records map[string]interface{}
	sync.RWMutex
}

// Get k string v interface{} 获取数据
func (p *Params) Get(k string) (v interface{}) {
	p.RLock()
	v = p.records[k]
	p.RUnlock()
	return
}

// Set k string v interface{} 赋值
func (p *Params) Set(k string, v interface{}) {
	p.Lock()
	p.records[k] = v
	p.Unlock()
}

func (p *Params) Exist(k string) (b bool) {
	p.RLock()
	_, b = p.records[k]
	p.RUnlock()
	return
}

func (p *Params) Remove(k string) {
	p.Lock()
	delete(p.records, k)
	p.Unlock()
}

func (p *Params) Int64(k string, def ...int64) (v int64, err error) {
	if len(def) > 0 {
		v = def[0]
	}

	m := p.Get(k)
	if m == nil {
		err = fmt.Errorf("%s is nil", k)
		return
	}

	if mv, ok := m.(string); ok {
		vv, erri := strconv.ParseInt(mv, 10, 64)
		if erri != nil {
			err = erri
			return
		}

		v = vv
	} else {
		err = fmt.Errorf("%s assert failed", k)
	}

	return
}

func (p *Params) Int32(k string, def ...int32) (v int32, err error) {
	if len(def) > 0 {
		v = def[0]
	}

	m := p.Get(k)
	if m == nil {
		err = fmt.Errorf("%s is nil", k)
		return
	}

	if mv, ok := m.(string); ok {
		vv, erri := strconv.ParseInt(mv, 10, 32)
		if erri != nil {
			err = erri
			return
		}

		v = int32(vv)
	} else {
		err = fmt.Errorf("%s assert failed", k)
	}

	return
}

func (p *Params) Int(k string, def ...int) (v int, err error) {
	if len(def) > 0 {
		v = def[0]
	}

	m := p.Get(k)
	if m == nil {
		err = fmt.Errorf("%s is nil", k)
		return
	}

	if mv, ok := m.(string); ok {
		vv, erri := strconv.Atoi(mv)
		if erri != nil {
			err = erri
			return
		}

		v = vv

	} else {
		err = fmt.Errorf("%s assert failed", k)
	}
	return
}

func (p *Params) String(k string, def ...string) (v string) {
	if len(def) > 0 {
		v = def[0]
	}

	m := p.Get(k)
	if m == nil {
		return
	}

	if mv, ok := m.(string); ok {
		v = mv
	}

	return
}

func (p *Params) Iterator(f func(k string, v interface{}) bool) {
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
