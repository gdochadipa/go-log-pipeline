package processor

import (
	"context"

	"github.com/ochadipa/goroutine-practice/models"
)


func ParsedLogs(ctx context.Context, id int, in <-chan models.RawLog, out chan <- models.ParsedLog, errCh <-chan models.ErrorLog){
	go func(){
		for val := range in {
			parsed := models.ParsedLog{
				Source: val.Source,
				Message: val.Message,
				Level: "INFO",
				Timestamp: val.Time,
			}

			out <- parsed
		}
	}()

	go func(){
		for val := range errCh {
			parsed := models.ParsedLog{
				Source: val.Source,
				Message: val.Error.Error(),
				Level: "ERROR",
				Timestamp: val.Time,
			}

			out <- parsed
		}
	}()
}
