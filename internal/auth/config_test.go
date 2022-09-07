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

			def := auth.DefaultUserConfig()
			assert.Equal(t, def.MonthlyStartDate, cfg.MonthlyStartDate)
		})
	})
}

func TestGetByUserId(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupUserConfigService(db)

		t.Run("return existing user config", func(t *testing.T) {
			userId := 1
			initCfg := &auth.UserConfig{
				UserId:           userId,
				MonthlyStartDate: 25,
			}

			db.Create(&initCfg)

			cfg, err := service.GetByUserId(userId)
			assert.NoError(t, err)
			assert.Equal(t, initCfg.MonthlyStartDate, cfg.MonthlyStartDate)
		})

		t.Run("create default user config if not exists and return it", func(t *testing.T) {
			cfg, err := service.GetByUserId(2)
			assert.NoError(t, err)
			def := auth.DefaultUserConfig()
			assert.Equal(t, def.MonthlyStartDate, cfg.MonthlyStartDate)
		})
	})
}

func TestUpdate(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupUserConfigService(db)

		t.Run("update config if exists", func(t *testing.T) {
			userId := 1
			initCfg := &auth.UserConfig{
				UserId:           userId,
				MonthlyStartDate: 25,
			}

			db.Create(&initCfg)

			payload := &auth.UpdateUserConfigPayload{
				UserId:           initCfg.UserId,
				MonthlyStartDate: 5,
			}

			err := service.Update(payload)
			assert.NoError(t, err)

			cfg, err := service.GetByUserId(userId)
			assert.NoError(t, err)
			assert.Equal(t, payload.MonthlyStartDate, cfg.MonthlyStartDate)
		})

		t.Run("create config if not exists", func(t *testing.T) {
			userId := 3
			payload := &auth.UpdateUserConfigPayload{
				UserId:           userId,
				MonthlyStartDate: 5,
			}

			err := service.Update(payload)
			assert.NoError(t, err)

			cfg, err := service.GetByUserId(userId)
			assert.NoError(t, err)
			assert.Equal(t, payload.MonthlyStartDate, cfg.MonthlyStartDate)
		})
	})
}
