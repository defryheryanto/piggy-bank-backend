package main

import (
	"github.com/defryheryanto/piggy-bank-backend/config"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	auth_storage "github.com/defryheryanto/piggy-bank-backend/internal/auth/sql"
	"github.com/defryheryanto/piggy-bank-backend/internal/encrypt/aes"
	jwt_service "github.com/defryheryanto/piggy-bank-backend/internal/token/jwt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

func SetupAuthService(db *gorm.DB) *auth.AuthService {
	userStorage := auth_storage.NewUserStorage(db)
	tokenService := jwt_service.NewJWTTokenService[*auth.AuthSession](jwt.SigningMethodHS256, config.JWTSecretKey(), nil)
	encryptor, err := aes.NewAESEncryptor(config.AESSecretKey())
	userConfigService := SetupUserConfigService(db)
	if err != nil {
		panic(err)
	}
	return auth.NewAuthService(userStorage, tokenService, encryptor, userConfigService)
}

func SetupUserConfigService(db *gorm.DB) *auth.UserConfigService {
	configStorage := auth_storage.NewUserConfigStorage(db)

	return auth.NewUserConfigService(configStorage)
}
