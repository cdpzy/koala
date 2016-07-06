package mediasession

import (
    "fmt"
    "sync"
    "errors"
    "strings"
)

type SubMediaSession interface{
    GetReferencecounter() int
    IncrementReferencecounter()
    DecrementReferencecounter()
}

type ServerMediaSession struct {
    mux              sync.RWMutex
    MediaName        string
    Description      string
    sdp              string
    subMediaSessions map[string]SubMediaSession
}

func (serverMediaSession *ServerMediaSession) AddSubMediaSession( name string, session SubMediaSession ) error {
    serverMediaSession.mux.Lock()
    defer serverMediaSession.mux.Unlock()

    if serverMediaSession.IsSubMediaSessionExists( name ) {
        return errors.New(fmt.Sprintf("SubMediaSession: %s exists"))
    }

    serverMediaSession.subMediaSessions[name] = session
    session.IncrementReferencecounter()
    return nil
}

func (serverMediaSession *ServerMediaSession) IsSubMediaSessionExists( name string ) bool {
    if _,ok := serverMediaSession.subMediaSessions[name]; ok {
        return true
    }

    return false
}

func (serverMediaSession *ServerMediaSession) RemoveSubMediaSession( name string ) error {
    serverMediaSession.mux.Lock()
    defer serverMediaSession.mux.Unlock()
    
    session, ok := serverMediaSession.subMediaSessions[name];
    if !ok {
        return nil
    }

    if session.GetReferencecounter() < 1 {
        delete(serverMediaSession.subMediaSessions, name)
    } else {
        session.DecrementReferencecounter()
    }

    return nil
}

func (serverMediaSession *ServerMediaSession) GenerateSDPDescription() string {

    return "test"
}

func NewServerMediaSession( name string, description string) *ServerMediaSession {
    return &ServerMediaSession{
        MediaName        : name,
        Description      : description,
        subMediaSessions : make(map[string]SubMediaSession, 0),
    }
}



type ServerMediaSessionManager struct {
    mux sync.RWMutex
    sessions map[string]*ServerMediaSession
}

func (serverMediaSessionManager *ServerMediaSessionManager) Add( name string, sess *ServerMediaSession ) error {
    serverMediaSessionManager.mux.Lock()
    defer serverMediaSessionManager.mux.Unlock()

    if serverMediaSessionManager.IsServerMediaSessionExists( name ) {
        return errors.New("ServerMediaSession exists")
    }


    return nil
}

func (serverMediaSessionManager *ServerMediaSessionManager) Get( name string ) *ServerMediaSession {
    serverMediaSessionManager.mux.RLock()
    defer serverMediaSessionManager.mux.RUnlock()

    s, ok := serverMediaSessionManager.sessions[name]
    if !ok {
        return nil
    }

    return s
}

func (serverMediaSessionManager *ServerMediaSessionManager) Remove( name string ) {
    serverMediaSessionManager.mux.Lock()
    defer serverMediaSessionManager.mux.Unlock()

    if serverMediaSessionManager.IsServerMediaSessionExists( name ){
        delete(serverMediaSessionManager.sessions, name)
    }
}

func (serverMediaSessionManager *ServerMediaSessionManager) IsServerMediaSessionExists( name string ) bool {
    if  _, ok := serverMediaSessionManager.sessions[name]; ok {
        return true
    }

    return false
}

func (serverMediaSessionManager *ServerMediaSessionManager) GetServerMediaSessions() map[string]*ServerMediaSession {
    return serverMediaSessionManager.sessions
}

func (serverMediaSessionManager *ServerMediaSessionManager) Create( name string ) (*ServerMediaSession , error) {
    serverMediaSessionManager.mux.Lock()
    defer serverMediaSessionManager.mux.Unlock()

    if s,ok := serverMediaSessionManager.sessions[name]; ok {
        return s, nil
    }

    part := strings.Split(name, ".")
    if len(part) < 1 {
        return nil, errors.New("Bad media name")
    }

    extension := part[len(part) - 1]

    var session *ServerMediaSession
    switch extension {
        case "264" :
             session = NewServerMediaSession(name, "H.246 Video")
             serverMediaSessionManager.sessions[name] = session
             return session, nil
    }

    return nil, errors.New("Bad media name")
}

func NewServerMediaSessionManager() *ServerMediaSessionManager {
    return &ServerMediaSessionManager{ sessions:make(map[string]*ServerMediaSession) }
}