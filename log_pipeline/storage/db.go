package storage

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func NewPostgressConnection(user *string, password *string, host *string, dbName *string) (*sql.DB, error) {

	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", user,  password, host, dbName)

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
	return sqlDB, nil
}
