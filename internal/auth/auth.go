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
	GetByUsername(username string) *User
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
	userByUsername := s.userStorage.GetByUsername(payload.Username)
	if userByUsername != nil {
		return ErrUsernameHasTaken
	}

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

func (s *AuthService) Login(username, password string) (string, error) {
	currentUser := s.userStorage.GetByUsername(username)
	if currentUser == nil {
		return "", ErrInvalidCredential
	}
	isMatch, err := s.encryptor.Check(password, currentUser.Password)
	if err != nil {
		return "", err
	}
	if !isMatch {
		return "", ErrInvalidCredential
	}

	session := &AuthSession{
		UserID:   currentUser.UserID,
		Username: currentUser.Username,
	}
	token, err := s.tokenService.GenerateToken(session)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Authenticate(token string) (bool, error) {
	isValid, err := s.tokenService.CheckValidity(token)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

func (s *AuthService) GetCurrentUser(token string) (*AuthSession, error) {
	currentSession, err := s.tokenService.Parse(token)
	if err != nil {
		return nil, err
	}

	currentUser := s.userStorage.GetByUsername(currentSession.Username)
	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	return currentSession, nil
}
