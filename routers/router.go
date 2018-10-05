package routers

import (
	"github.com/Ryan0520/go-mmall/pkg/setting"
	"github.com/Ryan0520/go-mmall/routers/api/v1/backend"
	"github.com/Ryan0520/go-mmall/routers/api/v1/portal"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(setting.ServerSetting.RunMode)

	apiV1 := r.Group("/api/v1/")
	{
		user := apiV1.Group("/user/")
		{
			user.POST("/login", portal.Login)
			user.POST("/register", portal.Register)
			user.GET("/check_valid", portal.CheckValid)
			user.GET("/user_info", portal.GetUserInfo)
			user.GET("/forget_get_question", portal.ForgetGetQuestion)
			user.POST("/forget_check_answer", portal.ForgetCheckAnswer)
			user.POST("/forget_reset_password", portal.ForgetResetPassword)
			user.POST("/reset_password", portal.ResetPassword)
			user.POST("/update_information", portal.UpdateUserInfo)
			user.GET("/logout", portal.Logout)
		}

		manager := apiV1.Group("/manage/")
		{
			manager.GET("/category", backend.GetCategory)
			manager.POST("/category", backend.AddCategory)
			manager.PUT("/category", backend.UpdateCategory)
			manager.GET("/deep_category_ids", backend.GetDeepCategoryId)
		}
	}

	return r
}
