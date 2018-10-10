package app

import (
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	if data == nil {
		g.C.JSON(httpCode, gin.H{
			"code": errCode,
			"msg":  e.GetMsg(errCode),
		})
	} else {
		g.C.JSON(httpCode, gin.H{
			"code": errCode,
			"msg":  e.GetMsg(errCode),
			"data": data,
		})
	}
}
