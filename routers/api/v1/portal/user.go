package portal

import (
	"fmt"
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/app"
	"github.com/Ryan0520/go-mmall/pkg/e"
	"github.com/Ryan0520/go-mmall/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 用户登录
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	type RequestParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var params RequestParams
	err := c.Bind(&params)
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	username := params.Username
	password := params.Password

	valid := validation.Validation{}
	valid.Required(username, "username").Message("用户名不能为空")
	valid.Required(password, "password").Message("密码不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.CheckUsernameExist(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_CHECK_USERNAME, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_USERNAME, nil)
		return
	}

	md5Pwd := util.MD5(password)
	user, err := models.QueryUserByUsernameAndPassword(username, md5Pwd)
	if err != nil {
		log.Error(err)
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_PASSWORD_NOT_CORRECT, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_QUERY_USER_BY_USERNAME_AND_PASSWORD, nil)
		}
		return
	}
	user.Password = ""
	util.WriteLoginUser(&user)
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

// 用户注册
func Register(c *gin.Context) {
	var user models.User
	appG := app.Gin{C: c}
	err := c.Bind(&user)
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(user.Username, "username").Message("用户名不能为空")
	valid.Required(user.Password, "password").Message("密码不能为空")
	valid.Phone(user.Phone, "phone").Message("手机号码格式不正确")
	valid.Email(user.Email, "email").Message("邮箱地址格式不正确")
	valid.Required(user.Question, "question").Message("问题不能为空")
	valid.Required(user.Answer, "answer").Message("答案不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.CheckUsernameExist(user.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_CHECK_USERNAME, nil)
		return
	}

	if exist {
		appG.Response(http.StatusOK, e.ERROR_EXIST_USERNAME, nil)
		return
	}

	exist, err = models.CheckPhoneExist(user.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_CHECK_PHONE, nil)
		return
	}

	if exist {
		appG.Response(http.StatusOK, e.ERROR_EXIST_PHONE, nil)
		return
	}

	exist, err = models.CheckEmailExist(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_CHECK_EMAIL, nil)
		return
	}

	if exist {
		appG.Response(http.StatusOK, e.ERROR_EXIST_EMAIL, nil)
		return
	}

	user.Role = 0
	user.Password = util.MD5(user.Password)
	err = user.Save()
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_SAVE_USER, nil)
		return
	}

	user.Password = ""
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

// 检查用户名或邮箱是否有效
func CheckValid(c *gin.Context) {
	appG := app.Gin{C: c}

	typeStr := c.DefaultQuery("type", "username")
	str := c.Query("str")

	valid := validation.Validation{}
	valid.Required(str, "str").Message("校验内容不能为空")
	if typeStr == "email" {
		valid.Email(str, " email").Message("邮箱地址格式不正确")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if typeStr == "username" {
		exist, err := models.CheckUsernameExist(str)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
			appG.Response(http.StatusOK, e.ERROR_CHECK_USERNAME, nil)
			return
		}

		if !exist {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}

		appG.Response(http.StatusOK, e.ERROR_EXIST_USERNAME, nil)
	} else if typeStr == "email" {
		exist, err := models.CheckEmailExist(str)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
			appG.Response(http.StatusOK, e.ERROR_CHECK_EMAIL, nil)
			return
		}

		if !exist {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}

		appG.Response(http.StatusOK, e.ERROR_EXIST_EMAIL, nil)
	}
}

// 获取登录的用户信息
func GetUserInfo(c *gin.Context) {
	loginUser := CheckLogin(c)
	if loginUser == nil {
		return
	}

	appG := app.Gin{C: c}
	username := loginUser.Username
	user, err := models.QueryUserWithUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_USER_NOT_EXIST, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_QUERY_USER, nil)
		}
		return
	}

	user.Password = ""
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

// 忘记密码中的获取问题
func ForgetGetQuestion(c *gin.Context) {
	username := c.Query("username")
	valid := validation.Validation{}
	valid.Required(username, "username").Message("用户名不能为空")

	appG := app.Gin{C: c}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.CheckUsernameExist(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_CHECK_USERNAME, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_USERNAME, nil)
		return
	}

	question, err := models.QueryQuestionWithUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_QUERY_QUESTION, nil)
		return
	}

	if question == "" || err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_QUESTION, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, question)
}

// 忘记密码中的检查答案
func ForgetCheckAnswer(c *gin.Context) {
	appG := app.Gin{C: c}
	type RequestParams struct {
		Username string
		Question string
		Answer   string
	}
	var requestParams RequestParams
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(requestParams.Username, "username").Message("用户名不能为空")
	valid.Required(requestParams.Question, "question").Message("找回密码的问题不能为空")
	valid.Required(requestParams.Answer, "answer").Message("找回密码的答案不能为空")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.CheckUsernameExist(requestParams.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_CHECK_USERNAME, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_USERNAME, nil)
		return
	}

	correct, err := models.CheckQuestionAndAnswerCorrect(requestParams.Username, requestParams.Question, requestParams.Answer)
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_CHECK_QUESTION_ANSWER, nil)
		return
	}

	if !correct {
		appG.Response(http.StatusOK, e.ERROR_QUESTION_ANSWER_NOT_CORRECT, nil)
		return
	}

	uuidStr := uuid.Must(uuid.NewV4())
	util.WriteResetPasswordToken(fmt.Sprintf("%s", uuidStr), requestParams.Username)
	data := make(map[string]interface{})
	data["token"] = uuidStr
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 忘记密码中的重置密码
func ForgetResetPassword(c *gin.Context) {
	type RequestParams struct {
		Token       string `json:"token"`
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}
	var requestParams RequestParams
	appG := app.Gin{C: c}
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	redisToken := util.ReadResetPasswordToken(requestParams.Username)
	if requestParams.Token != redisToken {
		appG.Response(http.StatusOK, e.ERROR_FORGET_RESET_PASSWORD_TOKEN, nil)
		return
	}

	user, err := models.QueryUserWithUsername(requestParams.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_QUERY_USER_BY_USERNAME_AND_PASSWORD, nil)
		return
	}

	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_EXIST, nil)
		return
	}

	user.Password = util.MD5(requestParams.NewPassword)
	err = user.Update()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_RESET_PASSWORD, nil)
		return
	}

	util.RemoveResetPasswordToken(requestParams.Username)
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 登录状态的重置密码
func ResetPassword(c *gin.Context) {
	loginUser := CheckLogin(c)
	if loginUser == nil {
		return
	}

	type RequestParams struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	var requestParams RequestParams
	appG := app.Gin{C: c}
	err := c.Bind(&requestParams)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	username := loginUser.Username
	md5Pwd := util.MD5(requestParams.OldPassword)
	user, err := models.QueryUserByUsernameAndPassword(username, md5Pwd)
	if err != nil && err != gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_QUERY_USER_BY_USERNAME_AND_PASSWORD, nil)
		return
	}

	if err == gorm.ErrRecordNotFound {
		appG.Response(http.StatusOK, e.ERROR_PASSWORD_NOT_CORRECT, nil)
		return
	}

	user.Password = util.MD5(requestParams.NewPassword)
	err = user.Update()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_RESET_PASSWORD, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	loginUser := CheckLogin(c)
	if loginUser == nil {
		return
	}

	var user models.User
	appG := app.Gin{C: c}
	err := c.Bind(&user)
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	if user.Email != "" {
		valid.Email(user.Email, "email").Message("邮箱地址格式不正确")
	}
	if user.Phone != "" {
		valid.Phone(user.Phone, "phone").Message("手机号码格式不正确")
	}
	valid.Required(user.Question, "question").Message("找回密码的问题不能为空")
	valid.Required(user.Answer, "answer").Message("找回密码的答案不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if user.Email != "" {
		exist, err := models.CheckEmailExist(user.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
			appG.Response(http.StatusOK, e.ERROR_CHECK_EMAIL, nil)
			return
		}

		if exist {
			log.Info("邮箱地址已经存在")
			appG.Response(http.StatusOK, e.ERROR_EXIST_EMAIL, nil)
			return
		}
	}

	if user.Phone != "" {
		exist, err := models.CheckPhoneExist(user.Phone)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
			appG.Response(http.StatusOK, e.ERROR_CHECK_PHONE, nil)
			return
		}

		if exist {
			log.Info("手机号码已存在")
			appG.Response(http.StatusOK, e.ERROR_EXIST_PHONE, nil)
			return
		}
	}

	username := loginUser.Username
	dbUser, err := models.QueryUserWithUsername(username)
	if err != nil {
		log.Error(err)
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.ERROR_USER_NOT_EXIST, nil)
		} else {
			appG.Response(http.StatusOK, e.ERROR_QUERY_USER, nil)
		}
		return
	}

	if user.Email != "" {
		dbUser.Email = user.Email
	}
	if user.Phone != "" {
		dbUser.Phone = user.Phone
	}
	dbUser.Question = user.Question
	dbUser.Answer = user.Answer
	err = dbUser.Update()
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusOK, e.ERROR_UPDATE_USER, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 退出登录
func Logout(c *gin.Context) {
	util.RemoveLoginUser()
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 检查当前是否有登录的用户
func CheckLogin(c *gin.Context) *models.User {
	appG := app.Gin{C: c}
	user := util.ReadLoginUser()
	if user == nil {
		appG.Response(http.StatusOK, e.ERROR_USER_NOT_LOGIN, nil)
		return nil
	}
	return user
}
