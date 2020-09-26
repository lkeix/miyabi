package miyabi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type (
	// Context is request, writer context
	Context struct {
		Response *Response
		Request  RequestContent
		Handler  HandlerFunc
	}
	// RequestContent is lapping http.Request
	RequestContent struct {
		Base        *http.Request
		Data        interface{}
		PathParams  map[string]string
		QueryParams map[string][]string
	}
)

// NewContext create Context instance.
func NewContext(w *http.ResponseWriter, r *http.Request) Context {
	var ctx Context
	ctx.Response = NewResponse(w)
	ctx.Request.Base = r
	ctx.Request.QueryParams = make(map[string][]string)
	return ctx
}

// Parse parse ctx request data.
func (ctx *Context) Parse() {
	contentType := ctx.Request.Base.Header.Get("Content-Type")
	if contentType == "text/plain" {
		ctx.textParser()
		return
	}
	if contentType == "application/json" {
		ctx.jsonParser()
		return
	}
	if strings.HasPrefix(contentType, "image/jpeg") ||
		strings.HasPrefix(contentType, "image/png") ||
		strings.HasPrefix(contentType, "image/gif") ||
		strings.HasPrefix(contentType, "image/bmp") ||
		strings.HasPrefix(contentType, "image/svg+xml") {
		ctx.imageParser()
		return
	}
	if strings.HasPrefix(contentType, "multipart/form-data") {
		ctx.fileParser()
	}
	ctx.quaryParser()
}

func (ctx *Context) textParser() {
	buf := make([]byte, ctx.Request.Base.ContentLength)
	ctx.Request.Base.Body.Read(buf)
	ctx.Request.Data = string(buf)
}

func (ctx *Context) jsonParser() {
	buf := make([]byte, ctx.Request.Base.ContentLength)
	ctx.Request.Base.Body.Read(buf)
	json.Unmarshal(buf, &ctx.Request.Data)
}

func (ctx *Context) imageParser() {
	buf := make([]byte, ctx.Request.Base.ContentLength)
	ctx.Request.Base.Body.Read(buf)
	// image convert to base64 string
	ctx.Request.Data = base64.StdEncoding.EncodeToString(buf)
}

func (ctx *Context) fileParser() {
	buf := make([]byte, ctx.Request.Base.ContentLength)
	ctx.Request.Base.Body.Read(buf)
	ctx.Request.Data = buf
}

func (ctx *Context) quaryParser() {
	req := ctx.Request.Base
	req.ParseForm()
	for key, value := range req.Form {
		ctx.Request.QueryParams[key] = value
	}
}
