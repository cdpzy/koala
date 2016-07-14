package mediasession


const (
    VPS_MAX_SIZE uint = 1000
    SPS_MAX_SIZE      = 1000
    SEI_MAX_SIZE      = 5000
   
)
type H264VideoStreamParser struct {
    StreamParser
}

func (h264VideoStreamParser *H264VideoStreamParser) Parse() {}

func NewH264VideoStreamParser() *H264VideoStreamParser {
    return &H264VideoStreamParser{}
}