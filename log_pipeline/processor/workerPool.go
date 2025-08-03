package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/ochadipa/log_pipeline/models"
	"gorm.io/gorm"
)

type Job struct {
	Log models.LogType
}

type WorkerPool struct {
	maxWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
	done       chan error
	ctx        context.Context    // To cancel the pool's context
	cancel     context.CancelFunc // Context for the pool's lifecycle
	db         *gorm.DB
}

func NewWorkerPool(db *gorm.DB, maxWorkers, queueSize int) *WorkerPool {
	jobQueue := make(chan Job, queueSize)
	done := make(chan error, maxWorkers) // Create buffered channel
	workerPool := &WorkerPool{
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		done:       done,
		db:         db,
	}

	return workerPool
}

// worker goroutine
func (pool *WorkerPool) Start(ctx context.Context) {
	pool.wg.Add(pool.maxWorkers)
	for i := range pool.maxWorkers {
		// when combine with loop, every goroutine have their work id
		go func(workId int) {
			// what ever happend, if the job queue less than worker, the worker still run
			defer pool.wg.Done()
			fmt.Printf("Worker started %d", workId)
			for job := range pool.jobQueue {
				// processLog(workId, job.Log)
				if err := pool.InsertLogs(ctx, &job.Log); err != nil {
					fmt.Printf("Worker %d: error inserting log: %v\n", workId, err)
					pool.done <- err // Send error to the done channel
					return
				}
				pool.done <- nil // Signal successful completions
			}
			fmt.Printf("worker done", workId)
		}(i + 1)
	}
}

// processLog simulates the work of processing a log entry.
func processLog(workerID int, log models.LogType) {
	fmt.Printf("Worker %d processing log: %s\n", workerID, log.Message)
	// Simulate a slow operation like writing to a database or file.
	time.Sleep(200 * time.Millisecond)
}

func (pool *WorkerPool) Submit(ctx context.Context, job Job) error {
	// this is select channel
	select {
	// when channel context is done() then we can return context an error
	case <-ctx.Done():
		return ctx.Err()
	case pool.jobQueue <- job:
		return nil
	default:
		return nil
	}
}

func (pool *WorkerPool) Shutdown(ctx context.Context) {
	close(pool.jobQueue)
	pool.wg.Wait()
}

func (pool *WorkerPool) InsertLogs(ctx context.Context, log *models.LogType) error {
	tx := pool.db.Begin()
	sql, args, err := sq.Insert("logs").Columns("timestamp", "service", "level", "message", "metadata").Values(log.Timestamp, log.Service, log.Level, log.Message, log.Metadata).ToSql()

	if err != nil {
		fmt.Errorf("failed to insert ledger entry: %w", err)
		return err
	}
	result := tx.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		tx.Rollback()
		fmt.Errorf("failed to insert ledger entry: %w", result.Error)
		return result.Error
	}

	// Use context with the transaction if GORM supports it for Exec
	// (GORM typically uses context on the DB instance for queries)
	// For raw SQL with Exec, you might need to ensure the underlying driver respects context.
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert log: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no rows affected for log insert: %v", log.Message)
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	// fmt.Printf("Log inserted successfully: %s\n", log.Message) // For

	fmt.Printf("done.\n")

	return nil
}
