package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ochadipa/log_pipeline/models"
)

type Job struct {
	log models.LogType
}

type WorkerPool struct {
	maxWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
}

func NewWorkerPool(maxWorkers, queueSize int) *WorkerPool {
	jobQueue := make(chan Job, queueSize)
	workerPool := &WorkerPool{
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
	}

	return workerPool
}

// worker goroutine
func (pool *WorkerPool) Start() {
	pool.wg.Add(pool.maxWorkers)
	for i := range pool.maxWorkers {
		go func(workId int) {
			defer pool.wg.Done()
			fmt.Printf("Worker started %d", workId)
			for job := range pool.jobQueue {
				processLog(workId, job.log)
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
	select {
		case <-ctx.Done() :
			return ctx.Err()
		case pool.jobQueue <- job :
			return nil
		default:
			return nil
	}
}

func (pool *WorkerPool) Shutdown(ctx context.Context) {
	close(pool.jobQueue)
	pool.wg.Wait()
}
