package util

import (
	"crypto/md5"
	"encoding/hex"
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

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

//func main() {
//	md5Str := MD5("123456")
//	log.Info("md5: %s", md5Str)
//}
