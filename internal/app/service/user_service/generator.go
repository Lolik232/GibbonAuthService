package user_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateEmailConfToken(userID, email string) (string, string) {
	s := hmac.New(sha256.New, []byte("superkey"))
	s.Write([]byte("message"))
	sum := s.Sum(nil)
	msg := "net"
	if hmac.Equal(sum,
		[]byte("message")) {
		msg = "lol"
	}
	return base64.StdEncoding.EncodeToString(sum), msg
}
