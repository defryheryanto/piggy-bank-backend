package transaction_test

import (
	"testing"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTransactionService(db *gorm.DB) *transaction.TransactionService {
	transactionStorage := sql.NewTransactionStorage(db)
	participantService := setupParticipantService(db)

	return transaction.NewTransactionService(transactionStorage, participantService)
}

func setupDatabase(t *testing.T) *gorm.DB {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	return db
}

func TestCreateBasic(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupTransactionService(db)

		t.Run("insert transaction", func(t *testing.T) {
			payload := &transaction.CreateBasicTransactionPayload{
				UserId:          1,
				AccountId:       1,
				CategoryId:      1,
				TransactionDate: time.Now(),
				Description:     "Makan sate",
				Notes:           "Ber lima",
				Amount:          200000,
				TransactionType: transaction.ExpenseType,
			}

			err := service.CreateBasic(payload)
			assert.NoError(t, err)

			existing := &transaction.Transaction{}
			db.Order("transaction_id DESC").First(&existing)

			assert.Equal(t, payload.UserId, existing.UserId)
			assert.Equal(t, payload.AccountId, existing.AccountId)
			assert.Equal(t, payload.CategoryId, existing.CategoryId)
			assert.Equal(t, payload.TransactionDate.Unix(), existing.TransactionDate.Unix())
			assert.Equal(t, payload.Description, existing.Description)
			assert.Equal(t, payload.Notes, existing.Notes)
			assert.Equal(t, payload.Amount, existing.Amount)
			assert.Equal(t, payload.TransactionType, transaction.TransactionType(existing.TransactionType))
		})

		t.Run("insert participants", func(t *testing.T) {
			payload := &transaction.CreateBasicTransactionPayload{
				UserId:          1,
				AccountId:       1,
				CategoryId:      1,
				TransactionDate: time.Now(),
				Description:     "Makan sate",
				Notes:           "Ber lima",
				Amount:          200000,
				TransactionType: transaction.ExpenseType,
				Participants: []*transaction.CreateParticipantPayload{
					{
						Name:   "Kevin",
						Amount: 20000,
					},
					{
						Name:   "Jeremy",
						Amount: 45000,
					},
				},
			}

			err := service.CreateBasic(payload)
			assert.NoError(t, err)

			existingTrx := &transaction.Transaction{}
			db.Order("transaction_id DESC").First(&existingTrx)

			existing := []*transaction.Participant{}
			db.Where("transaction_id = ?", existingTrx.TransactionId).Find(&existing)

			assert.Equal(t, len(payload.Participants), len(existing))

			for _, pt := range payload.Participants {
				existing = []*transaction.Participant{}
				db.Where(
					"transaction_id = ? AND name = ? AND amount = ?",
					existingTrx.TransactionId,
					pt.Name,
					pt.Amount,
				).First(&existing)
				assert.Equal(t, 1, len(existing))
			}
		})
	})
}
