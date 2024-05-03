package main

import (
	"reflect"
	"testing"
)

func TestRouterGroup(t *testing.T) {
	g := &RouterGroup{}
	middleware1 := func(c *Context) {}
	middleware2 := func(c *Context) {}

	g.Use(middleware1, middleware2)

	expectedHandlers := []HandlerFunc{middleware1, middleware2}
	if !reflect.DeepEqual(len(g.Handlers), len(expectedHandlers)) {
		t.Fatalf("expected %d handlers, got %d", len(expectedHandlers), len(g.Handlers))
	}
}
