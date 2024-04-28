package main

import (
	"net/http"
)

func hello(c *Context) {
	// _ = c.SendString(http.StatusOK, "PostForm:%s, Query:%s\n", c.PostForm("name"), c.Query("name"))
	c.SendFile("mini-gin.go")
}

func main() {
	mini := New()
	mini.GET("/hello", hello)

	err := http.ListenAndServe(":8080", mini)
	if err != nil {
		panic(err)
	}
}
