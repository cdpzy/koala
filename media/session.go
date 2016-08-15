package media

import (
	"errors"
	"fmt"
	"sync"

	"github.com/doublemo/koala/helper"
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
	mux               sync.RWMutex
	Name              string
	Description       string
	subsessions       map[string]MediaSubSession
	IsSSM             bool   // Source Specific Multicast (SSM) Sessions
	MiscSDPLines      string // miscellaneous session SDP lines (if any)
	subsessionCounter int
	CreateAt          *helper.Time
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
	mediaSession.subsessionCounter = 0
	mediaSession.CreateAt = helper.GetNowTime()
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

// GenerateSDP RTCP SDP
func (mediaSession *MediaSession) GenerateSDP() string {
	var (
		sourceFilterLine string
		rangeLine        string
	)

	// For a SSM sessions, we need a "a=source-filter: incl ..." line also:
	// eg: a=source-filter:incl IN IP4 233.252.0.2 198.51.100.1
	if mediaSession.IsSSM {
		sourceFilterLine = fmt.Sprintf("a=source-filter: incl IN IP4 * %s\r\na=rtcp-unicast: reflection\r\n", "192.168.18.152")
	}

	duration := mediaSession.Duration()
	if duration == 0.0 {
		rangeLine = "a=range:npt=0-\r\n"
	} else if duration > 0.0 {
		rangeLine = fmt.Sprintf("a=range:npt=0-%.3f\r\n", duration)
	}

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

	sdp := fmt.Sprintf(sdpPrefixFmt,
		mediaSession.CreateAt.SEC,
		mediaSession.CreateAt.USEC,
		1,
		"192.168.18.152",
		mediaSession.Description,
		mediaSession.Name,
		helper.SERVER_NAME,
		helper.Version,
		sourceFilterLine,
		rangeLine,
		mediaSession.Description,
		mediaSession.Name,
		mediaSession.MiscSDPLines)

	for _, sub := range mediaSession.subsessions {
		sub.SetParentDuration(duration)
		sdp += sub.GenerateSDP()
	}

	return sdp
}

// Duration time
func (mediaSession *MediaSession) Duration() float64 {
	var (
		minSubsessionDuration float64
		maxSubsessionDuration float64
		absStartTime          string
		duration              float64
	)

	j := 0
	for _, sub := range mediaSession.subsessions {
		absStartTime, _ = sub.AbsoluteTimeRange()
		if absStartTime != "" {
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

	mediaSession.subsessionCounter++
	sess.SetID(mediaSession.subsessionCounter)
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
