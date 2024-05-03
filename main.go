package main

import (
	"log"
)

func onlyForV2() HandlerFunc {
	return func(c *Context) {
		log.Println("onlyForV2...")
	}
}

func onlyForV22() HandlerFunc {
	return func(c *Context) {
		log.Println("onlyForV22...")
	}
}

func main() {
	mini := New()
	v2 := mini.Group("/v2")
	v2.Use(onlyForV2(), onlyForV22()) // v2 group middleware

	{
		v2.GET("/hello/:name", func(c *Context) {
			log.Println("/hello/", c.Param("name"))
		})
		v2.GET("/chat", func(c *Context) {
			log.Println("/chat")
		})
	}

	mini.Run(":8080")
}
