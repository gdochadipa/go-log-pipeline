package producer

import (
	"context"
	"time"

	"github.com/ochadipa/goroutine-practice/models"
)

/**
Simulates log messages from a fake service like ServiceA, ServiceB.
*/
func GenerateRawLogs(ctx context.Context, source string, out chan<- models.RawLog, errCh chan<- models.ErrorLog) {
  rawLog := models.RawLog{
 	Source: source,
	  Message: "DataA",
	  Time: time.Now(),
  }

  out <- rawLog
}
