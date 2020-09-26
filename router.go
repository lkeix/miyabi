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

	// PathParam record path parameter.
	PathParam struct {
		key   string
		value string
	}

	// Route is node of routing
	Route struct {
		handler  HandlerFunc
		children map[string]*Route
		param    PathParam
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

func (r *Router) search(method, path string) (*HandlerFunc, map[string]string) {
	currentNode := r.tree[method]
	if path == separator {
		return &currentNode.handler, nil
	}
	comparePath := ""
	params := make(map[string]string)
	for _, separatedStr := range strings.Split(path, separator) {
		for charIdx := 0; charIdx < len(separatedStr); charIdx++ {
			nextNode, exist := currentNode.children[string(separatedStr[charIdx])]
			comparePath += string(separatedStr[charIdx])
			if exist {
				if comparePath == path && nextNode.handler != nil {
					return &nextNode.handler, params
				}
				currentNode = nextNode
				continue
			}
			// children have path parameter delimiter coron
			if nextNode, exist := currentNode.children[coron]; exist {
				comparePath = comparePath[0:len(comparePath)-1] + separatedStr
				params[nextNode.param.key] = separatedStr
				charIdx = len(separatedStr) - 1
				if comparePath == path && nextNode.handler != nil {
					return &nextNode.handler, params
				}
				currentNode = nextNode
			}
		}
		comparePath += separator
	}
	return nil, nil
}

func (r *Router) insert(method, path string, handler HandlerFunc) {
	currentNode := r.tree[method]
	if path == separator {
		*currentNode = newRoute(handler, "")
		return
	}
	comparePath := ""
	for _, separatedStr := range strings.Split(path, separator) {
		for charIdx := 0; charIdx < len(separatedStr); charIdx++ {
			nextNode, exist := currentNode.children[string(separatedStr[charIdx])]
			comparePath += string(separatedStr[charIdx])
			// pathparameter
			if string(separatedStr[0]) == coron {
				charIdx = len(separatedStr) - 1
				route := newRoute(nil, separatedStr[1:charIdx+1])
				currentNode.children[coron] = &route
				currentNode = currentNode.children[coron]
				comparePath += separatedStr[1 : charIdx+1]
				if comparePath == path {
					currentNode.handler = handler
					return
				}
				continue
			}
			if exist {
				currentNode = nextNode
				continue
			}
			// target path
			if comparePath == path && string(separatedStr[0]) != coron {
				route := newRoute(handler, "")
				currentNode.children[string(separatedStr[charIdx])] = &route
				return
			}
			// children nil
			route := newRoute(nil, "")
			currentNode.children[string(separatedStr[charIdx])] = &route
			currentNode = currentNode.children[string(separatedStr[charIdx])]
		}
		comparePath += separator
	}
}

func newRoute(handler HandlerFunc, key string) Route {
	if key == "" {
		return Route{
			handler:  handler,
			children: make(map[string]*Route),
		}
	}
	route := Route{
		handler:  handler,
		children: make(map[string]*Route),
	}
	route.param.key = key
	return route
}

// GET set http handler on method GET
func (r *Router) GET(path string, handler HandlerFunc) {
	r.insert(http.MethodGet, path, handler)
}

// POST set http handler on method POST
func (r *Router) POST(path string, handler HandlerFunc) {
	r.insert(http.MethodPost, path, handler)
}
