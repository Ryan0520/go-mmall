package portal

import (
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/service/portal"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func PayOrder(c *gin.Context)  {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	type RequestParams struct{
		OrderNo int `json:"order_no"`
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

	orderService := portal.Order{
		OrderNo: requestParams.OrderNo,
		UserId: user.ID,
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

func QueryOrderPayStatus(c *gin.Context)  {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	orderNo := com.StrTo(c.Query("order_no")).MustInt()
	if orderNo <= 0 {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	orderService := portal.Order{
		OrderNo: orderNo,
		UserId: user.ID,
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