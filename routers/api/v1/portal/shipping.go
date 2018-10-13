package portal

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/service/portal"
	"github.com/Ryan0520/go-mmall/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func AddShipping(c *gin.Context)  {
	appG := app.Gin{C: c}
	var s models.Shipping
	err := c.Bind(&s)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(s.ReceiverName, "receiver_name").Message("收件人名字不能为空")
	valid.Phone(s.ReceiverPhone, "receiver_phone").Message("收件人手机号码格式不正确")
	if s.ReceiverMobile != "" {
		valid.Mobile(s.ReceiverMobile, "receiver_mobile").Message("收件人电话号码格式不正确")
	}
	valid.Required(s.ReceiverProvince, "receiver_province").Message("收件人省份不能为空")
	valid.Required(s.ReceiverCity, "receiver_city").Message("收件人省市不能为空")
	valid.Required(s.ReceiverDistrict, "receiver_district").Message("收件人区域不能为空")
	valid.Required(s.ReceiverAddress, "receiver_address").Message("收件人地址不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	shippingService := portal.Shipping{
		UserId: user.ID,
		ReceiverName: s.ReceiverName,
		ReceiverPhone: s.ReceiverPhone,
		ReceiverMobile: s.ReceiverMobile,
		ReceiverProvince: s.ReceiverProvince,
		ReceiverCity: s.ReceiverCity,
		ReceiverAddress: s.ReceiverAddress,
		ReceiverDistrict: s.ReceiverDistrict,
		ReceiverZip: s.ReceiverZip,
	}
	shippingId, err := shippingService.Add()
	if err != nil || shippingId == 0 {
		appG.Response(http.StatusOK, e.ERROR_ADD_SHIPPING, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, shippingId)
}

func DeleteShipping(c *gin.Context)  {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	shippingId := com.StrTo(c.Param("id")).MustInt()
	shippingService := portal.Shipping{
		ID: shippingId,
	}
	err := shippingService.Delete()
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_DELETE_SHIPPING, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func UpdateShipping(c *gin.Context)  {
	appG := app.Gin{C: c}
	var s models.Shipping
	err := c.Bind(&s)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	valid := validation.Validation{}
	shippingId := com.StrTo(c.Param("id")).MustInt()
	valid.Min(shippingId, 1, "id").Message("shipping_id必须大于零")
	if s.ReceiverName != "" {
		valid.Required(s.ReceiverName, "receiver_name").Message("收件人名字不能为空")
	}
	if s.ReceiverPhone != "" {
		valid.Phone(s.ReceiverPhone, "receiver_phone").Message("收件人手机号码格式不正确")
	}
	if s.ReceiverMobile != "" {
		valid.Mobile(s.ReceiverMobile, "receiver_mobile").Message("收件人电话号码格式不正确")
	}
	if s.ReceiverProvince != "" {
		valid.Required(s.ReceiverProvince, "receiver_province").Message("收件人省份不能为空")
	}
	if s.ReceiverCity != "" {
		valid.Required(s.ReceiverCity, "receiver_city").Message("收件人省市不能为空")
	}
	if s.ReceiverDistrict != "" {
		valid.Required(s.ReceiverDistrict, "receiver_district").Message("收件人区域不能为空")
	}
	if s.ReceiverAddress != "" {
		valid.Required(s.ReceiverAddress, "receiver_address").Message("收件人地址不能为空")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	shippingService := portal.Shipping{
		ID: shippingId,
		UserId: user.ID,
		ReceiverName: s.ReceiverName,
		ReceiverPhone: s.ReceiverPhone,
		ReceiverMobile: s.ReceiverMobile,
		ReceiverProvince: s.ReceiverProvince,
		ReceiverCity: s.ReceiverCity,
		ReceiverAddress: s.ReceiverAddress,
		ReceiverDistrict: s.ReceiverDistrict,
		ReceiverZip: s.ReceiverZip,
	}
	err = shippingService.Update()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_UPDATE_SHIPPING, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetShipping(c *gin.Context)  {
	appG := app.Gin{C: c}
	shippingId := com.StrTo(c.Query("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(shippingId, 1, "id").Message("shipping_id必须大于零")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	shippingService := portal.Shipping{
		ID: shippingId,
		UserId: user.ID,
	}
	modelShipping, err := shippingService.Get()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_SHIPPING, nil)
 		} else {
			appG.Response(http.StatusOK, e.ERROR_GET_SHIPPING, nil)
		}
		return
	}

	if modelShipping.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_SHIPPING, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, modelShipping)
}

func GetShippingList(c *gin.Context)  {
	appG := app.Gin{C: c}
	user := GetCurrentUser()
	if user.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return
	}

	shippingService := portal.Shipping{
		UserId: user.ID,
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}

	list, total, err := shippingService.GetList()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_SHIPPINGS, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}