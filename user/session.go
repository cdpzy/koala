package user

import (
	"errors"
	"sync"

	"github.com/satori/go.uuid"
)

var SessionManager = NewUserSessionManager()

// UserSession 用户session
type UserSession struct {
	SessionID   string
	sessionData map[string]interface{}
	Expire      int
}

// UserSessionManager session 管理
type UserSessionManager struct {
	mux      sync.RWMutex
	sessions map[string]*UserSession
}

// NewUserSession 创建用户Session
func NewUserSession() *UserSession {
	session := new(UserSession)
	session.sessionData = make(map[string]interface{})
	session.Expire = 600
	session.RegenerateID()
	return session
}

// NewUserSessionManager 创建session管理器
func NewUserSessionManager() *UserSessionManager {
	manager := new(UserSessionManager)
	manager.sessions = make(map[string]*UserSession)
	return manager
}

// RegenerateID 重新生成
func (userSession *UserSession) RegenerateID() {
	userSession.SessionID = uuid.NewV4().String()
}

func (userSession *UserSession) Set(key string, val interface{}) {
	userSession.sessionData[key] = val
}

func (userSession *UserSession) Get(key string) (interface{}, error) {
	if data, ok := userSession.sessionData[key]; ok {
		return data, nil
	}

	return nil, errors.New("NotFound")
}

func (userSession *UserSession) Remove(key string) {
	delete(userSession.sessionData, key)
}

func (userSessionManager *UserSessionManager) Register(sid string, sess *UserSession) error {
	userSessionManager.mux.Lock()
	defer userSessionManager.mux.Unlock()

	if userSessionManager.Registered(sid) {
		return errors.New("Exists:" + sid)
	}

	userSessionManager.sessions[sid] = sess
	return nil
}

func (userSessionManager *UserSessionManager) Registered(sid string) bool {
	if _, ok := userSessionManager.sessions[sid]; ok {
		return true
	}

	return false
}

func (userSessionManager *UserSessionManager) UnRegister(sid string) {
	userSessionManager.mux.Lock()
	defer userSessionManager.mux.Unlock()

	if !userSessionManager.Registered(sid) {
		return
	}

	delete(userSessionManager.sessions, sid)
}

func (userSessionManager *UserSessionManager) Get(sid string) (*UserSession, error) {
	userSessionManager.mux.RLock()
	defer userSessionManager.mux.RUnlock()

	if sess, ok := userSessionManager.sessions[sid]; ok {
		return sess, nil
	}

	return nil, errors.New("NotFound")
}

func CreateUserSession() *UserSession {
	sess := NewUserSession()
	SessionManager.Register(sess.SessionID, sess)
	return sess
}
