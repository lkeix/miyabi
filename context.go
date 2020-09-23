package miyabi

import "net/http"

type (
	// Context is request, writer context
	Context struct {
		Response *Response
		Request  *http.Request
		Handler  HandlerFunc
	}
)

// NewContext create Context instance.
func NewContext(w *http.ResponseWriter, r *http.Request) Context {
	var ctx Context
	ctx.Response = NewResponse(w)
	ctx.Request = r
	return ctx
}
