package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 自定义个日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("my custom middle")
	}
}

func main() {
	r := gin.New()
	// 通过use设置全局中间件
	// 设置日志中间件，主要用于打印请求日志
	r.Use(Logger())

	r.GET("/testMyMid", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})
	r.Run(":8080")
}
