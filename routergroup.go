package main

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	Handlers []HandlerFunc
	basePath string
	engine   *Engine
}

func (g *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	newGroup := &RouterGroup{
		Handlers: g.combineHandlers(handlers),
		basePath: g.calculateAbsolutePath(relativePath),
		engine:   g.engine,
	}
	g.engine.groups = append(g.engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	g.Handlers = append(g.Handlers, middleware...)
}

// POST is a shortcut for router.Handle("POST", path, handlers).
func (g *RouterGroup) POST(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodPost, relativePath, handler)
}

// GET is a shortcut for router.Handle("GET", path, handlers).
func (g *RouterGroup) GET(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodGet, relativePath, handler)
}

func (g *RouterGroup) handle(httpMethod, relativePath string, handler HandlerFunc) {
	absolutePath := g.calculateAbsolutePath(relativePath)
	g.engine.router.addRoute(httpMethod, absolutePath, handler)
}

func (g *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(g.basePath, relativePath)
}

func (g *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(g.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], handlers)
	return mergedHandlers
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	return finalPath
}
