package main

import (
	"net/http"
	"web"
)

func main() {
	r := web.New()
	r.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Web</h1>")
	})
	r.GET("/hello", func(c *web.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":9999")
}
