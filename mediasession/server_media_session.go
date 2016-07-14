package mediasession

import (
    "fmt"
    "sync"
    "errors"
    "strings"
    "github.com/doublemo/koala/helper"
)

type SubMediaSession interface{
    IncrementTrackId()
    GetTrackId() string
    SDPLines()   string
}

type ServerMediaSession struct {
    mux              sync.RWMutex
    MediaName        string
    Description      string
    sdp              string
    SSM              bool   //是否为指定播放源
    IPAddr           string 
    subMediaSessions map[string]SubMediaSession
    creationTime     *helper.Time
}

func (serverMediaSession *ServerMediaSession) AddSubMediaSession( name string, session SubMediaSession ) error {
    serverMediaSession.mux.Lock()
    defer serverMediaSession.mux.Unlock()

    if serverMediaSession.IsSubMediaSessionExists( name ) {
        return errors.New(fmt.Sprintf("SubMediaSession: %s exists"))
    }

    serverMediaSession.subMediaSessions[name] = session
    session.IncrementTrackId()
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
    
    _, ok := serverMediaSession.subMediaSessions[name];
    if !ok {
        return nil
    }

    delete(serverMediaSession.subMediaSessions, name)
    return nil
}

func (serverMediaSession *ServerMediaSession) GenerateSDPDescription() string {
    var (
        sourceFilterLine string 
        rangeLine        string
    )

    if serverMediaSession.SSM {
        sourceFilterLine = fmt.Sprintf("a=source-filter: incl IN IP4 * %s\r\na=rtcp-unicast: reflection\r\n", serverMediaSession.IPAddr)
    } 

    rangeLine     = "a=range:npt=0-\r\n"

    sdpPrefixFmt := "v=0\r\n" +
                    "o=- %d%06d %d IN IP4 %s\r\n" +
                    "s=%s\r\n" +
                    "i=%s\r\n" +
                    "t=0 0\r\n" +
                    "a=tool:%s%s\r\n" +
                    "a=type:broadcast\r\n" +
                    "a=control:*\r\n" +
                    "%s" +
                    "%s" +
                    "a=x-qt-text-nam:%s\r\n" +
                    "a=x-qt-text-inf:%s\r\n" +
                    "%s"

    sdp := fmt.Sprintf( sdpPrefixFmt,
                        serverMediaSession.creationTime.SEC,
                        serverMediaSession.creationTime.USEC,
                        1,
                        serverMediaSession.IPAddr,
                        serverMediaSession.Description,
                        serverMediaSession.MediaName,
                        helper.SERVER_NAME, 
                        helper.Version,
                        sourceFilterLine,
                        rangeLine,
                        serverMediaSession.Description,
                        serverMediaSession.MediaName,
                        "")

    for _, subsession := range serverMediaSession.subMediaSessions {
        sdpLines := subsession.SDPLines()
        sdp += sdpLines
    }

    return sdp
}

func NewServerMediaSession( name string, description string) *ServerMediaSession {
    return &ServerMediaSession{
        MediaName        : name,
        Description      : description,
        subMediaSessions : make(map[string]SubMediaSession, 0),
        creationTime     : helper.GetNowTime(),
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
             session.AddSubMediaSession("H264", NewH264FileMediaSubSession(name, false))
             serverMediaSessionManager.sessions[name] = session
             return session, nil
    }

    return nil, errors.New("Bad media name")
}

func NewServerMediaSessionManager() *ServerMediaSessionManager {
    return &ServerMediaSessionManager{ sessions:make(map[string]*ServerMediaSession) }
}