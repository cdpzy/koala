package mongodb

import (
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var (
	_defautls_sm *SessionManager
)

type SessionManager struct {
	sync.RWMutex
	records map[string]*mgo.Session
}

func (sm *SessionManager) Register(k string, s *mgo.Session) {
	sm.Lock()
	sm.records[k] = s
	sm.Unlock()
}

func (sm *SessionManager) Unregister(k string) {
	sm.Lock()
	defer sm.Unlock()

	if m, ok := sm.records[k]; ok {
		sm.Unlock()
		m.Close()
		sm.Lock()
	}

	delete(sm.records, k)
}

func (sm *SessionManager) Get(k string) (s *mgo.Session) {
	sm.RLock()
	s = sm.records[k]
	sm.RUnlock()
	return
}

func NewSessionManager() *SessionManager {
	return &SessionManager{records: make(map[string]*mgo.Session)}
}

func init() {
	_defautls_sm = NewSessionManager()
}

func Create(key, addr, name, user, pass string) (*mgo.Session, error) {
	session, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}

	_defautls_sm.Register(key, session)
	return session, nil
}

func Get(key string) *mgo.Session {
	s := _defautls_sm.Get(key)
	if s == nil {
		return nil
	}
	return s.Copy()
}

func Close(key string) {
	_defautls_sm.Unregister(key)
}
