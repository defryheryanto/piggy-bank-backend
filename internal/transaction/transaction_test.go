package transaction_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/storage"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTransactionService(db *gorm.DB) *transaction.TransactionService {
	transactionStorage := sql.NewTransactionStorage(db)
	participantService := setupParticipantService(db)
	sqlManager := storage.NewSQLManager(db)

	return transaction.NewTransactionService(transactionStorage, participantService, sqlManager)
}

func setupDatabase(t *testing.T) *gorm.DB {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	return db
}

func TestCreateBasic_InsertTransaction(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()

	service := setupTransactionService(db)

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.Participant{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
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

		err := service.CreateBasic(ctx, payload)
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
}

func TestCreateBasic_InsertParticipants(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()

	service := setupTransactionService(db)

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.Participant{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
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

		err := service.CreateBasic(ctx, payload)
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
}

func TestCreateBasic_RollbackIfFailed(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.Participant{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
		errorString := "mock error"
		mockParticipantStorage := sql.NewParticipantMockStorage(func() error {
			return fmt.Errorf(errorString)
		})
		mockTransactionStorage := sql.NewTransactionMockStorage(
			func() error { return nil },
			func() error { return nil },
			func() error { return nil },
		)
		participantService := transaction.NewParticipantService(mockParticipantStorage)
		sqlManager := storage.NewSQLManager(db)
		transactionService := transaction.NewTransactionService(mockTransactionStorage, participantService, sqlManager)

		payload := &transaction.CreateBasicTransactionPayload{
			UserId:          1,
			AccountId:       1,
			CategoryId:      1,
			TransactionDate: time.Now(),
			Description:     "test rollback",
			Notes:           "test rollback",
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

		var c int64
		db.Model(&transaction.Transaction{}).Order("transaction_id DESC").Count(&c)

		err := transactionService.CreateBasic(ctx, payload)
		assert.Error(t, err)

		existingTrx := &transaction.Transaction{}
		db.Order("transaction_id DESC").First(&existingTrx)
		assert.Empty(t, existingTrx)

		existingParticipants := []*transaction.Participant{}
		db.Where("transaction_id = ?", existingTrx.TransactionId).Find(&existingParticipants)
		assert.Empty(t, existingParticipants)
	})
}

func TestCreateTransfer_InsertDetails(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()
	service := setupTransactionService(db)

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.TransferDetail{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
		payload := &transaction.CreateTransferPayload{
			UserId:          1,
			SourceAccountId: 1,
			TargetAccountId: 2,
			TransactionDate: time.Now(),
			Description:     "Transfer ke tabungan",
			Notes:           "Tabungan bulan septemer",
			Amount:          500000,
		}

		err := service.CreateTransfer(ctx, payload)
		assert.NoError(t, err)

		existingTrx := &transaction.Transaction{}
		db.Last(&existingTrx)
		assert.Equal(t, payload.UserId, existingTrx.UserId)
		assert.Equal(t, payload.SourceAccountId, existingTrx.AccountId)
		assert.Equal(t, payload.TransactionDate.Unix(), existingTrx.TransactionDate.Unix())
		assert.Equal(t, payload.Description, existingTrx.Description)
		assert.Equal(t, payload.Notes, existingTrx.Notes)
		assert.Equal(t, payload.Amount, existingTrx.Amount)

		existingDetail := &transaction.TransferDetail{}
		db.Last(&existingDetail)
		assert.Equal(t, existingTrx.TransactionId, existingDetail.TransactionId)
		assert.Equal(t, payload.TargetAccountId, existingDetail.TargetAccountId)
	})
}

func TestCreateTransfer_RollbackIfError(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()
	errorString := "mock error"
	mockParticipantStorage := sql.NewParticipantMockStorage(func() error {
		return fmt.Errorf(errorString)
	})
	mockTransactionStorage := sql.NewTransactionMockStorage(
		func() error { return nil },
		func() error { return fmt.Errorf(errorString) },
		func() error { return nil },
	)
	participantService := transaction.NewParticipantService(mockParticipantStorage)
	sqlManager := storage.NewSQLManager(db)
	service := transaction.NewTransactionService(mockTransactionStorage, participantService, sqlManager)

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.TransferDetail{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
		payload := &transaction.CreateTransferPayload{
			UserId:          1,
			SourceAccountId: 1,
			TargetAccountId: 2,
			TransactionDate: time.Now(),
			Description:     "Transfer ke tabungan",
			Notes:           "Tabungan bulan septemer",
			Amount:          500000,
		}

		err := service.CreateTransfer(ctx, payload)
		assert.Error(t, err)
		assert.Equal(t, errorString, err.Error())

		existingTrx := &transaction.Transaction{}
		db.Last(&existingTrx)
		assert.Empty(t, existingTrx)

		existingDetail := &transaction.TransferDetail{}
		db.Last(&existingDetail)
		assert.Empty(t, existingDetail)
	})
}

func TestCreateSaving_InsertDetails(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()
	service := setupTransactionService(db)

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.SavingDetail{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
		payload := &transaction.CreateSavingPayload{
			UserId:          1,
			SourceAccountId: 1,
			TargetAccountId: 2,
			TransactionDate: time.Now(),
			Description:     "Nabung Bibit",
			Notes:           "Tabungan Bulan ini",
			Amount:          5000000,
		}

		err := service.CreateSaving(ctx, payload)
		assert.NoError(t, err)

		trx := &transaction.Transaction{}
		db.Last(&trx)
		assert.Equal(t, payload.UserId, trx.UserId)
		assert.Equal(t, payload.SourceAccountId, trx.AccountId)
		assert.Equal(t, payload.TransactionDate.Unix(), trx.TransactionDate.Unix())
		assert.Equal(t, payload.Description, trx.Description)
		assert.Equal(t, payload.Notes, trx.Notes)
		assert.Equal(t, payload.Amount, trx.Amount)

		savingDetail := &transaction.SavingDetail{}
		db.Last(&savingDetail)
		assert.Equal(t, trx.TransactionId, savingDetail.TransactionId)
		assert.Equal(t, payload.TargetAccountId, savingDetail.TargetAccountId)
	})
}

func TestCreateSaving_RollbackIfError(t *testing.T) {
	db := setupDatabase(t)
	ctx := context.TODO()
	errorString := "mock error"
	mockParticipantStorage := sql.NewParticipantMockStorage(func() error {
		return fmt.Errorf(errorString)
	})
	mockTransactionStorage := sql.NewTransactionMockStorage(
		func() error { return nil },
		func() error { return nil },
		func() error { return fmt.Errorf(errorString) },
	)
	participantService := transaction.NewParticipantService(mockParticipantStorage)
	sqlManager := storage.NewSQLManager(db)
	service := transaction.NewTransactionService(mockTransactionStorage, participantService, sqlManager)

	usedTableNames := []string{
		transaction.Transaction{}.TableName(),
		transaction.SavingDetail{}.TableName(),
	}

	test.TruncateAfterTest(t, db, usedTableNames, func() {
		payload := &transaction.CreateSavingPayload{
			UserId:          1,
			SourceAccountId: 1,
			TargetAccountId: 2,
			TransactionDate: time.Now(),
			Description:     "Nabung Bibit",
			Notes:           "Tabungan Bulan ini",
			Amount:          5000000,
		}

		err := service.CreateSaving(ctx, payload)
		assert.Error(t, err)
		assert.Equal(t, errorString, err.Error())

		existingTrx := &transaction.Transaction{}
		db.Last(&existingTrx)
		assert.Empty(t, existingTrx)

		existingDetail := &transaction.SavingDetail{}
		db.Last(&existingDetail)
		assert.Empty(t, existingDetail)
	})
}
