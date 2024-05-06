package web

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, url string, handler HandlerFunc) {
	engine.router.addRouter(method, url, handler)
}

func (engine *Engine) GET(url string, handler HandlerFunc) {
	engine.addRoute("GET", url, handler)
}

func (engine *Engine) POST(url string, handler HandlerFunc) {
	engine.addRoute("POST", url, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
