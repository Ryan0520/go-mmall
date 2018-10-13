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
	appG := app.Gin{C: c}
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
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
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
	type RequestParams struct {
		Name string `json:"name"`
	}
	appG := app.Gin{C: c}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Required(requestParams.Name, "name").Message("分类名字不能为空")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	category, err := models.GetCategoryById(id)
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

func CheckAdminLogin(c *gin.Context) *models.User {
	loginUser := portal.CheckLogin(c)
	if loginUser == nil {
		return nil
	}

	appG := app.Gin{C: c}
	if loginUser.Role != models.Admin {
		appG.Response(http.StatusOK, e.ERROR_NOT_ADMIN, nil)
		return nil
	}

	return loginUser
}
