package media

import (
	"fmt"
)

// MediaSubSession 针对不同的媒体格式，指定特定的session
type MediaSubSession interface {
	TrackId() string
	AbsoluteTimeRange() (float64, float64)
	Duration() float64
}

// BaseMediaSubSession 基础session
type BaseMediaSubSession struct {
	ID       int
	FileName string
	FileSize string
}

func (baseMediaSubSession *BaseMediaSubSession) TrackId() string {
	return fmt.Sprintf("track%d", baseMediaSubSession.ID)
}
