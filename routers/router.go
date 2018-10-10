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
			manager.GET("/categories/:id", backend.GetCategory)
			manager.POST("/categories", backend.AddCategory)
			manager.PUT("/categories/:id", backend.UpdateCategory)
			manager.GET("/categories_deep_ids", backend.GetDeepCategoryId)

			manager.GET("/products/:id", backend.GetProduct)
			manager.GET("/products", backend.GetProducts)
			manager.GET("/product/search", backend.SearchProduct)
			manager.POST("/product/sale_status", backend.UpdateProductSaleStatus)
			manager.POST("/products", backend.SaveOrUpdate)
			manager.POST("/product/upload", backend.UploadProductImage)
			manager.POST("/product/rich_text_img_upload", backend.UploadProductRichTextImage)
		}
	}

	return r
}
