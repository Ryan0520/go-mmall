package util

import (
	"github.com/Ryan0520/go-mmall/pkg/setting"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int  {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * GetPageSize(c)
	}
	return result
}

func GetPageSize(c *gin.Context) int {
	result := setting.AppSetting.PageSize
	pageSize, _ := com.StrTo(c.Query("pageSize")).Int()
	if pageSize > 0 {
		result = pageSize
	}
	return result
}