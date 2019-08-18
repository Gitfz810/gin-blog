package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func GetTime() time.Time {
	//return time.Now().UnixNano() / int64(time.Millisecond)
	return time.Now()
}

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}