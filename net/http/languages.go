package http

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// AcceptLanguage is a single language from the Accept-Language HTTP header.
type AcceptLanguage struct {
	Language string
	Quality  float32
}

// AcceptLanguages is collection of sortable AcceptLanguage instances.
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

// ResolveAcceptLanguage returns a sorted list of Accept-Language
// header values.
//
// The results are sorted using the quality defined in the header for each
// language range with the most qualified language range as the first
// element in the slice.
//
// See the HTTP header fields specification
// (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) for more details.
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
				log.Errorf("Detected malformed Accept-Language header quality in '%s', assuming quality is 1", languageRange)
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
