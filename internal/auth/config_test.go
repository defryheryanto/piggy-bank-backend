package auth_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	auth_storage "github.com/defryheryanto/piggy-bank-backend/internal/auth/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupUserConfigService(db *gorm.DB) *auth.UserConfigService {
	configStorage := auth_storage.NewUserConfigStorage(db)

	return auth.NewUserConfigService(configStorage)
}

func setupDatabase(t *testing.T) *gorm.DB {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	return db
}

func TestCreateDefault(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupUserConfigService(db)

		t.Run("create default user config", func(t *testing.T) {
			err := service.CreateDefault(1)
			assert.NoError(t, err)

			cfg := &auth.UserConfig{}
			db.Where("user_id = 1").First(&cfg)

			//check default config values
			assert.Equal(t, 1, cfg.MonthlyStartDate)
		})
	})
}
