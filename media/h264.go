package media

// H264MediaSubSession H264
type H264MediaSubSession struct {
	BaseMediaSubSession
}

// NewH264MediaSubSession H264
func NewH264MediaSubSession(FileName string) *H264MediaSubSession {
	sess := new(H264MediaSubSession)
	sess.FileName = FileName
	return sess
}

func (h264MediaSubSession *H264MediaSubSession) AbsoluteTimeRange() (float64, float64) {
	return 0.0, 0.0
}

func (h264MediaSubSession *H264MediaSubSession) Duration() float64 {
	return 0
}
