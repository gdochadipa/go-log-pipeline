package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/ochadipa/log_pipeline/models"
)

type LogQueue struct {
	queueData chan models.LogType
	wg        sync.WaitGroup
}

// setup new queue
func NewLogQueue(bufferSize int) *LogQueue {
	lq := &LogQueue{
		queueData: make(chan models.LogType, bufferSize),
	}


	// add waiting on wait group
	lq.wg.Add(1)
	// run Worker function in go
	go lq.Worker()

	return lq

}

// worker is the internal goroutine that continuously processes logs.
// It uses a for...range loop that automatically exits when the channel is closed.
func (lq *LogQueue) Worker() {
	defer lq.wg.Done()

	for log := range lq.queueData {
		// Process the log message.
		fmt.Printf("Processing log: %s\n", log.Message)

		fmt.Println(time.Now().Format(time.RFC3339), log.Message)
		// To simulate work and see the queue in action.
		time.Sleep(100 * time.Millisecond)
	}
}

// add value to queue
func (lq *LogQueue) Enqueue(log models.LogType) {
	select {
	case lq.queueData <- log:
	default:
		fmt.Println("queue has already full, will throw the lock")
	}
}

func (lq *LogQueue) ShutDown() {
	close(lq.queueData)
	lq.wg.Wait()

}

// ini berarti
func Run() {
	lq := NewLogQueue(100)

	for i := 0; i < 10; i++ {
		// ini berarti bakal run dibawah go orutine ?
		lq.Enqueue(models.LogType{
			Service:   "example-service",
			Timestamp: time.Now(),
			Level:     "INFO",
			Message:   "pesan",
			Metadata: map[string]interface{}{
				"iteration": i,
			},
		})
	}
	time.Sleep(1 * time.Second)
	lq.ShutDown()
}

/*
 *
 * pake worker pool
 */
