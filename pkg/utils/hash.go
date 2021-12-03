package utils

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func Sha512Encode(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
