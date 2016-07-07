package mediasession

import (
    "fmt"
)

type H264FileMediaSubSession struct {
    ServerMediaSubSession
}

func (h264FileMediaSubSession *H264FileMediaSubSession) SDPLines() string {
    sdpFmt := "m=%s %d RTP/AVP %d\r\n" +
               "c=IN IP4 %s\r\n" +
               "b=AS:%d\r\n" +
               "%s" +
               "%s" +
               "%s" +"a=control:%s\r\n"

    sdp    := fmt.Sprintf(sdpFmt, 
                          "Video",
                          0,
                          96,
                          "0.0.0.0",
                          500,
                          "",
                          "",
                          "",
                          h264FileMediaSubSession.GetTrackId(),
    )

    return sdp
}

func NewH264FileMediaSubSession( file string ) *H264FileMediaSubSession{
    session := new(H264FileMediaSubSession)
    session.FileName = file
    return session
}