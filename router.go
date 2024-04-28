package main

import "net/http"

type Router struct {
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	r.handlers[method+"-"+pattern] = handler
}

func (r *Router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
		return
	}
	http.Error(c.Writer, "404 Not Found", http.StatusNotFound)
}
