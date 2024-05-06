package web

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, url string, handler HandlerFunc) {
	key := method + "-" + url
	engine.router[key] = handler
}

func (engine *Engine) GET(url string, handler HandlerFunc) {
	engine.addRoute("GET", url, handler)
}

func (engine *Engine) POST(url string, handler HandlerFunc) {
	engine.addRoute("POST", url, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if hanler, ok := engine.router[key]; ok {
		hanler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND", req.URL)
	}
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
