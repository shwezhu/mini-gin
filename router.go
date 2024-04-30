package main

import (
	"net/http"
	"strings"
)

// router
// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	r.handlers[method+"-"+pattern] = handler

	// root 的子节点为 GET, POST 等
	// 真正的路径为 POST /a/b/c 或 GET /a/b/c, 注意到区别了吗, 在根节点怎么分叉的
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	parts := parsePattern(pattern)
	r.roots[method].insertChild(pattern, parts, 0)
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	params := make(map[string]string)
	itemsOfPath := parsePattern(path)
	nod := root.matchPattern(itemsOfPath, 0)
	if nod != nil {
		itemsOfPattern := parsePattern(nod.pattern)
		// 理解这段代码 需要知道动态 url 是怎么用的, 可参考 gin
		for index, item := range itemsOfPattern {
			if item[0] == ':' {
				params[item[1:]] = itemsOfPath[index]
				continue
			}
			if item[0] == '*' && len(item) > 1 {
				params[item[1:]] = strings.Join(itemsOfPath[index:], "/")
				break
			}
		}
		return nod, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	nod, params := r.getRoute(c.Method, c.Path)
	if nod != nil {
		c.Params = params
		key := c.Method + "-" + nod.pattern
		r.handlers[key](c)
	} else {
		http.Error(c.Writer, "404 Not Found", http.StatusNotFound)
	}
}

// Only one * is allowed
// /a/*/c/* => /a
func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}
