// 简单测试服务器
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("启动简单测试服务器...")

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "测试服务器运行正常",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "机器人路径编辑器测试版",
		})
	})

	fmt.Println("服务器启动在 :8080")
	r.Run(":8080")
}
