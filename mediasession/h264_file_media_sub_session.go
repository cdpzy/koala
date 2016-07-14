
// PT 	Encoding Name 	Audio/Video (A/V) 	Clock Rate (Hz) 	Channels 	Reference 
// 0	PCMU	             A	               8000	               1	    [RFC3551]
// 1	Reserved				
// 2	Reserved				
// 3	GSM                  A	               8000	               1	    [RFC3551]
// 4	G723	             A	               8000                1	    [Vineet_Kumar][RFC3551]
// 5	DVI4	             A                 8000	               1	    [RFC3551]
// 6	DVI4	             A	               16000	           1	    [RFC3551]
// 7	LPC	                 A	               8000	               1	    [RFC3551]
// 8	PCMA	             A	               8000                1	    [RFC3551]
// 9	G722	             A	               8000	               1	    [RFC3551]
// 10	L16	                 A	               44100	           2	    [RFC3551]
// 11	L16	                 A	               44100	           1	    [RFC3551]
// 12	QCELP	             A	               8000	               1	    [RFC3551]
// 13	CN	                 A	               8000	               1	    [RFC3389]
// 14	MPA	                 A	               90000		                [RFC3551][RFC2250]
// 15	G728	             A	               8000	               1	    [RFC3551]
// 16	DVI4	             A	               11025	           1	    [Joseph_Di_Pol]
// 17	DVI4	             A	               22050	           1	    [Joseph_Di_Pol]
// 18	G729	             A	               8000	               1	    [RFC3551]
// 19	Reserved             A			
// 20	Unassigned	         A			
// 21	Unassigned	         A			
// 22	Unassigned	         A			
// 23	Unassigned	         A			
// 24	Unassigned	         V			
// 25	CelB	             V	               90000		                [RFC2029]
// 26	JPEG	             V	               90000		                [RFC2435]
// 27	Unassigned	         V			
// 28	nv	                 V	               90000		                [RFC3551]
// 29	Unassigned	         V			
// 30	Unassigned	         V			
// 31	H261	             V	               90000		                [RFC4587]
// 32	MPV	                 V	               90000		                [RFC2250]
// 33	MP2T	             AV                90000		                [RFC2250]
// 34	H263	             V	               90000		                [Chunrong_Zhu]
// 35-71	Unassigned	     ?			
// 72-76	Reserved for RTCP conflict avoidance				            [RFC3551]
// 77-95	Unassigned	     ?			
// 96-127	dynamic	         ?			                                    [RFC3551]

package mediasession

import (
    "fmt"
)



type H264FileMediaSubSession struct {
    ServerMediaSubSession
    RtpPayloadType uint
    RtpTimestampFrequency uint
    RtpPayloadFormatName  string
    NumChannels           uint
    InitialPortNum        int
}

func (h264FileMediaSubSession *H264FileMediaSubSession) SDPLines() string {
    h264FileMediaSubSession.RtpPayloadType = uint(96 + h264FileMediaSubSession.Id - 1)
    sdpFmt := "m=%s %d RTP/AVP %d\r\n" +
               "c=IN IP4 %s\r\n" +
               "b=AS:%d\r\n" +
               "%s" +
               "%s" +
               "%s" +"a=control:%s\r\n"

    sdp    := fmt.Sprintf(sdpFmt, 
                          h264FileMediaSubSession.SdpMediaType(),
                          0,
                          h264FileMediaSubSession.RtpPayloadType,
                          "0.0.0.0",
                          500,
                          "",

                          "",
                          "",
                          h264FileMediaSubSession.GetTrackId(),
    )

    return sdp
}

func (h264FileMediaSubSession *H264FileMediaSubSession) SdpMediaType() string {
    return "video"
}

func (h264FileMediaSubSession *H264FileMediaSubSession) RtpmapLine() string {
    if h264FileMediaSubSession.RtpPayloadType > 96 {
        encodingParamsPart := ""
        if h264FileMediaSubSession.NumChannels != 1 {
            encodingParamsPart = fmt.Sprintf("/%d", h264FileMediaSubSession.NumChannels)
        } 
        return fmt.Sprintf("a=rtpmap:%d %s/%d%s\r\n",
                            h264FileMediaSubSession.RtpPayloadType,
                            h264FileMediaSubSession.RtpPayloadFormatName,
                            h264FileMediaSubSession.RtpTimestampFrequency,
                            encodingParamsPart,
                          )
    }
    return ""
}

func (h264FileMediaSubSession *H264FileMediaSubSession) rangeSDPLine( absStartTime,  absEndTime string) string {
    // 针对绝对时间支持
    if absStartTime != "" {
        if absEndTime != "" {
            return fmt.Sprintf("a=range:clock=%s-%s\r\n", absStartTime, absEndTime)
        } else {
            return fmt.Sprintf("a=range:clock=%s-\r\n", absStartTime)
        }
    }

    duration := h264FileMediaSubSession.duration()
    if duration == 0.0 {
        return "a=range:npt=0-\r\n"
    }

    return fmt.Sprintf("a=range:npt=0-%.3f\r\n", duration)
}

func (h264FileMediaSubSession *H264FileMediaSubSession) duration() float64 {
    return 0.0
}

func (h264FileMediaSubSession *H264FileMediaSubSession) auxSDPLine() string {
    return ""
}

func (h264FileMediaSubSession *H264FileMediaSubSession) getVPSandSPSandPPS() (vps, sps, pps string, vpsSize, spsSize, ppsSize int) {
    //setVPSandSPSandPPS
    return
}

func NewH264FileMediaSubSession( file string, multiplexRTCPWithRTP bool) *H264FileMediaSubSession{
    session := new(H264FileMediaSubSession)
    session.FileName = file
    session.RtpPayloadFormatName  = "H264"
    session.RtpTimestampFrequency = 90000
    session.NumChannels =  1

    // 是否复用原来端口
    if multiplexRTCPWithRTP {
        session.InitialPortNum = 6970
    } else {
        // 确保RTP端口为偶数
        session.InitialPortNum = (6970 + 1) &^1;
    }
    return session
}