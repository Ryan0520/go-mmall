package backend

import (
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/routers/api/v1/portal"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AddCategory(c *gin.Context) {
	loginUser := portal.CheckLogin(c)
	if loginUser == nil {
		return
	}

	appG := app.Gin{C: c}
	if loginUser.Role != 1 {
		appG.Response(http.StatusOK, e.ERROR_NOT_ADMIN, nil)
		return
	}

	type RequestParams struct {
		ParentId int    `json:"parent_id"`
		Name     string `json:"name"`
	}

	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	parentId := 0
	if requestParams.ParentId > 0 {
		parentId = requestParams.ParentId
		exist, err := models.ExistCategoryByParentId(parentId)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_CHECK_CATEGORY_BY_PARENT_ID, nil)
			return
		}

		if !exist {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
			return
		}
	}

	valid := validation.Validation{}
	valid.Required(requestParams.Name, "name").Message("标签名称不能为空")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	category := models.Category{
		ParentId: parentId,
		Name:     requestParams.Name,
	}

	err = category.Save()
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_SAVE_CATEGORY, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetCategory(c *gin.Context) {
	loginUser := portal.CheckLogin(c)
	if loginUser == nil {
		return
	}

	appG := app.Gin{C: c}
	categoryId := c.DefaultQuery("category_id", "0")
	id, _ := com.StrTo(categoryId).Int()
	categories, err := models.GetCategoriesByParentId(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_GET_CATEGORY, nil)
		}
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, categories)
}

func UpdateCategory(c *gin.Context) {
	loginUser := portal.CheckLogin(c)
	if loginUser == nil {
		return
	}

	type RequestParams struct {
		CategoryId int    `json:"category_id"`
		Name       string `json:"name"`
	}
	appG := app.Gin{C: c}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(requestParams.Name, "name").Message("分类名字不能为空")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	category, err := models.GetCategoryById(requestParams.CategoryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_GET_CATEGORY, nil)
		}
		return
	}

	category.Name = requestParams.Name
	err = category.Update()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_UPDATE_CATEGORY, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetDeepCategoryId(c *gin.Context) {

}
