package test

import (
	"fmt"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func RunTestWithDB(db *gorm.DB, t *testing.T, testFunction func(t *testing.T, db *gorm.DB)) {
	db = db.Begin()
	testFunction(t, db)
	db.Rollback()
}

func TruncateAfterTest(t *testing.T, db *gorm.DB, tableNames []string, actionFunc func()) {
	actionFunc()
	res := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", strings.Join(tableNames, ", ")))
	if res.Error != nil {
		fmt.Println(res.Error)
	}
}
