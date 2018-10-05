package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Model

	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Role     int    `json:"role"`
}

func CheckUsernameExist(username string) (bool, error) {
	var count int
	err := db.Where("username = ?", username).Find(&User{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("CheckUsernameExist error: %v", err)
	}
	return count > 0, nil
}

func CheckPhoneExist(phone string) (bool, error) {
	var count int
	err := db.Where("phone = ?", phone).Find(&User{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("CheckPhoneExist error: %v", err)
	}
	return count > 0, nil
}

func CheckEmailExist(email string) (bool, error) {
	var count int
	err := db.Where("email= ?", email).Find(&User{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("CheckEmailExist error: %v", err)
	}
	return count > 0, nil
}

func QueryUserByUsernameAndPassword(username string, password string) (User, error) {
	var user User
	err := db.Where("username = ? AND password = ?", username, password).Find(&User{}).Scan(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("QueryUserByUsernameAndPassword error: %v", err)
	}
	return user, err
}

func QueryQuestionWithUsername(username string) (string, error) {
	type Result struct {
		Question string
	}
	var result Result
	err := db.Select("question").Where("username = ?", username).Find(&User{}).Scan(&result).Error
	return result.Question, err
}

func CheckQuestionAndAnswerCorrect(username string, question string, answer string) (bool, error) {
	var count int
	err := db.Where("username = ? AND question = ? AND answer = ?", username, question, answer).Find(&User{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("CheckQuestionAndAnswerCorrect error: %v", err)
		return false, err
	}
	return count > 0, nil
}

func QueryUserWithUsername(username string) (User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("QueryUserWithUsername error: %v", err)
	}
	return user, err
}

func (user *User) Save() error {
	err := db.Create(user).Error
	return err
}

func (user *User) Update() error {
	err := db.Model(&user).Updates(user).Error
	return err
}
