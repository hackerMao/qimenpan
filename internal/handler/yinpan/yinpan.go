package yinpan

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qimenpan/pkg/pan"
	"strings"
	"time"
)

func NewPanHandler(ctx *gin.Context) {
	datetime := ctx.DefaultQuery("datetime", time.Now().Format(time.DateTime))
	datetime = strings.Replace(datetime, "T", " ", -1)
	parser, err := time.Parse(time.DateTime, datetime)
	if err != nil {
		panic(err)
	}
	p := pan.New(parser)
	panData := p.Parse()
	ctx.JSON(http.StatusOK, panData)
}
