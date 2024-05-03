package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
	// dynamic url parameters
	Params map[string]string

	// middleware
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:   w,
		Req:      r,
		Path:     r.URL.Path,
		Method:   r.Method,
		handlers: nil,
		index:    -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	// 1.application/ x-www-form-urlencoded form body
	// 2.query parameters in the URL (this is the same as c.Query())
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// SendString Returns an error if the write fails, the reason could be the following:
// 1. The connection is closed by the client
// 2. If content length was specified, and you attempt to write more than that: http.ErrContentLength
// 3...
// Suggestions: You can log the error, or ignore it.
// Learn more: https://stackoverflow.com/a/43976633/16317008
func (c *Context) SendString(code int, format string, values ...interface{}) error {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	// If err!=nil, http.Error(c.Writer, err.Error(), 500) won't work
	return err
}

func (c *Context) SendJSON(code int, obj interface{}) error {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	return encoder.Encode(obj)
}

// SendFile parameter filepath has to be an absolute path.
func (c *Context) SendFile(filePath string) {
	filename := filepath.Base(filePath)
	c.SetHeader("Content-Type", "application/octet-stream")
	// strconv.Quote ensures the filename can be handled safely because of special characters.
	c.SetHeader("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	http.ServeFile(c.Writer, c.Req, filePath)
}
