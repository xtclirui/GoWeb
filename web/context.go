package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string // params存放的是注册路由时，含有模糊匹配时，匹配的对应关系
	StatusCode int

	mid   []HandlerFunc // 中间件
	index int           // index是记录当前执行到第几个中间件
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) PostForm(kye string) string {
	// FormValue 方法可以用于获取 HTTP 请求的表单参数值
	return c.Req.FormValue(kye)
}

func (c *Context) Query(kye string) string {
	// 获取请求 URL 中名为 key 的查询参数的值
	return c.Req.URL.Query().Get(kye)
}

func (c *Context) Status(status int) {
	c.StatusCode = status
	c.Writer.WriteHeader(status)
}

func (c *Context) SetHeader(key string, value string) {
	// 设置http头信息
	c.Writer.Header().Set(key, value)
}

// 提供了快速构造String/Data/JSON/HTML响应的方法。
func (c *Context) String(status int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(status)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 空接口可以接受任意类型的参数  interface{}能接收map[string]interface{}类型

func (c *Context) JSON(status int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(status)
	// 作用是创建一个 JSON 编码器，用于将数据编码为 JSON 格式并写入指定的 http.ResponseWriter
	encoder := json.NewEncoder(c.Writer)
	// 调用 Encode() 方法后，编码器会将数据对象 obj 转换为 JSON 字符串，并将其写入到编码器指向的 io.Writer 接口对象中。
	// obj 是map类型
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(status int, data []byte) {
	c.Status(status)
	c.Writer.Write(data)
}

func (c *Context) HTML(status int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(status)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 调用中间件
func (c *Context) Next() {
	c.index++
	s := len(c.mid)

	for ; c.index < s; c.index++ {
		c.mid[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.mid)
	c.JSON(code, H{"message": err})
}
