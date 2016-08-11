package media

import (
	"errors"
	"sync"
)

// ServerMediaSessionManager 默认创建媒体管理器
var ServerMediaSessionManager = NewMediaSessionManager()

// MediaSessionManager 媒体session 管理器
// 用于媒体
type MediaSessionManager struct {
	mux       sync.RWMutex
	container map[string]*MediaSession
}

// MediaSession 媒体Session
type MediaSession struct {
	mux         sync.RWMutex
	Name        string
	Description string
	subsessions map[string]MediaSubSession
}

// NewMediaSessionManager 创建session 管理器
func NewMediaSessionManager() *MediaSessionManager {
	sessionManager := new(MediaSessionManager)
	sessionManager.container = make(map[string]*MediaSession)

	return sessionManager
}

// NewMediaSession media session
func NewMediaSession(name, description string) *MediaSession {
	mediaSession := new(MediaSession)
	mediaSession.Name = name
	mediaSession.Description = description
	mediaSession.subsessions = make(map[string]MediaSubSession)
	return mediaSession
}

// Register 注册管理器
func (mediaSessionManager *MediaSessionManager) Register(mediaName string, sess *MediaSession) bool {
	mediaSessionManager.mux.Lock()
	defer mediaSessionManager.mux.Unlock()

	if mediaSessionManager.Registered(mediaName) {
		return false
	}

	mediaSessionManager.container[mediaName] = sess
	return true
}

// UnRegister 取消注册
func (mediaSessionManager *MediaSessionManager) UnRegister(mediaName string) bool {
	mediaSessionManager.mux.Lock()
	defer mediaSessionManager.mux.Unlock()

	if !mediaSessionManager.Registered(mediaName) {
		return false
	}

	delete(mediaSessionManager.container, mediaName)
	return true
}

// Registered 是否已经注册
func (mediaSessionManager *MediaSessionManager) Registered(mediaName string) bool {
	if _, ok := mediaSessionManager.container[mediaName]; ok {
		return true
	}

	return false
}

// Get 获取媒体session
func (mediaSessionManager *MediaSessionManager) Get(mediaName string) (*MediaSession, error) {
	mediaSessionManager.mux.RLock()
	defer mediaSessionManager.mux.RUnlock()

	if s, ok := mediaSessionManager.container[mediaName]; ok {
		return s, nil
	}

	return nil, errors.New("NotFound")
}

// GenerateSDPDescription RTCP SDP
func (mediaSession *MediaSession) GenerateSDPDescription() string {
	return ""
}

// Duration time
func (mediaSession *MediaSession) Duration() float64 {
	var (
		minSubsessionDuration float64
		maxSubsessionDuration float64
		absStartTime          float64
		duration              float64
	)

	j := 0
	for _, sub := range mediaSession.subsessions {
		absStartTime, _ = sub.AbsoluteTimeRange()
		if absStartTime != 0.0 {
			return -1.0
		}

		duration = sub.Duration()

		if j == 0 {
			minSubsessionDuration = duration
			maxSubsessionDuration = duration
		} else if duration < minSubsessionDuration {
			minSubsessionDuration = duration
		} else if duration > maxSubsessionDuration {
			maxSubsessionDuration = duration
		}

		j++
	}

	if maxSubsessionDuration != minSubsessionDuration {
		return -maxSubsessionDuration
	}

	return maxSubsessionDuration
}

func (mediaSession *MediaSession) RegisterSubSession(mediaName string, sess MediaSubSession) bool {
	mediaSession.mux.Lock()
	defer mediaSession.mux.Unlock()

	if mediaSession.RegisteredSubSession(mediaName) {
		return false
	}

	mediaSession.subsessions[mediaName] = sess
	return true
}

func (mediaSession *MediaSession) RegisteredSubSession(mediaName string) bool {
	if _, ok := mediaSession.subsessions[mediaName]; ok {
		return true
	}

	return false
}

func (mediaSession *MediaSession) UnRegisterSubSession(mediaName string) bool {
	mediaSession.mux.Lock()
	defer mediaSession.mux.Unlock()

	if !mediaSession.RegisteredSubSession(mediaName) {
		return false
	}

	delete(mediaSession.subsessions, mediaName)
	return true
}
