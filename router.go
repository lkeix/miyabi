package miyabi

import (
	"net/http"
	"strings"
)

type (

	tree struct {
		body map[string]*Route
	}

	// Router control routing
	Router struct {
		Tree *tree
		Groups []*Group
		middlewares []*HandlerFunc
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

	// Group is sub-routing data structure
	Group struct {
		basePath string
		middlewares []*HandlerFunc
		Tree *tree
	}
)

const (
	separator = "/"
	coron     = ":"
)

// NewRouter create router instance.
func NewRouter() *Router {
	return &Router{
		Tree: newTree(),
	}
}


func newTree() *tree {
	return &tree{
		map[string]*Route{
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

// NewGroup create group instance
func (r *Router) NewGroup(group string) *Group {
	return &Group{
		basePath: group,
		Tree: newTree(),
	}
}

// Apply middleware handler on group
func (r *Router) Apply(handlers ...HandlerFunc) {
	for _, handler := range handlers {
		r.middlewares = append(r.middlewares, &handler)
	}
}

// Apply middleware handler on group
func (g *Group) Apply(handlers ...HandlerFunc) {
	for _, handler := range handlers {
		g.middlewares = append(g.middlewares, &handler)
	}
}

func (tr *tree) search(method, path string) (*HandlerFunc, map[string]string) {
	currentNode := tr.body[method]
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

func (tr *tree) insert(method, path string, handler HandlerFunc) {
	currentNode := tr.body[method]
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

// AppendGroup append group for router groups
func (r *Router) AppendGroup(group *Group) {
	r.Groups = append(r.Groups, group)
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

// RunMiddleware run middleware resisterd on router
func (r *Router) RunMiddleware(ctx *Context) {
	for _, middleware := range r.middlewares {
		handler := *middleware
		handler(ctx)
	}
}

// RunMiddleware run middleware resisterd on group
func (g *Group) RunMiddleware(ctx *Context) {
	for _, middleware := range g.middlewares {
		handler := *middleware
		handler(ctx)
	}
}

// GET set http handler on method GET
func (r *Router) GET(path string, handler HandlerFunc) {
	r.Tree.insert(http.MethodGet, path, handler)
}

// POST set http handler on method POST
func (r *Router) POST(path string, handler HandlerFunc) {
	r.Tree.insert(http.MethodPost, path, handler)
}

// GET set http handler on method GET in Group
func (g *Group) GET(path string, handler HandlerFunc) {
	g.Tree.insert(http.MethodGet, g.basePath + path, handler)
}

// POST set http handler on method POST in Group
func (g *Group) POST(path string, handler HandlerFunc) {
	g.Tree.insert(http.MethodPost, g.basePath + path, handler)
}
