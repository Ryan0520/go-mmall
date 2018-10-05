package util

import (
	"crypto/md5"
	"fmt"
)

const Salt = "jWcgWrFyiKKiyl6gvAqy"

func MD5(content string) string {
	content = content + Salt
	data := []byte(content)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

//func main() {
//	md5Str := MD5("123456")
//	log.Info("md5: %s", md5Str)
//}
