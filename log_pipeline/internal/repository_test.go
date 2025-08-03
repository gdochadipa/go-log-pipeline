package internal

import (
	// "context"
	"sync"
	"testing"

	"github.com/ochadipa/log_pipeline/models"
	"github.com/ochadipa/log_pipeline/processor"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type MockWorkerPool struct {
	maxWorkers int
	jobQueue   chan processor.Job
	wg         sync.WaitGroup
	db         *gorm.DB
}

func NewMockWorkerPool(db *gorm.DB, maxWorker, queueSize int) *MockWorkerPool {
	pool := &MockWorkerPool{
		jobQueue: make(chan processor.Job, queueSize),
		maxWorkers: maxWorker,
		db: db,
	}
	return pool
}

// worker goroutine
// func (pool *MockWorkerPool) Start(ctx context.Context) {
// 	pool.wg.Add(pool.maxWorkers)
// 	for i := range pool.maxWorkers {
// 		go func(workId int) {
// 			defer pool.wg.Done()
// 			fmt.Printf("Worker started %d", workId)
// 			for job := range pool.jobQueue {
// 				// processLog(workId, job.Log)
// 				// pool.InsertLogs(ctx, &job.Log)
// 			}
// 			fmt.Printf("worker done", workId)
// 		}(i + 1)
// 	}
// }


func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// AutoMigrate will create the table, fields, indexes, and constraints based on the LogType model
	err = db.AutoMigrate(&models.LogType{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
