package user_service

import (
	errors "auth-server/internal/app/errors/types"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"
)

//DecodeEmailConfToken returns userID, tokenDeadTime
func decodeEmailConfToken(token, key string) (string, int64, error) {
	if len(token) < 0 {
		return "", 0, errors.ErrInvalidArgument.New("Invalid token!")
	}
	tokenHex, err := hex.DecodeString(token)
	if err != nil {
		return "", 0, errors.ErrInvalidArgument.New("Invalid token!")
	}
	keyHex, err := hex.DecodeString(key)
	if err != nil {
		return "", 0, err
	}
	block, err := aes.NewCipher(keyHex)
	if err != nil {
		return "", 0, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", 0, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(tokenHex) < nonceSize {
		return "", 0, errors.ErrInvalidArgument.New("Invalid token!")
	}
	nonce, tokenHex := tokenHex[:nonceSize], tokenHex[nonceSize:]
	decodedTokenBytes, err := aesGCM.Open(nil, nonce, tokenHex, nil)
	if err != nil {
		return "", 0, err
	}
	decodedToken := base64.StdEncoding.EncodeToString(decodedTokenBytes)
	data := strings.Split(decodedToken, "/")
	userID := data[0]
	seconds, err := time.ParseDuration(data[1] + "s")
	if err != nil {
		return "", 0, err
	}
	tokenDeadTime := time.Unix(int64(seconds.Seconds()), 0).Unix()
	return userID, tokenDeadTime, nil
}

func generateEmailConfToken(userID, key string) (string, error) {
	data := fmt.Sprintf("%s/%d", userID, time.Now().Add(24*time.Hour).Unix())
	//log.Printf("time is %s", data)
	keyHex, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}
	plainText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(keyHex)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plainText, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}
