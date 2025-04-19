package util

import (
	"crypto/md5"
	"fmt"
)

// MakeMd5 生成MD5哈希
func MakeMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}