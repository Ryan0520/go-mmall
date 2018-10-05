package util

import (
	"encoding/json"
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/gredis"
	log "github.com/sirupsen/logrus"
)

var LoginToken = "mmall_login_token"
var TokenPrefix = "Token_"

func WriteLoginUser(user *models.User) error {
	return gredis.Set(LoginToken, user, 60*10)
}

func RemoveLoginUser() {
	result, err := gredis.Delete(LoginToken)
	if err != nil {
		log.Error(err)
	}
	log.Info("RemoveLoginUser result: %d", result)
}

func ReadLoginUser() *models.User {
	var user *models.User
	data, err := gredis.Get(LoginToken)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal(data, &user)
	return user
}

func WriteResetPasswordToken(token string, username string) error {
	return gredis.Set(TokenPrefix+username, token, 60*10)
}

func ReadResetPasswordToken(username string) string {
	var token string
	data, err := gredis.Get(TokenPrefix + username)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal(data, &token)
	return token
}

func RemoveResetPasswordToken(username string) {
	result, err := gredis.Delete(TokenPrefix + username)
	if err != nil {
		log.Error(err)
	}
	log.Info("RemoveResetPasswordToken result: %d", result)
}
