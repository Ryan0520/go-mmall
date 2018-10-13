package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/Ryan0520/go-mmall/pkg/setting"
)

func MD5BySalt(content string) string {
	content = content + setting.AppSetting.PasswordSalt
	data := []byte(content)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}