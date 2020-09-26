package miyabi

import (
	"html/template"
	"net/http"
	"sync"
)

type (
	// HandlerFunc is this framework handler
	HandlerFunc func(*Context)

	// Miyabi is this web framework base class.
	Miyabi struct {
		FuncMap template.FuncMap
		Routing *Router
		pool    sync.Pool
		Debug   bool
		server  http.Server
	}
)

// New create Miyabi instance, return it. if you want to use debug mode, you write true on arg.
func New(debug ...bool) *Miyabi {
	var doDebug bool
	if len(debug) == 1 {
		doDebug = debug[0]
	}
	myb := &Miyabi{
		FuncMap: template.FuncMap{},
		Debug:   doDebug,
	}
	myb.pool.New = func() interface{} {
		return NewContext(nil, nil)
	}
	return myb
}

func (myb *Miyabi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := myb.pool.Get().(Context)
	ctx = NewContext(&w, r)
	method := ctx.Request.Base.Method
	url := ctx.Request.Base.URL.Path
	handler, params := myb.Routing.search(method, url)
	if handler != nil {
		ctx.Handler = *handler
		ctx.Request.PathParams = params
		ctx.Parse()
		ctx.Handler(&ctx)
		myb.pool.Put(ctx)
	} else {
		ctx.Handler = noRoute()
		ctx.Handler(&ctx)
	}
}

func noRoute() HandlerFunc {
	return func(ctx *Context) {
		ctx.Response.WriteResponse("404 Not Found.")
	}
}

// Serve start http server
func (myb *Miyabi) Serve(port string) {
	var s http.Server
	s.Handler = myb
	s.Addr = port
	s.ListenAndServe()
}
