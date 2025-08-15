package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMACSHA256Hex(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
