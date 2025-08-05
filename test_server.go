// 绠€鍗曟祴璇曟湇鍔″櫒
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("鍚姩绠€鍗曟祴璇曟湇鍔″櫒...")

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "娴嬭瘯鏈嶅姟鍣ㄨ繍琛屾甯?,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "鏈哄櫒浜鸿矾寰勭紪杈戝櫒娴嬭瘯鐗?,
		})
	})

	fmt.Println("鏈嶅姟鍣ㄥ惎鍔ㄥ湪 :8080")
	r.Run(":8080")
}
