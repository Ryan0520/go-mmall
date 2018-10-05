package app

import (
	"github.com/astaxie/beego/validation"
	log "github.com/sirupsen/logrus"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Info(err.Key, err.Message)
	}

	return
}
