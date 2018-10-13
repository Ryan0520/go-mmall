package auth

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/routers/api/v1/portal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUser := portal.CheckLogin(c)
		if loginUser == nil {
			c.Abort()
			return
		}

		if loginUser.Role != models.Admin {
			appG := app.Gin{C: c}
			appG.Response(http.StatusOK, e.ERROR_NOT_ADMIN, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUser := portal.CheckLogin(c)
		if loginUser == nil {
			c.Abort()
			return
		}

		c.Next()
	}
}
