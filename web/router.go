package web

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRouter(method string, url string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, url)
	key := method + "-" + url
	r.handlers[key] = handler
}
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNoContent, "404 NOT FOUND: %s\n", c.Path)
	}
}
