package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTCustomClaims[T any] struct {
	Payload T
	jwt.StandardClaims
}
