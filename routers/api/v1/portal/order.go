package portal

import (
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	protalServie "github.com/Ryan0520/go-mmall/service/portal"
	"github.com/Ryan0520/go-mmall/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func PayOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	type RequestParams struct {
		OrderNo int64 `json:"order_no"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(requestParams.OrderNo, 1, "order_no").Message("订单号必须大于零")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := protalServie.Order{
		OrderNo: requestParams.OrderNo,
		UserId:  user.ID,
	}
	payUrl, err := orderService.PayOrder()
	if err != nil || len(payUrl) == 0 {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ORDER, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_ORDER_PAY, nil)
		}
		return
	}

	data := make(map[string]interface{})
	data["order_no"] = requestParams.OrderNo
	data["pay_url"] = payUrl
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func QueryOrderPayStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	orderNo := com.StrTo(c.Query("order_no")).MustInt64()
	if orderNo <= 0 {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := protalServie.Order{
		OrderNo: orderNo,
		UserId:  user.ID,
	}
	result, err := orderService.QueryOrderPayStatus()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ORDER, nil);
		} else {
			appG.Response(http.StatusOK, e.ERROR_QUERY_ORDER_PAY_STATUS, nil)
		}
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

func CreateOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	type RequestParams struct {
		ShippingId int `json:"shipping_id"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(requestParams.ShippingId, 1, "shipping_id").Message("shipping_id必须大于零")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := protalServie.Order{
		ShippingId: requestParams.ShippingId,
		UserId:     user.ID,
	}
	orderVo, err := orderService.Create()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CREATE_ORDER, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orderVo)
}

func GetOrderCartProduct(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}
	orderService := protalServie.Order{
		UserId: user.ID,
	}
	orderCartProducts, err := orderService.GetOrderCartProducts()
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_GET_ORDER_CART_PRODUCT, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orderCartProducts)
}

func GetOrderList(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}
	orderService := protalServie.Order{
		UserId:   user.ID,
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	orders, err := orderService.GetOrderList()
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_GET_ORDER_LIST, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orders)
}

func GetOrderDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}
	orderNo := com.StrTo(c.Param("order_no")).MustInt64()
	if orderNo <= 0 {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := protalServie.Order{
		OrderNo: orderNo,
		UserId:  user.ID,
	}
	orderVo, err := orderService.GetOrderDetail()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ORDER_DETAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orderVo)
}

func CancelOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	type RequestParams struct {
		OrderNo int64 `json:"order_no"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(requestParams.OrderNo, 1, "order_no").Message("order_no必须大于零")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := protalServie.Order{
		OrderNo: requestParams.OrderNo,
		UserId:  user.ID,
	}
	if err := orderService.Cancel(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_CANCEL_ORDER, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 后台管理员接口
func ManageGetOrderList(c *gin.Context) {
	appG := app.Gin{C: c}
	orderService := protalServie.Order{
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	orders, err := orderService.ManageGetOrderList()
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_GET_ORDER_LIST, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orders)
}

func ManageOrderSearch(c *gin.Context) {
	appG := app.Gin{C: c}
	type RequestParams struct {
		OrderNo int64 `json:"order_no"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(requestParams.OrderNo, 1, "order_no").Message("order_no必须大于零")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := protalServie.Order{
		OrderNo: requestParams.OrderNo,
	}
	orderVo, err := orderService.ManageOrderSearch()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ORDER_DETAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orderVo)
}

func ManageGetOrderDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	orderNo := com.StrTo(c.Param("order_no")).MustInt64()
	if orderNo <= 0 {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	orderService := protalServie.Order{
		OrderNo: orderNo,
	}
	orderVo, err := orderService.ManageGetOrderDetail()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ORDER_DETAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, orderVo)
}

func ManageSendGoods(c *gin.Context)  {
	appG := app.Gin{C: c}
	type RequestParams struct {
		OrderNo int64 `json:"order_no"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(requestParams.OrderNo, 1, "order_no").Message("order_no必须大于零")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	orderService := protalServie.Order{
		OrderNo: requestParams.OrderNo,
	}
	if err := orderService.SendGoods(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_SEND_GOODS, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "发货成功")
}
