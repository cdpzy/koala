package mongodb

import (
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var (
	_defautls_sm *SessionManager
)

type Session struct {
	Session *mgo.Session
	Db      string
}

type SessionManager struct {
	sync.RWMutex
	records map[string]*Session
}

func (sm *SessionManager) Register(k string, s *Session) {
	sm.Lock()
	sm.records[k] = s
	sm.Unlock()
}

func (sm *SessionManager) Unregister(k string) {
	sm.Lock()
	defer sm.Unlock()

	if m, ok := sm.records[k]; ok {
		sm.Unlock()
		m.Session.Close()
		sm.Lock()
	}

	delete(sm.records, k)
}

// UnregisterAll unregister all at one batch
func (sm *SessionManager) UnregisterAll() {
	sm.Lock()
	defer sm.Unlock()

	for _, m := range sm.records {
		m.Session.Close()
	}

	sm.records = make(map[string]*Session)
}

func (sm *SessionManager) Get(k string) (s *Session) {
	sm.RLock()
	s = sm.records[k]
	sm.RUnlock()
	return
}

func NewSessionManager() *SessionManager {
	return &SessionManager{records: make(map[string]*Session)}
}

func init() {
	_defautls_sm = NewSessionManager()
}

func Create(key, addr, name, user, pass string) (*Session, error) {
	session, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	s := &Session{Session: session, Db: name}
	_defautls_sm.Register(key, s)
	return s, nil
}

func Get(key string) *Session {
	s := _defautls_sm.Get(key)
	if s == nil {
		return nil
	}
	return s
}

func Close(key string) {
	_defautls_sm.Unregister(key)
}

// CloseAll close all db that allready registered
func CloseAll() {
	_defautls_sm.UnregisterAll()
}

func GetCopySession(key string) *mgo.Session {
	s := Get(key)
	if s == nil {
		return nil
	}

	return s.Session.Copy()
}

func GetDb(key string) string {
	s := Get(key)
	if s == nil {
		return ""
	}
	return s.Db
}
