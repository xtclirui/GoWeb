package main

import (
	"net/http"
	"web"
)

func main() {
	w := web.New()
	w.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	w.GET("/hello", func(c *web.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	w.GET("/hello/:name", func(c *web.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	w.GET("/assets/*filepath", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{"filepath": c.Param("filepath")})
	})

	_ = w.Run(":9999")
}
