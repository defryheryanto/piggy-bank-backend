package auth_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	user_storage "github.com/defryheryanto/piggy-bank-backend/internal/auth/storage"
	"github.com/defryheryanto/piggy-bank-backend/internal/encrypt/aes"
	jwt_service "github.com/defryheryanto/piggy-bank-backend/internal/token/jwt"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

func setupService(t *testing.T, db *gorm.DB) *auth.AuthService {
	userStorage := user_storage.NewUserStorage(db)
	tokenService := jwt_service.NewJWTTokenService[*auth.AuthSession](jwt.SigningMethodHS256, "testsecret", nil)
	encryptor, err := aes.NewAESEncryptor("secret_need_to_be_32_characters!")
	if err != nil {
		t.Errorf("failed to initialize encryptor %v", err)
	}

	return auth.NewAuthService(userStorage, tokenService, encryptor)
}

func truncateAuthTables(t *testing.T, db *gorm.DB) {
	user := &auth.User{}
	tables := []string{
		user.TableName(),
	}
	test.TruncateTables(t, db, tables)
}

func TestRegister(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	truncateAuthTables(t, db)
	service := setupService(t, db)

	t.Run("insert user to db", func(t *testing.T) {
		user := &auth.User{
			Username: "TestUser",
			Password: "123123",
		}
		err := service.Register(user)
		if err != nil {
			t.Errorf("failed to register user %v", err)
		}

		var insertedData *auth.User
		db.Where("username = ?", user.Username).First(&insertedData)
		if insertedData.Username != user.Username {
			t.Errorf("user not registered")
		}
	})

	t.Run("return error if username is already taken", func(t *testing.T) {
		user := &auth.User{
			Username: "TestUser",
			Password: "123123",
		}
		err := service.Register(user)
		if err != auth.UsernameHasTakenError {
			t.Errorf("should raise error if username is exists")
		}
	})
	truncateAuthTables(t, db)
}
