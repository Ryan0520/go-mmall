package routers

import (
	"github.com/Ryan0520/go-mmall/middleware/auth"
	"github.com/Ryan0520/go-mmall/pkg/setting"
	"github.com/Ryan0520/go-mmall/routers/api/v1/backend"
	"github.com/Ryan0520/go-mmall/routers/api/v1/portal"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(setting.ServerSetting.RunMode)
	r.LoadHTMLGlob("templates/*")
	r.GET("/alipay/return", ReturnHandle)
	r.POST("/alipay/notify", NotifyHandle)

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

		product := apiV1.Group("/products/")
		product.Use(auth.Admin())
		{
			product.GET("/list", portal.GetProducts)
			product.GET("/", portal.GetProduct)
		}

		manager := apiV1.Group("/manage/")
		manager.Use(auth.Admin())
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

		cart := apiV1.Group("/cart/")
		cart.Use(auth.Login())
		{
			cart.GET("/list", portal.GetCartList)
			cart.POST("/add_product", portal.AddCartProduct)
			cart.POST("/update_product_num", portal.UpdateCartProduct)
			cart.POST("/delete_product", portal.DeleteCartProduct)
			cart.POST("/select_product", portal.SelectCartProduct)
			cart.POST("/un_select_product", portal.UnSelectCartProduct)
			cart.GET("/cart_product_count", portal.GetCartProductCount)
			cart.POST("/select_all", portal.SelectAllCart)
			cart.POST("/un_select_all", portal.UnSelectAllCart)
		}

		shipping := apiV1.Group("/shipping")
		shipping.Use(auth.Login())
		{
			shipping.POST("/", portal.AddShipping)
			shipping.DELETE("/:id", portal.DeleteShipping)
			shipping.PUT("/:id", portal.UpdateShipping)
			shipping.GET("/", portal.GetShipping)
			shipping.GET("/list", portal.GetShippingList)
		}

		order := apiV1.Group("/order/")
		order.Use(auth.Login())
		{
			order.POST("/pay", portal.PayOrder)
			order.GET("/order_pay_status", portal.QueryOrderPayStatus)
		}

	}

	return r
}
