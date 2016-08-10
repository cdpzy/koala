package mediasession


const (
    VPS_MAX_SIZE uint = 1000
    SPS_MAX_SIZE      = 1000
    SEI_MAX_SIZE      = 5000
   
)
type H264VideoStreamParser struct {
    StreamParser
    haveSeenFirstStartCode bool
    outputStartCodeSize    int
}

func (h264VideoStreamParser *H264VideoStreamParser) Parse() uint {
    if !h264VideoStreamParser.haveSeenFirstStartCode {
        for first4Bytes := h264VideoStreamParser.test4Bytes(); first4Bytes != 0x00000001; {
            h264VideoStreamParser.Get1Bytes()
            h264VideoStreamParser.saveParserState()
        }

        h264VideoStreamParser.skipBytes(4)
        h264VideoStreamParser.haveSeenFirstStartCode = true
    }

    if h264VideoStreamParser.outputStartCodeSize > 0 && !h264VideoStreamParser.HaveSeenEOF() {
        //h264VideoStreamParser.sav(0x00000001)
    }

    return 0
}

func NewH264VideoStreamParser() *H264VideoStreamParser {
    return &H264VideoStreamParser{}
}