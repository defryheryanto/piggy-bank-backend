package auth

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/encrypt"
	"github.com/defryheryanto/piggy-bank-backend/internal/token"
)

type User struct {
	UserID   int    `gorm:"primaryKey;autoIncrement;column:user_id" json:"user_id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}

func (u *User) TableName() string { return "users" }

type UserRepository interface {
	Create(*User) error
}

type AuthSession struct {
	UserID   int
	Username string
}

type AuthService struct {
	userStorage  UserRepository
	tokenService token.TokenIService[*AuthSession]
	encryptor    encrypt.Encryptor
}

func NewAuthService(
	userStorage UserRepository,
	tokenService token.TokenIService[*AuthSession],
	encryptor encrypt.Encryptor,
) *AuthService {
	return &AuthService{userStorage, tokenService, encryptor}
}

func (s *AuthService) Register(payload *User) error {
	encryptedPassword, err := s.encryptor.Encrypt(payload.Password)
	if err != nil {
		return err
	}
	user := &User{
		Username: payload.Username,
		Password: encryptedPassword,
	}

	err = s.userStorage.Create(user)
	if err != nil {
		return err
	}

	return nil
}
