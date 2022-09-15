package storage

import (
	"context"

	"gorm.io/gorm"
)

type key string

const databaseContextKey key = "database_context"

func DatabaseFromContext(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(databaseContextKey).(*gorm.DB)
	if !ok {
		return nil
	}

	return db
}

type SQLManager struct {
	db *gorm.DB
}

func NewSQLManager(db *gorm.DB) *SQLManager {
	return &SQLManager{db}
}

func (m *SQLManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	db := m.db.Begin()
	if db.Error != nil {
		return db.Error
	}
	ctx = context.WithValue(ctx, databaseContextKey, db)

	err := fn(ctx)
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}
