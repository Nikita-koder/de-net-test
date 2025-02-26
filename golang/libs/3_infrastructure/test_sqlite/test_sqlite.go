package test_sqlite

import (
	"de-net/libs/4_common/smart_context"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateUniqueSqliteInMemoryForUnitTests() *gorm.DB {
	dbName := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared", time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		// TODO: хотя лучше бы его тогда и возвращать
		Logger: smart_context.NewGormLoggerWrapper(smart_context.NewSmartContext()),
	})
	if err != nil {
		panic(fmt.Sprintf("error creating sqlite in memory: %v", err))
	}

	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")

	// Limit connections to avoid "database is locked" errors
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("error get db.DB(): %v", err))
	}
	sqlDB.SetMaxOpenConns(1) // Allow only 1 concurrent writer
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	return db
}
