package miyabi

import (
	"fmt"
)

type (
	node struct {
		prefix         string
		parent         *node
		staticChildren []*node
		paramsKey      []string
		originPath     string
		kind           kind
		middlewares    []*HandlerFunc
		paramsChild    *node
		anyChild       *node
		handler        *HandlerFunc
	}

	// PathParam record path parameter.
	PathParam struct {
		key   string
		value string
	}

	// Route is node of routing
	Router struct {
		handler  HandlerFunc
		tree     *node
		param    PathParam
		maxParam int
	}

	kind uint8

	// Info is route info, used setup log.
	Info struct {
		path   string
		method string
	}

	Param struct {
		Key string
		Val string
	}
)

const (
	separator = '/'
	coron     = ':'
	asta      = '*'

	coronLabel = byte(coron)
	astaLabel  = byte(asta)

	isStatic = iota
	isParam
	isAny
)

// NewRouter create router instance.
func NewRouter() *Router {
	return &Router{
		tree: newTree(),
	}
}

func newTree() *node {
	return &node{}
}

func newNode(t kind, prefix string, parent *node, staticChildren []*node, handler HandlerFunc, originPath string, paramsKey []string, paramChildren, anyChildren *node) *node {
	return &node{
		kind:           t,
		prefix:         prefix,
		parent:         parent,
		staticChildren: staticChildren,
		originPath:     originPath,
		paramsKey:      paramsKey,
		handler:        &handler,
		paramsChild:    paramChildren,
		anyChild:       anyChildren,
	}
}

func (o *node) reset(t kind, prefix string, staticChildren []*node, handler HandlerFunc, originPath string, paramsKey []string, paramsChild, anyChild *node) {
	o.kind = t
	o.prefix = prefix
	o.staticChildren = staticChildren
	o.handler = &handler
	o.originPath = originPath
	o.paramsKey = paramsKey
	o.paramsChild = paramsChild
	o.anyChild = anyChild
}

func (o *node) findChild(label byte) *node {
	for _, child := range o.staticChildren {
		if child.prefix[0] == label {
			return child
		}
	}
	if label == coronLabel {
		return o.paramsChild
	}
	if label == astaLabel {
		return o.anyChild
	}
	return nil
}

func (r *Router) Insert(method, path string, handler HandlerFunc) error {
	path = pathValidate(path)

	pnames := []string{}
	originPath := path

	if err := handlerValidate(handler, originPath); err != nil {
		return err
	}

	for i, lcpIndex := 0, len(path); i < lcpIndex; i++ {
		if path[i] == coron {
			if i > 0 && path[i-1] == '\\' {
				path = path[:i-1] + path[i:]
				i--
				lcpIndex--
				continue
			}
			j := i + 1

			r.insert(method, path[:i], nil, isStatic, "", nil)
			for ; i < lcpIndex && path[i] != separator; i++ {
			}

			pnames = append(pnames, path[j:i])

			path = path[:j] + path[i:]

			i, lcpIndex = j, len(path)

			if i == lcpIndex {
				fmt.Println("insert A")
				fmt.Println(i)
				fmt.Println(path)
				fmt.Println(pnames)
			} else {
				fmt.Println("insert B")
				fmt.Println(i)
				fmt.Println(path)
				fmt.Println(pnames)
			}
		} else if path[i] == asta {
			fmt.Println("insert C")
			pnames = append(pnames, "*")
			fmt.Println(i)
			fmt.Println(path)
			fmt.Println(pnames)
		}
	}

	fmt.Println("insert D")
	fmt.Println(path)
	fmt.Println(pnames)
	return nil
}

func (r *Router) insert(method, path string, handler HandlerFunc, k kind, originPath string, paramsKey []string) {
	currentNode := r.tree

	search := path

	for {
		searchLen := len(search)
		prefixLen := len(currentNode.prefix)
		lcpLen := 0

		max := prefixLen
		max = maxi(max, searchLen)

		for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		}

		if lcpLen == 0 {
			currentNode.prefix = search
			if handlerValidate(handler, originPath) != nil {
				currentNode.handler = &handler
				currentNode.paramsKey = paramsKey
			}
		} else if lcpLen < prefixLen {
			node := newNode(
				currentNode.kind,
				currentNode.prefix[lcpLen:],
				currentNode,
				currentNode.staticChildren,
				*currentNode.handler,
				currentNode.originPath,
				currentNode.paramsKey,
				currentNode.paramsChild,
				currentNode.anyChild,
			)

			for _, child := range currentNode.staticChildren {
				child.parent = node
			}

			if currentNode.paramsChild != nil {
				currentNode.paramsChild.parent = node
			}
			if currentNode.anyChild != nil {
				currentNode.anyChild.parent = node
			}

			currentNode.reset(
				isStatic,
				currentNode.prefix[:lcpLen],
				nil,
				handler,
				"",
				nil,
				nil,
				nil,
			)

			currentNode.staticChildren = append(currentNode.staticChildren, node)

			if lcpLen == searchLen {
				currentNode.kind = k
				currentNode.handler = &handler
				currentNode.originPath = originPath
				currentNode.paramsKey = paramsKey
			} else {
				node = newNode(k, search[lcpLen:], currentNode, nil, handler, originPath, paramsKey, nil, nil)
				node.handler = &handler
				currentNode.staticChildren = append(currentNode.staticChildren, node)
			}
		} else if lcpLen < searchLen {
			search = search[lcpLen:]

			childNode := currentNode.findChild(search[0])
			if childNode != nil {
				currentNode = childNode
				continue
			}

			node := newNode(k, search, currentNode, nil, handler, originPath, paramsKey, nil, nil)
			switch k {
			case isStatic:
				currentNode.staticChildren = append(currentNode.staticChildren, node)
			case isParam:
				currentNode.paramsChild = node
			case isAny:
				currentNode.anyChild = node
			}
		} else {
			if handler != nil {
				currentNode.handler = &handler
				currentNode.originPath = originPath
				if len(currentNode.originPath) == 0 {
					currentNode.paramsKey = paramsKey
				}
			}
		}
		break
	}
}

/*
func (tr *Router) insert(method, path string, handler HandlerFunc) {
	currentNode := tr.body[method]
	if path == separator {
		*currentNode = newRoute(handler, "")
		return
	}
	comparePath := ""
	for _, separatedStr := range strings.Split(path, separator) {
		for charIdx := 0; charIdx < len(separatedStr); charIdx++ {
			nextNode, exist := currentNode.children[string(separatedStr[charIdx])]
			comparePath = strings.Join([]string{comparePath, string(separatedStr[charIdx])}, "")
			// pathparameter
			if string(separatedStr[0]) == coron {
				charIdx = len(separatedStr)
				route := newRoute(nil, separatedStr[1:charIdx])
				currentNode.children[coron] = &route
				currentNode = currentNode.children[coron]
				comparePath = strings.Join([]string{comparePath, separatedStr[1:charIdx]}, "")
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
			if string(comparePath) == path && string(separatedStr[0]) != coron {
				route := newRoute(handler, "")
				currentNode.children[string(separatedStr[charIdx])] = &route
				return
			}
			// children nil
			route := newRoute(nil, "")
			currentNode.children[string(separatedStr[charIdx])] = &route
			currentNode = currentNode.children[string(separatedStr[charIdx])]
		}
		comparePath = strings.Join([]string{comparePath, separator}, "")
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

func (tr *Router) search(method, path string) (*HandlerFunc, []Param) {
	currentNode := tr.body[method]
	if path == separator {
		return &currentNode.handler, nil
	}
	comparePath := ""
	params := []Param{}
	for _, separatedStr := range strings.Split(path, separator) {
		for charIdx := 0; charIdx < len(separatedStr); charIdx++ {
			nextNode, exist := currentNode.children[string(separatedStr[charIdx])]
			comparePath = strings.Join([]string{comparePath, string(separatedStr[charIdx])}, "")
			if exist {
				if string(comparePath) == path && nextNode.handler != nil {
					return &nextNode.handler, params
				}
				currentNode = nextNode
				continue
			}
			// children have path parameter delimiter coron
			if nextNode, exist := currentNode.children[coron]; exist {
				comparePath = strings.Join([]string{comparePath[0:max(0, len(comparePath)-1)], string(separatedStr)}, "")
				var param Param
				param.Key = nextNode.param.key
				param.Val = separatedStr
				params = append(params, param)
				charIdx = len(separatedStr) - 1
				if string(comparePath) == path && nextNode.handler != nil {
					return &nextNode.handler, params
				}
				currentNode = nextNode
			}
		}
		comparePath = strings.Join([]string{comparePath, separator}, "")
	}
	return nil, nil
}
*/

func joinByte(base []byte, target byte) []byte {
	return append(base, target)
}

func joinBytes(base []byte, targets []byte) []byte {
	return append(base, targets...)
}

// RunMiddleware run middleware resisterd on router
func (o *node) RunMiddleware(ctx *Context) {
	for _, middleware := range o.middlewares {
		handler := *middleware
		handler(ctx)
	}
}

/*
// RunMiddleware run middleware resisterd on group
func (g *Group) RunMiddleware(ctx *Context) {
	for _, middleware := range g.middlewares {
		handler := middleware
		handler(ctx)
	}
}

// GET set http handler on method GET
func (r *Router) GET(path string, handler HandlerFunc) {
	r.Tree.insert(http.MethodGet, path, handler)
	r.RouterInfo = append(r.RouterInfo, &Info{path: path, method: http.MethodGet})
}

// POST set http handler on method POST
func (r *Router) POST(path string, handler HandlerFunc) {
	r.Tree.insert(http.MethodPost, path, handler)
	r.RouterInfo = append(r.RouterInfo, &Info{path: path, method: http.MethodPost})
}

// GET set http handler on method GET in Group
func (g *Group) GET(path string, handler HandlerFunc) {
	g.Tree.insert(http.MethodGet, g.originPath+path, handler)
	g.GroupInfo = append(g.GroupInfo, &Info{path: path, method: http.MethodGet})
}

// POST set http handler on method POST in Group
func (g *Group) POST(path string, handler HandlerFunc) {
	g.Tree.insert(http.MethodPost, g.originPath+path, handler)
	g.GroupInfo = append(g.GroupInfo, &Info{path: path, method: http.MethodPost})
}
*/

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
