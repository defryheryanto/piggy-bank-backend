package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type AESEncryptor struct {
	secret []byte
}

func NewAESEncryptor(secret string) (*AESEncryptor, error) {
	secretBytes := []byte(secret)
	if len(secretBytes) != 32 {
		return nil, errors.New("secret must be 32 characters length")
	}
	return &AESEncryptor{secretBytes}, nil
}

func (e *AESEncryptor) Encrypt(text string) (string, error) {
	cph, err := aes.NewCipher(e.secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cph)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	textBytes := []byte(text)
	encryptedBytes := gcm.Seal(nonce, nonce, textBytes, nil)

	return base64.URLEncoding.EncodeToString(encryptedBytes), nil
}

func (e *AESEncryptor) Check(realString, encryptedString string) (bool, error) {
	c, err := aes.NewCipher(e.secret)
	if err != nil {
		return false, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return false, err
	}

	encryptedBytes, err := base64.URLEncoding.DecodeString(encryptedString)
	if err != nil {
		return false, err
	}
	nonceSize := gcm.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return false, nil
	}

	nonce, encryptedBytes := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
	textBytes, err := gcm.Open(nil, nonce, encryptedBytes, nil)
	if err != nil {
		return false, err
	}

	text := string(textBytes[:])
	if realString != text {
		return false, nil
	}

	return true, nil
}
