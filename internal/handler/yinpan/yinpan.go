package yinpan

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qimenpan/pkg/pan"
	"time"
)

func NewPanHandler(ctx *gin.Context) {
	datetime := ctx.DefaultQuery("datetime", time.Now().Format(time.DateTime))
	parser, err := time.Parse(time.DateTime, datetime)
	if err != nil {
		panic(err)
	}
	p := pan.New(parser)
	panData := p.Parse()
	ctx.JSON(http.StatusOK, panData)
}
