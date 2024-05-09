package web

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
		groups       []*RouterGroup // store all groups，至少含有一个，web本身
		htmlTemplate *template.Template
		funcMap      template.FuncMap
	}
)

// New is the constructor of gee.Engine
func New() *Web {
	web := &Web{router: newRouter()}
	web.RouterGroup = &RouterGroup{web: web}
	// 初始化时，添加web实例的web.RouterGroup
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
	// 首先根据前缀初始化中间件
	var middlewares []HandlerFunc
	for _, group := range web.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 一个http连接对应一个Context
	c := newContext(w, req)
	c.mid = middlewares
	c.web = web
	web.router.handle(c)
}

// 创建静态处理函数
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 服务器静态文件
func (group *RouterGroup) Static(relativePath string, root string) {
	// http.Dir(root) 返回实现了http.FileSystem接口的Open 方法
	// http.Dir用于指定一个文件系统目录，以便在HTTP服务器中提供静态文件服务
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	url := path.Join(relativePath, "/*filepath")
	group.GET(url, handler)
}

func (web *Web) SetFuncMap(funcMap template.FuncMap) {
	web.funcMap = funcMap
}

func (web *Web) LoadHTMLGlob(url string) {
	web.htmlTemplate = template.Must(template.New("").Funcs(web.funcMap).ParseGlob(url))
}
