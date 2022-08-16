package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenService[T any] struct {
	signMethod jwt.SigningMethod
	secret     string
	expiryTime *time.Duration
}

func NewJWTTokenService[T any](signMethod jwt.SigningMethod, secret string, expiryTime *time.Duration) *JWTTokenService[T] {
	return &JWTTokenService[T]{
		signMethod: signMethod,
		secret:     secret,
		expiryTime: expiryTime,
	}
}

func (s *JWTTokenService[T]) GenerateToken(payload T) (string, error) {
	standardClaims := jwt.StandardClaims{}
	if s.expiryTime != nil {
		standardClaims.ExpiresAt = time.Now().Add(*s.expiryTime).Unix()
	}

	customClaims := &JWTCustomClaims[T]{
		Payload:        payload,
		StandardClaims: standardClaims,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenStr, err := jwtToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *JWTTokenService[T]) Parse(token string) (T, error) {
	_, payload, err := s.parseToken(token)
	if err != nil {
		var t T
		return t, err
	}

	return payload, nil
}

func (s *JWTTokenService[T]) CheckValidity(token string) (bool, error) {
	isValid, _, err := s.parseToken(token)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

func (s *JWTTokenService[T]) parseToken(token string) (validity bool, payload T, err error) {
	customClaims := JWTCustomClaims[T]{}
	jwtToken, err := jwt.ParseWithClaims(token, &customClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		var t T
		return false, t, err
	}

	return jwtToken.Valid, customClaims.Payload, nil
}
