package portal

import (
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/service/portal"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCartList(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	cartVO, err := cartService.GetCartList()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_CART_LIST, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func AddCartProduct(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	type RequestParams struct {
		ProductId int `json:"product_id"`
		Count int `json:"count"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Min(requestParams.ProductId, 1, "product_id").Message("商品ID必须大于0")
	valid.Min(requestParams.Count, 1, "count").Message("商品个数必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	cartVO, err := cartService.AddCartProduct(requestParams.ProductId, requestParams.Count)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_CART_PRODUCT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func UpdateCartProduct(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	type RequestParams struct {
		ProductId int `json:"product_id"`
		Count int `json:"count"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Min(requestParams.ProductId, 1, "product_id").Message("商品ID必须大于0")
	valid.Min(requestParams.Count, 1, "count").Message("商品个数必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	cartVO, err := cartService.UpdateCartProductCount(requestParams.ProductId, requestParams.Count)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_UPDATE_CART_PRODUCT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func DeleteCartProduct(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	type RequestParams struct {
		ProductIds string `json:"product_ids"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Required(requestParams.ProductIds, "product_ids").Message("商品ID不能为空")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	cartVO, err := cartService.DeleteCartProducts(requestParams.ProductIds)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_CART_PRODUCT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func SelectCartProduct(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	type RequestParams struct {
		ProductId int `json:"product_id"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Min(requestParams.ProductId, 1, "product_id").Message("商品ID必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	cartVO, err := cartService.SelectCartProduct(requestParams.ProductId)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_SELECT_CART_PRODUCT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func UnSelectCartProduct(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	type RequestParams struct {
		ProductId int `json:"product_id"`
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Min(requestParams.ProductId, 1, "product_id").Message("商品ID必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	cartVO, err := cartService.UnSelectCartProduct(requestParams.ProductId)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_UN_SELECT_CART_PRODUCT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func GetCartProductCount(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	count, err := cartService.GetCartProductCount()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_CART_PRODUCT_COUNT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, count)
}

func SelectAllCart(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	cartVO, err := cartService.SelectAllCart()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_SELECT_ALL_CART, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

func UnSelectAllCart(c *gin.Context)  {
	cartService := portal.CartProductVo{
		UserId: GetCurrentUser().ID,
	}
	appG := app.Gin{C: c}
	cartVO, err := cartService.UnSelectAllCart()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_UN_SELECT_ALL_CART, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cartVO)
}

