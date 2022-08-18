package test

import (
	"testing"

	"gorm.io/gorm"
)

func RunTestWithDB(db *gorm.DB, t *testing.T, testFunction func(t *testing.T, db *gorm.DB)) {
	db = db.Begin()
	testFunction(t, db)
	db.Rollback()
}
