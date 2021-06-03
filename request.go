package miyabi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type (
	// RequestContent is lapping http.Request
	RequestContent struct {
		Base        *http.Request
		Data        interface{}
		PathParams  []Param
		QueryParams map[string][]string
	}
)

// NewRequest create RequestContent instance
func NewRequest(r *http.Request) *RequestContent {
	return &RequestContent{
		Base:        r,
		PathParams:  []Param{},
		QueryParams: make(map[string][]string),
	}
}

// Parse parse request data.
func (req *RequestContent) Parse() {
	contentType := req.Base.Header.Get("Content-Type")
	if contentType == "text/plain" {
		req.textParser()
		return
	}
	if contentType == "applicati/json" {
		req.jsonParser()
		return
	}
	if strings.HasPrefix(contentType, "image/jpeg") ||
		strings.HasPrefix(contentType, "image/png") ||
		strings.HasPrefix(contentType, "imageif") ||
		strings.HasPrefix(contentType, "image/bmp") ||
		strings.HasPrefix(contentType, "image/svg+xml") {
		req.imageParser()
		return
	}
	if strings.HasPrefix(contentType, "multipart/form-data") {
		req.fileParser()
	}
	req.quaryParser()
}

func (req *RequestContent) textParser() {
	buf := make([]byte, req.Base.ContentLength)
	req.Base.Body.Read(buf)
	req.Data = string(buf)
}

func (req *RequestContent) jsonParser() {
	buf := make([]byte, req.Base.ContentLength)
	req.Base.Body.Read(buf)
	json.Unmarshal(buf, &req.Data)
}

func (req *RequestContent) imageParser() {
	buf := make([]byte, req.Base.ContentLength)
	req.Base.Body.Read(buf)
	// image convert to base64 string
	req.Data = base64.StdEncoding.EncodeToString(buf)
}

func (req *RequestContent) fileParser() {
	buf := make([]byte, req.Base.ContentLength)
	req.Base.Body.Read(buf)
	req.Data = buf
}

func (req *RequestContent) quaryParser() {
	r := req.Base
	r.ParseForm()
	for key, value := range r.Form {
		req.QueryParams[key] = value
	}
}
