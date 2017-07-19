package client

import (
	"sync"
)

type EvtCallBack func(*Client)

type Evt struct {
	records map[string]EvtCallBack
	sync.RWMutex
}

func (e *Evt) Set(k string, v EvtCallBack) {
	e.Lock()
	e.records[k] = v
	e.Unlock()
}

func (e *Evt) Get(k string) (v interface{}, ok bool) {
	e.RLock()
	v, ok = e.records[k]
	e.RUnlock()

	return
}

func (e *Evt) Remove(k string) {
	e.Lock()
	delete(e.records, k)
	e.Unlock()
}

func (e *Evt) Iterator(f func(string, EvtCallBack) bool) {
	e.RLock()
	defer e.RUnlock()
	for k, v := range e.records {
		e.RUnlock()
		b := f(k, v)
		e.RLock()
		if !b {
			break
		}
	}
}

func (e *Evt) Count() (count int) {
	e.RLock()
	count = len(e.records)
	e.RUnlock()
	return
}

func NewEvt() *Evt {
	return &Evt{records: make(map[string]EvtCallBack)}
}
