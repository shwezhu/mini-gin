package main

import (
	"net/http"
)

func hello(c *Context) {
	c.String(http.StatusOK, "PostForm:%s, Query:%s\n", c.PostForm("name"), c.Query("name"))
}

func main() {
	mini := New()
	mini.GET("/", hello)

	err := http.ListenAndServe(":8080", mini)
	if err != nil {
		panic(err)
	}
}
