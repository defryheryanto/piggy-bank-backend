package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	trx_sql "github.com/defryheryanto/piggy-bank-backend/internal/transaction/sql"
	"gorm.io/gorm"
)

func setupParticipantService(db *gorm.DB) *transaction.ParticipantService {
	participantSQL := trx_sql.NewParticipantStorage(db)
	return transaction.NewParticipantService(participantSQL)
}

func setupTransactionService(db *gorm.DB) *transaction.TransactionService {
	trxSQL := trx_sql.NewTransactionStorage(db)
	participantService := setupParticipantService(db)

	return transaction.NewTransactionService(trxSQL, participantService)
}
