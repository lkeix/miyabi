package miyabi

import (
	"net/http"
	"strings"
)

type (
	// Router control routing
	Router struct {
		tree map[string]*Route
	}
	// Route is node of routing
	Route struct {
		handler  HandlerFunc
		children map[string]*Route
	}
)

const (
	separator = "/"
	coron     = ":"
)

// NewRouter create router instance.
func NewRouter() *Router {
	return &Router{
		tree: map[string]*Route{
			http.MethodGet: {
				handler:  nil,
				children: make(map[string]*Route),
			},
			http.MethodPost: {
				handler:  nil,
				children: make(map[string]*Route),
			},
		},
	}
}

func (r *Router) search(method, path string) *HandlerFunc {
	currentNode := r.tree[method]
	if path == "/" {
		return &currentNode.handler
	}
	comparePath := ""
	for _, separatedStr := range strings.Split(path, "/") {
		for charIdx := range separatedStr {
			comparePath += string(separatedStr[charIdx])
			if nextNode, exist := currentNode.children[string(separatedStr[charIdx])]; exist {
				if comparePath == path {
					return &nextNode.handler
				}
				currentNode = nextNode
				continue
			}
		}
		comparePath += "/"
	}
	return nil
}

func (r *Router) insert(method, path string, handler HandlerFunc) {
	currentNode := r.tree[method]
	if path == "/" {
		*currentNode = newRoute(handler)
		return
	}
	comparePath := ""
	for _, separatedStr := range strings.Split(path, "/") {
		for charIdx := range separatedStr {
			comparePath += string(separatedStr[charIdx])
			if nextNode, exist := currentNode.children[string(separatedStr[charIdx])]; exist {
				currentNode = nextNode
				continue
			}
			// children nil
			if comparePath == path {
				route := newRoute(handler)
				currentNode.children[string(separatedStr[charIdx])] = &route
				return
			}
			route := newRoute(handler)
			currentNode.children[string(separatedStr[charIdx])] = &route
			currentNode = currentNode.children[string(separatedStr[charIdx])]
		}
		comparePath += "/"
	}
}

func newRoute(handler HandlerFunc) Route {
	return Route{
		handler:  handler,
		children: make(map[string]*Route),
	}
}

// GET set http handler on method GET
func (r *Router) GET(path string, handler HandlerFunc) {
	r.insert(http.MethodGet, path, handler)
}
