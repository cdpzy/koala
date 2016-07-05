package koala

import (
	"io"
	"io/ioutil"
    "strconv"
    "sort"
    "fmt"
    "bytes"
    "log"
    "strings"
    "net/http"
	"net/url"
)

type HTTPRequest struct {
	BaseRequest
    *http.Request
    ContentType     string
    Format          string 
    AcceptLanguages AcceptLanguages
    Locale          string
}

func (httpRequest *HTTPRequest) GetMethod() string {
    return httpRequest.Method
}

func (httpRequest *HTTPRequest) GetHeader() http.Header {
	return httpRequest.Header
}

func (httpRequest *HTTPRequest) GetBody() io.ReadCloser {
	return httpRequest.Body
}


func (httpRequest *HTTPRequest) GetURL() *url.URL {
	return httpRequest.URL
}


func (httpRequest *HTTPRequest) String() string {
    s := fmt.Sprintf("%s %s %s/%d.%d\r\n", httpRequest.Method, httpRequest.URL, httpRequest.Proto, httpRequest.ProtoMajor, httpRequest.ProtoMinor)
    for k, v := range httpRequest.Header {
        for _, v := range v {
            s += fmt.Sprintf("%s: %s\r\n", k, v)
        }
    }
    s += "\r\n"
    if httpRequest.Body != nil {
        str, _ := ioutil.ReadAll(httpRequest.Body)
        s += string(str)
    }
    return s
}

func NewHTTPRequest( r *http.Request ) *HTTPRequest {
    req := new(HTTPRequest)
    req.Request = r
    req.ContentType     = ResolveContentType( r )
    req.Format          = ResolveFormat( r )
    req.AcceptLanguages = ResolveAcceptLanguage( r )
    return req
}

func ResolveContentType( req *http.Request ) string {
    contentType := req.Header.Get("Content-Type")
    if contentType == "" {
        return "text/html"
    }

    return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

func ResolveFormat(req *http.Request) string {
    accept := req.Header.Get("accept")

    switch {
	case accept == "",
		strings.HasPrefix(accept, "*/*"), // */
		strings.Contains(accept, "application/xhtml"),
		strings.Contains(accept, "text/html"):
		return "html"
	case strings.Contains(accept, "application/json"),
		strings.Contains(accept, "text/javascript"),
		strings.Contains(accept, "application/javascript"):
		return "json"
	case strings.Contains(accept, "application/xml"),
		strings.Contains(accept, "text/xml"):
		return "xml"
	case strings.Contains(accept, "text/plain"):
		return "txt"
	}

	return "html"
}

type AcceptLanguage struct {
	Language string
	Quality  float32
}

type AcceptLanguages []AcceptLanguage

func (al AcceptLanguages) Len() int           { return len(al) }
func (al AcceptLanguages) Swap(i, j int)      { al[i], al[j] = al[j], al[i] }
func (al AcceptLanguages) Less(i, j int) bool { return al[i].Quality > al[j].Quality }
func (al AcceptLanguages) String() string {
	output := bytes.NewBufferString("")
	for i, language := range al {
		output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
		if i != len(al)-1 {
			output.WriteString(", ")
		}
	}
	return output.String()
}

func ResolveAcceptLanguage(req *http.Request) AcceptLanguages {
	header := req.Header.Get("Accept-Language")
	if header == "" {
		return nil
	}

	acceptLanguageHeaderValues := strings.Split(header, ",")
	acceptLanguages := make(AcceptLanguages, len(acceptLanguageHeaderValues))

	for i, languageRange := range acceptLanguageHeaderValues {
		if qualifiedRange := strings.Split(languageRange, ";q="); len(qualifiedRange) == 2 {
			quality, error := strconv.ParseFloat(qualifiedRange[1], 32)
			if error != nil {
				log.Printf("Detected malformed Accept-Language header quality in '%s', assuming quality is 1", languageRange)
				acceptLanguages[i] = AcceptLanguage{qualifiedRange[0], 1}
			} else {
				acceptLanguages[i] = AcceptLanguage{qualifiedRange[0], float32(quality)}
			}
		} else {
			acceptLanguages[i] = AcceptLanguage{languageRange, 1}
		}
	}

	sort.Sort(acceptLanguages)
	return acceptLanguages
}
