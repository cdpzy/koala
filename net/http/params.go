package http

import (
	"mime/multipart"
	"net/url"
	"os"
	"reflect"

	log "github.com/Sirupsen/logrus"
)

// Params 参数
type Params struct {
	url.Values
	Fixed    url.Values                         // Fixed parameters from the route, e.g. App.Action("fixed param")
	Route    url.Values                         // Parameters extracted from the route,  e.g. /customers/{id}
	Query    url.Values                         // Parameters from the query string, e.g. /index?limit=10
	Form     url.Values                         // Parameters from the request body.
	Files    map[string][]*multipart.FileHeader // Files uploaded in a multipart form
	tmpFiles []*os.File                         // Temp files used during the request.
}

// Parse 解析
func (params *Params) Parse(req *Request) {
	params.Query = req.URL.Query()

	switch req.ContentType {
	case "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			log.Errorln("Error parsing request body:", err)
		} else {
			params.Form = req.Form
		}

	case "multipart/form-data":
		if err := req.ParseMultipartForm(32 << 20 /* 32 MB */); err != nil {
			log.Errorln("Error parsing request body:", err)
		} else {
			params.Form = req.MultipartForm.Value
			params.Files = req.MultipartForm.File
		}
	}
}

func (p *Params) Bind(dest interface{}, name string) {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		panic("params: non-pointer passed to Bind: " + name)
	}
	value = value.Elem()
	if !value.CanSet() {
		panic("params: non-settable variable passed to Bind: " + name)
	}
	value.Set(Bind(p, name, value.Type()))
}

// calcValues 计算参数内容
func (params *Params) calcValues() url.Values {
	numParams := len(params.Query) + len(params.Fixed) + len(params.Route) + len(params.Form)
	if numParams == 0 {
		return make(url.Values, 0)
	}

	switch numParams {
	case len(params.Query):
		return params.Query
	case len(params.Route):
		return params.Route
	case len(params.Fixed):
		return params.Fixed
	case len(params.Form):
		return params.Form
	}

	values := make(url.Values, numParams)
	for k, v := range params.Fixed {
		values[k] = append(values[k], v...)
	}
	for k, v := range params.Query {
		values[k] = append(values[k], v...)
	}
	for k, v := range params.Route {
		values[k] = append(values[k], v...)
	}
	for k, v := range params.Form {
		values[k] = append(values[k], v...)
	}

	return values
}
