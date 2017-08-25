package http

import (
	"net/url"
	"sync"
)

type RouteParams struct {
	records url.Values
	sync.RWMutex
}

func (rp *RouteParams) Get(k string) (v string) {
	rp.RLock()
	v = rp.records.Get(k)
	rp.RUnlock()
	return
}

func (rp *RouteParams) Set(k, v string) {
	rp.Lock()
	rp.records.Set(k, v)
	rp.Unlock()
}

func (rp *RouteParams) Del(k string) {
	rp.Lock()
	rp.records.Del(k)
	rp.Unlock()
}

func (rp *RouteParams) Add(k, v string) {
	rp.Lock()
	rp.records.Add(k, v)
	rp.Unlock()
}

func (rp *RouteParams) Encode() (v string) {
	rp.RLock()
	v = rp.records.Encode()
	rp.RUnlock()
	return
}

func (rp *RouteParams) Reset() {
	rp.Lock()
	rp.records = make(url.Values)
	rp.Unlock()
}

func (rp *RouteParams) Clone() url.Values {
	data := make(url.Values)
	rp.RLock()
	defer rp.RUnlock()

	for k, v := range rp.records {
		rp.RUnlock()
		data[k] = v
		rp.RLock()
	}
	return data
}

func NewRouteParams() *RouteParams {
	return &RouteParams{records: make(url.Values)}
}
