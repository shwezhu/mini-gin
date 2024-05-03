package main

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	*RouterGroup
	router *router
	// use a slice to store all groups
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
		RouterGroup: &RouterGroup{
			basePath: "/",
		},
	}
	engine.RouterGroup.engine = engine
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute(http.MethodGet, pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute(http.MethodPost, pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.basePath) {
			middlewares = append(middlewares, group.Handlers...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	e.router.handle(c)
}
