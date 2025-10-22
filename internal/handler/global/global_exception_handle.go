package global

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleGlobalException() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    5000,
					"message": "服务器开小车啦~",
				})
				c.Abort()
			}
		}()
	}
}
