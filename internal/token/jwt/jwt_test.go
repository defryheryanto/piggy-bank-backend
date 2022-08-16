package jwt_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/token"
	jwt_service "github.com/defryheryanto/piggy-bank-backend/internal/token/jwt"
	"github.com/dgrijalva/jwt-go"
)

type payload struct {
	UserID   int
	Username string
}

func TestJWTTokenService(t *testing.T) {
	var service token.TokenIService[*payload]
	p := &payload{1, "Defry Heryanto"}
	service = jwt_service.NewJWTTokenService[*payload](jwt.SigningMethodHS256, "secretkey", nil)

	//check generate process
	token, err := service.GenerateToken(p)
	if err != nil {
		t.Errorf("generate token failed - %v", err)
	}

	if token == "" {
		t.Errorf("empty token generated")
	}

	//check parse process
	parsedPayload, err := service.Parse(token)
	if err != nil {
		t.Errorf("error parsing token - %v", err)
	}

	if parsedPayload.UserID != p.UserID || parsedPayload.Username != p.Username {
		t.Errorf("parsed jwt payload is different from initial payload")
	}

	//check validity
	isValid, err := service.CheckValidity(token)
	if err != nil {
		t.Errorf("error when checking token validity %v", err)
	}

	if !isValid {
		t.Errorf("jwt token should be valid")
	}
}
