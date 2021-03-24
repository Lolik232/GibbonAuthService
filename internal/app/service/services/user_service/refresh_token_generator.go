package user_service

import (
	"crypto/rand"
	"encoding/hex"
)

func generateRefreshToken() string {
	bytes := make([]byte, 32)
	token := ""
	rand.Read(bytes)

	for _, v := range bytes {
		token += string(v)
	}
	token = hex.EncodeToString(bytes)
	return token
}
