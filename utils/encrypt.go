package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptPasswd(str string) string {
	srcByte := []byte(str)
	sha256New := sha256.New()
	sha256Bytes := sha256New.Sum(srcByte)
	sha256String := hex.EncodeToString(sha256Bytes)
	return sha256String
}
