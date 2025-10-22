package router

import (
	"github.com/gin-gonic/gin"
	"qimenpan/internal/handler/yinpan"
)

func includeYinPanRouter(engine *gin.Engine) {
	group := engine.Group("/qiMen")
	group.GET("/yinPan/v1", yinpan.NewPanHandler)
}
