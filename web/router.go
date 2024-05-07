package web

import (
	"strings"
)

// 注册了四个路由：
// "/"
// "/hello"
// "/hello/:name"
// "/assets/*filepath"
type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node // GET POST 等method分开存放在不同的树
}

func parseUrl(url string) []string {
	vs := strings.Split(url, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc), roots: make(map[string]*node)}
}

func (r *router) addRouter(method string, url string, handler HandlerFunc) {
	parts := parseUrl(url)

	key := method + "-" + url
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(url, parts, 0)

	r.handlers[key] = handler
}

func (r *router) getRoute(method string, url string) (*node, map[string]string) {
	searchParts := parseUrl(url)
	// params存放的是注册路由时，含有模糊匹配时，匹配的对应关系
	// "/p/:lang/doc" "/p/go/doc"
	// key : "lang", value : "go"
	// "/static/*filepath" "/static/js/jQuery.js"
	// key : "filepath", value : "js/jQuery.js"
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// 寻找到与url匹配的 *node
	n := root.search(searchParts, 0)

	if n != nil {
		// parts却注册的url去掉'/'后的切片
		parts := parseUrl(n.url)
		for index, part := range parts {
			if part[0] == ':' {

				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 对整个web没有影响，只是将匹配结果放入Context中，等待后续处理
		c.Params = params
		key := c.Method + "-" + n.url
		r.handlers[key](c)
	} else {
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
}
