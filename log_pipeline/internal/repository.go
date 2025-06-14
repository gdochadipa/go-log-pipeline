package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)


var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type LogRepository interface {
	Close()
	StreamLogs(ctx context.Context) error
}

type logRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) LogRepository {
	return &logRepository{db}
}


func (r *logRepository) StreamLogs(ctx context.Context) error {
	fmt.Println("StreamLogs running")
	return nil
}

func (r *logRepository) Close() {
	r.db.Close()
}
