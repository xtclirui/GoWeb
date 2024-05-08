package web

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
// 以前缀来区分分组,并且支持分组的嵌套
// 作用在/post分组上的中间件，也都会作用在子分组，子分组还可以应用自己特有的中间件。
// 分组是为了更好的注册路由
type (
	RouterGroup struct {
		prefix      string        // 前缀
		middlewares []HandlerFunc // 中间件
		parent      *RouterGroup  // 嵌套,分组父亲
		web         *Web          // 所以组共享一个web实例
	}

	Web struct {
		*RouterGroup // 字段没有变量名，默认使用类型作为字段名，模拟继承关系
		router       *router
		groups       []*RouterGroup // store all groups
	}
)

// New is the constructor of gee.Engine
func New() *Web {
	web := &Web{router: newRouter()}
	web.RouterGroup = &RouterGroup{web: web}
	web.groups = []*RouterGroup{web.RouterGroup}
	return web
}

// Group is defined to create a new RouterGroup
// 创建该分组的子分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// 所以组共享一个web实例
	web := group.web
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		web:    web,
	}
	web.groups = append(web.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, url string, handler HandlerFunc) {
	url = group.prefix + url
	log.Printf("Route %4s - %s", method, url)
	group.web.router.addRouter(method, url, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(url string, handler HandlerFunc) {
	group.addRoute("GET", url, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(url string, handler HandlerFunc) {
	group.addRoute("POST", url, handler)
}

func (group *RouterGroup) AddMid(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Run defines the method to start a http server
func (web *Web) Run(addr string) (err error) {
	return http.ListenAndServe(addr, web)
}

func (web *Web) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range web.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.mid = middlewares
	web.router.handle(c)
}
