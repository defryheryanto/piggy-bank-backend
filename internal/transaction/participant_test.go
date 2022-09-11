package transaction_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupParticipantService(db *gorm.DB) *transaction.ParticipantService {
	participantStorage := sql.NewParticipantStorage(db)

	return transaction.NewParticipantService(participantStorage)
}

func setupDatabase(t *testing.T) *gorm.DB {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	return db
}

func TestBulkCreate(t *testing.T) {
	db := setupDatabase(t)
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupParticipantService(db)

		t.Run("insert participants to db", func(t *testing.T) {
			payloads := []*transaction.CreateParticipantPayload{
				{
					Name:   "Kevin",
					Amount: 18000,
				},
				{
					Name:   "Jeremy",
					Amount: 45000,
				},
			}

			err := service.BulkCreate(1, payloads)
			assert.NoError(t, err)

			existing := []*transaction.Participant{}
			db.Where("transaction_id = 1").Find(&existing)

			assert.Equal(t, len(payloads), len(existing))

			db.Where("transaction_id = 1 AND name = ? AND amount = ?", payloads[0].Name, payloads[0].Amount).Find(&existing)
			assert.NotNil(t, existing)
			assert.Equal(t, 1, len(existing))

			db.Where("transaction_id = 1 AND name = ? AND amount = ?", payloads[1].Name, payloads[1].Amount).Find(&existing)
			assert.NotNil(t, existing)
			assert.Equal(t, 1, len(existing))
		})

		t.Run("return empty name error if name is empty", func(t *testing.T) {
			payloads := []*transaction.CreateParticipantPayload{
				{
					Name:   "",
					Amount: 18000,
				},
				{
					Name:   "Jeremy",
					Amount: 45000,
				},
			}
			err := service.BulkCreate(1, payloads)
			assert.ErrorIs(t, err, transaction.ErrEmptyParticipantName)
		})

		t.Run("return invalid amount err if amount is less than 0", func(t *testing.T) {
			payloads := []*transaction.CreateParticipantPayload{
				{
					Name:   "Kevin",
					Amount: 18000,
				},
				{
					Name:   "Jeremy",
					Amount: -1,
				},
			}
			err := service.BulkCreate(1, payloads)
			assert.ErrorIs(t, err, transaction.ErrInvalidAmount)
		})
	})
}
