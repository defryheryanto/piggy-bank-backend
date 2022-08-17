package aes_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/encrypt"
	"github.com/defryheryanto/piggy-bank-backend/internal/encrypt/aes"
)

func TestNewAESEncryptor(t *testing.T) {
	_, err := aes.NewAESEncryptor("secret_need_to_be_32_characters!")
	if err != nil {
		t.Errorf("error new aes encryptor %v", err)
	}

	_, err = aes.NewAESEncryptor("not_32_characters!")
	if err == nil {
		t.Errorf("should return error if secret not 32 chars length")
	}
}

func TestAESEncryptorShouldMatch(t *testing.T) {
	var encryptor encrypt.Encryptor
	encryptor, err := aes.NewAESEncryptor("secret_need_to_be_32_characters!")
	if err != nil {
		t.Errorf("error new aes encryptor %v", err)
	}

	text := "this text needs to be encrypted"
	encrypted, err := encryptor.Encrypt(text)
	if err != nil {
		t.Errorf("error encrypting text %v", err)
	}

	isMatch, err := encryptor.Check(text, encrypted)
	if err != nil {
		t.Errorf("error checking encrypted text %v", err)
	}

	if !isMatch {
		t.Errorf("initial text and decrypted value is different")
	}
}

func TestAESEncryptorShouldNotMatch(t *testing.T) {
	var encryptor encrypt.Encryptor
	encryptor, err := aes.NewAESEncryptor("secret_need_to_be_32_characters!")
	if err != nil {
		t.Errorf("error new aes encryptor %v", err)
	}

	text := "this text needs to be encrypted"
	encrypted, err := encryptor.Encrypt(text)
	if err != nil {
		t.Errorf("error encrypting text %v", err)
	}

	isMatch, err := encryptor.Check("different text", encrypted)
	if err != nil {
		t.Errorf("error checking encrypted text %v", err)
	}

	if isMatch {
		t.Errorf("initial text and decrypted value is same")
	}
}
