package http

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// RenderJSONResult JSON
type RenderJSONResult struct {
	obj      interface{}
	callback string
	Pretty   bool
}

func (r RenderJSONResult) Apply(req *Request, resp *Response) {
	var b []byte
	var err error
	if r.Pretty {
		b, err = json.MarshalIndent(r.obj, "", "  ")
	} else {
		b, err = json.Marshal(r.obj)
	}

	if err != nil {
		ErrorResult{Error: err}.Apply(req, resp)
		return
	}

	if r.callback == "" {
		resp.WriteHeader(http.StatusOK, "application/json; charset=utf-8")
		if _, err = resp.Out.Write(b); err != nil {
			log.Errorln("Response write failed:", err)
		}
		return
	}

	resp.WriteHeader(http.StatusOK, "application/javascript; charset=utf-8")
	if _, err = resp.Out.Write([]byte(r.callback + "(")); err != nil {
		log.Errorln("Response write failed:", err)
	}
	if _, err = resp.Out.Write(b); err != nil {
		log.Errorln("Response write failed:", err)
	}
	if _, err = resp.Out.Write([]byte(");")); err != nil {
		log.Errorln("Response write failed:", err)
	}
}

type RenderJSONTextResult struct {
	obj      string
	callback string
}

func (r RenderJSONTextResult) Apply(req *Request, resp *Response) {
	var err error
	if err != nil {
		ErrorResult{Error: err}.Apply(req, resp)
		return
	}

	if r.callback == "" {
		resp.WriteHeader(http.StatusOK, "application/json; charset=utf-8")
		if _, err = resp.Out.Write([]byte(r.obj)); err != nil {
			log.Errorln("Response write failed:", err)
		}
		return
	}

	resp.WriteHeader(http.StatusOK, "application/javascript; charset=utf-8")
	if _, err = resp.Out.Write([]byte(r.callback + "(")); err != nil {
		log.Errorln("Response write failed:", err)
	}
	if _, err = resp.Out.Write([]byte(r.obj)); err != nil {
		log.Errorln("Response write failed:", err)
	}
	if _, err = resp.Out.Write([]byte(");")); err != nil {
		log.Errorln("Response write failed:", err)
	}
}
