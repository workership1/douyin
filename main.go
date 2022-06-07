package main

import (
	"github.com/gin-gonic/gin"

)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		// 该函数直接返回JSON字符串
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
