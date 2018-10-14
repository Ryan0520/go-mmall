package portal

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/service/backend"
	"github.com/Ryan0520/go-mmall/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProducts(c *gin.Context) {
	appG := app.Gin{C: c}

	productService := backend.Product{
		Keyword:  c.DefaultQuery("keyword", ""),
		OrderBy:  c.DefaultQuery("order_by", ""),
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}

	total, err := productService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_PRODUCT_FAIL, nil)
		return
	}

	products, err := productService.GetAllFilterOffSale()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PRODUCT_LIST, nil)
		return
	}

	data := make(map[string]interface{})
	data["list"] = products
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func GetProduct(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Query("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	modelP := models.Product{
		Model: models.Model{ID: id},
	}
	productService := backend.Product{
		P: &modelP,
	}

	product, err := productService.GetFilterOffSale()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_PRODUCT_OFFSALE_OR_DELETED, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}
