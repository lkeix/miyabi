package miyabi

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type (
	// HandlerFunc is this framework handler
	HandlerFunc func(*Context)

	// Miyabi is this web framework base class.
	Miyabi struct {
		Router    *Router
		pool      sync.Pool
		server    http.Server
		isTLS     bool
		Timestamp map[string]time.Time
	}
)

const logBuffer = 1.0

// New create Miyabi instance, return it.
func New() *Miyabi {
	myb := &Miyabi{}
	myb.pool.New = func() interface{} {
		var w http.ResponseWriter
		var r *http.Request
		return NewContext(&w, r)
	}
	myb.Timestamp = make(map[string]time.Time)
	return myb
}

func (myb *Miyabi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := myb.pool.Get().(*Context)
	ctx.IsTSL = myb.isTLS
	ctx.Request.Base = r
	ctx.Response.Writer = &w
	method := ctx.Request.Base.Method
	url := ctx.Request.Base.URL.Path
	route := myb.Router
	handler, params := route.Tree.search(method, url)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	outlog := false
	now := time.Now()
	last, exist := myb.Timestamp[ip]
	if !exist || now.Sub(last).Seconds() >= logBuffer {
		myb.Timestamp[ip] = now
		outlog = true
	}
	if handler != nil {
		route.RunMiddleware(ctx)
		execHandler(ctx, handler, params)
		myb.pool.Put(ctx)
		if outlog {
			requestLog(url, method, 200)
		}
		return
	}
	for i := 0; i < len(route.Groups); i++ {
		group := myb.Router.Groups[i]
		handler, params := group.Tree.search(method, url)
		if handler != nil {
			group.RunMiddleware(ctx)
			execHandler(ctx, handler, params)
			if outlog {
				requestLog(url, method, 200)
			}
			myb.pool.Put(ctx)
			return
		}
	}
	if outlog {
		requestLog(url, method, 404)
	}
	ctx.Handler = noRoute()
	ctx.Handler(ctx)
}

func execHandler(ctx *Context, handler *HandlerFunc, params []Param) {
	ctx.Handler = *handler
	ctx.Request.PathParams = params
	ctx.Handler(ctx)
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
	routerLog(myb)
	s.ListenAndServe()
}

// ServeTLS start https server need cert, key file path
func (myb *Miyabi) ServeTLS(port, cert, key string) {
	var s http.Server
	s.Handler = myb
	s.Addr = port
	myb.isTLS = true
	routerLog(myb)
	s.ListenAndServeTLS(cert, key)
}
