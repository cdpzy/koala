package koala

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/doublemo/koala/media"
	"github.com/doublemo/koala/protocol/rtp"
)

const AllowedMethod = "OPTIONS, DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, GET_PARAMETER, SET_PARAMETER"

type HandleMethod struct {
	r Request
	w Response
}

// http
func (handleMethod *HandleMethod) GET() {
	header := handleMethod.w.GetHeader()
	csep := handleMethod.r.GetHeader().Get("CSeq")
	header.Set("CSeq", csep)
	header.Set("Date", time.Now().String())
	header.Set("Cache-Control", "no-cache")
	header.Set("Pragma", "no-cache")
	header.Set("Content-Type", "application/x-rtsp-tunnelled")
	handleMethod.w.Write("")
}

func (handleMethod *HandleMethod) POST() {
	log.Println(handleMethod.w.GetHeader())
}

func (handleMethod *HandleMethod) OPTIONS() {
	csep := handleMethod.r.GetHeader().Get("CSeq")
	header := handleMethod.w.GetHeader()
	header.Set("CSeq", csep)
	header.Set("Date", time.Now().String())
	header.Set("Public", AllowedMethod)
	handleMethod.w.Write("")
}

func (handleMethod *HandleMethod) DESCRIBE() {
	header := handleMethod.w.GetHeader()
	path := strings.Trim(handleMethod.r.GetURL().Path, "/")
	csep := handleMethod.r.GetHeader().Get("CSeq")

	session := media.NewMediaSession(path, "H264")
	h264MediaSession := media.NewH264MediaSubSession(path)
	h264MediaSession.SetMultiplexRTCPWithRTP(false)
	h264MediaSession.SetInitialPort(6970)
	session.RegisterSubSession("H264", h264MediaSession)
	media.ServerMediaSessionManager.Register(path, session)

	sdpDescription := session.GenerateSDP()
	sdpDescriptionSize := len(sdpDescription)
	if sdpDescriptionSize < 1 {
		handleMethod.w.NotFound()
		return
	}

	header.Set("CSeq", csep)
	header.Set("Date", time.Now().String())
	header.Set("Content-Base", handleMethod.r.GetURL().String())
	header.Set("Content-Type", "application/sdp")
	header.Set("Content-Length", strconv.Itoa(sdpDescriptionSize))
	handleMethod.w.Write(sdpDescription)
}

func (handleMethod *HandleMethod) SETUP() {
	header := handleMethod.w.GetHeader()
	csep := handleMethod.r.GetHeader().Get("CSeq")
	path := strings.Trim(handleMethod.r.GetURL().Path, "/")
	paths := strings.Split(path, "/")
	if len(paths) != 2 {
		handleMethod.w.NotFound()
		return
	}

	name := paths[0]
	trackId := paths[1]

	session, err := media.ServerMediaSessionManager.Get(name)
	if err != nil {
		handleMethod.w.NotFound()
		return
	}

	transport := rtp.ParseTransportHeader(handleMethod.r.GetHeader().Get("transport"))
	if transport.StreamingMode == rtp.RTP_TCP && transport.RTPChannelID == 0xFF {
		transport.StreamingMode = rtp.RTP_TCP
		transport.RTPChannelID = 1
		transport.RTCPChannelID = 2
	}

	if transport.DestinationAddr == "" {
		ip := strings.Split(handleMethod.r.GetRemoteAddr().String(), ":")
		transport.DestinationAddr = ip[0]
	}

	parameters, err := session.GetStreamParameters(transport, trackId)
	if err != nil {
		handleMethod.w.NotFound()
		return
	}

	sip := strings.Split(handleMethod.r.GetLocalAddr().String(), ":")
	sourceSev := sip[0]

	//rangeHeader, err := rtp.ParseRangeHeader(handleMethod.r.GetHeader().Get("range"))
	//if err == nil {

	//}
	//x-playNow

	header.Set("CSeq", csep)
	header.Set("Date", time.Now().String())
	fmt.Println("transport.StreamingMode =", transport.DestinationAddr)
	if parameters.IsMulticast {
		switch transport.StreamingMode {
		case rtp.RTP_UDP:
			transportString := fmt.Sprintf("RTP/AVP;multicast;destination=%s;source=%s;port=%d-%d;ttl=%d", parameters.DestinationAddr, sourceSev, parameters.ServerRTPPort, parameters.ServerRTCPPort, parameters.DestinationTTL)
			header.Set("Transport", transportString)
		case rtp.RAW_UDP:
			transportString := fmt.Sprintf("%s;multicast;destination=%s;source=%s;port=%d;ttl=%d", transport.StreamingModeName, parameters.DestinationAddr, sourceSev, parameters.ServerRTPPort, parameters.DestinationTTL)
			header.Set("Transport", transportString)
		default:
			handleMethod.w.NotSupported("Streaming Mode RTP_TCP")
			return
		}
	} else {
		switch transport.StreamingMode {
		case rtp.RTP_UDP:
			transportString := fmt.Sprintf("RTP/AVP;unicast;destination=%s;source=%s;client_port=%d-%d;server_port=%d-%d", parameters.DestinationAddr, sourceSev, parameters.ClientRTPPort, parameters.ClientRTCPPort, parameters.ServerRTPPort, parameters.ServerRTCPPort)
			header.Set("Transport", transportString)

		case rtp.RTP_TCP:
			transportString := fmt.Sprintf("RTP/AVP/TCP;unicast;destination=%s;source=%s;interleaved=%d-%d\r\n", parameters.DestinationAddr, sourceSev, transport.RTPChannelID, transport.RTCPChannelID)
			header.Set("Transport", transportString)

		case rtp.RAW_UDP:
			transportString := fmt.Sprintf("%s;unicast;destination=%s;source=%s;client_port=%d;server_port=%d", transport.StreamingModeName, parameters.DestinationAddr, sourceSev, parameters.ClientRTPPort, parameters.ServerRTPPort)
			header.Set("Transport", transportString)

		}
	}

	header.Set("Session", "45564666ffddf4rr;timeout=600")
	handleMethod.w.Write("")
}

func (handleMethod *HandleMethod) PLAY() {
	fmt.Println("PLAY")
}

func NewHandleMethod(r Request, w Response) *HandleMethod {
	handle := new(HandleMethod)
	handle.r = r
	handle.w = w
	return handle
}
