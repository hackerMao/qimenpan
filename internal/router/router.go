package router

import (
	"github.com/gin-gonic/gin"
	"qimenpan/internal/handler/global"
)

func InitRouter() *gin.Engine {
	gin.SetMode("debug")
	router := gin.Default()
	router.Use(global.CORSMiddleware())
	router.Use(global.HandleGlobalException())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	includeYinPanRouter(router)
	return router
}
