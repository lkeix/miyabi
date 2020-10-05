package miyabi

import (
	"html/template"
	"net/http"
)

type (
	// Context is request, writer context
	Context struct {
		Response  *Response
		Request   *RequestContent
		Handler   HandlerFunc
		aborted   bool
		HTMLtmpls map[string][]string
		IsTSL     bool
	}

	// Templates tmpl files ma
	Templates map[string][]string
)

// NewContext create Context instance.
func NewContext(w *http.ResponseWriter, r *http.Request) Context {
	var ctx Context
	ctx.Response = NewResponse(w)
	ctx.Request = NewRequest(r)
	ctx.HTMLtmpls = make(map[string][]string)
	ctx.aborted = false
	return ctx
}

// Execute parseemplate file.
func (ctx *Context) Execute(label string, data interface{}) {
	t, err := template.ParseFiles(ctx.HTMLtmpls[label]...)
	if err != nil {
		panic(err)
	}
	t.Execute(*ctx.Response.Writer, data)
}

// AddTemplates mplates files
func (ctx *Context) AddTemplates(label string, files ...string) {
	ctx.HTMLtmpls[label] = files
}

// Abort abort process of request.
func (ctx *Context) Abort() {
	ctx.aborted = true
}
