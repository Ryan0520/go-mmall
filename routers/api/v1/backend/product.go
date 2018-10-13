package backend

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/pkg/upload"
	"github.com/Ryan0520/go-mmall/service/backend"
	"github.com/Ryan0520/go-mmall/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetProducts(c *gin.Context) {
	appG := app.Gin{C: c}
	productService := backend.Product{
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}

	total, err := productService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_PRODUCT_FAIL, nil)
		return
	}

	products, err := productService.GetAll()
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
	id := com.StrTo(c.Param("id")).MustInt()

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

	product, err := productService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PRODUCT, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}

func SearchProduct(c *gin.Context) {
	appG := app.Gin{C: c}
	pName := c.Query("productName")
	pId := com.StrTo(c.Query("productId")).MustInt()
	if len(pName) == 0 && pId == 0 {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	modelP := models.Product{
		Model: models.Model{
			ID: pId,
		},
		Name: pName,
	}

	productService := backend.Product{
		P:        &modelP,
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	products, total, err := productService.Search()
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_SEARCH_PRODUCT, nil)
		return
	}

	data := make(map[string]interface{})
	data["list"] = products
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func UploadProductImage(c *gin.Context) {
	appG := app.Gin{C: c}

	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	if image == nil {
		log.Info("image is nil")
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				log.Error(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				log.Error(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["url"] = upload.GetImageFullUrl(imageName)
				data["uri"] = upload.GetImagePath() + imageName
			}
		}
	}

	appG.Response(http.StatusOK, code, data)
}

func UploadProductRichTextImage(c *gin.Context) {
	appG := app.Gin{C: c}

	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	if image == nil {
		log.Info("image is nil")
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				log.Error(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				log.Error(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["file_path"] = upload.GetImageFullUrl(imageName)
			}
		}
	}

	appG.Response(http.StatusOK, code, data)
}

func UpdateProductSaleStatus(c *gin.Context) {
	type RequestParams struct {
		Id     int `json:"id"`
		Status int `json:"status"`
	}
	appG := app.Gin{C: c}
	var requestParams RequestParams
	err := c.ShouldBindJSON(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(requestParams.Id, 1, "id").Message("ID必须大于0")
	valid.Range(requestParams.Status, models.OnSale, models.OffSale, "status").Message("上架状态必须是1或者2")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	p := models.Product{
		Model: models.Model{
			ID: requestParams.Id,
		},
		Status: requestParams.Status,
	}
	productService := backend.Product{P: &p}
	dbP, err := productService.Get()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_GET_PRODUCT, nil)
		}
		return
	}

	if dbP.ID <= 0 {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
		return
	}

	err = productService.SetSaleStatus()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_UPDATE_PRODUCT_SALE_STATUS, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func SaveOrUpdate(c *gin.Context) {
	var p *models.Product
	appG := app.Gin{C: c}
	err := c.ShouldBind(&p)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(p.Name, "name").Message("产品名字不能为空")
	valid.Min(p.CategoryId, 1, "category_id").Message("分类ID必须大于0")
	valid.Range(p.Status, models.OnSale, models.OffSale, "status").Message("上架状态必须是1或者2")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	dbP, err := models.GetProduct(p.ID)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PRODUCT, nil)
		return
	}
	if dbP.ID == 0 {
		log.Error("产品 productId: ? 不存在", dbP.ID)
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
		return
	}

	if _, err := models.GetCategoryById(p.CategoryId); err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_GET_CATEGORY, nil)
		}
		return
	}

	if dbP.ID == 0 {
		if err = p.Save(); err != nil {
			appG.Response(http.StatusOK, e.ERROR_SAVE_PRODUCT,nil)
			return
		}
	} else {
		if err = p.Update(); err != nil {
			appG.Response(http.StatusOK, e.ERROR_UPDATE_PRODUCT,nil)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}