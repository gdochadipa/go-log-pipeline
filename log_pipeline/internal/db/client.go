package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func NewPostgressConnection(user *string, password *string, host *string, dbName *string) (*gorm.DB, error) {

	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", *user,  *password, *host, *dbName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	// Test the connection
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected")
	return db, nil
}

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) *DB {
return &DB{db}
}
