package web

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Web struct {
	router *router
}

func New() *Web {
	return &Web{router: newRouter()}
}

func (web *Web) addRoute(method string, url string, handler HandlerFunc) {
	web.router.addRouter(method, url, handler)
}

func (web *Web) GET(url string, handler HandlerFunc) {
	web.addRoute("GET", url, handler)
}

func (web *Web) POST(url string, handler HandlerFunc) {
	web.addRoute("POST", url, handler)
}

func (web *Web) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	web.router.handle(c)
}

func (web *Web) Run(addr string) (err error) {
	return http.ListenAndServe(addr, web)
}
