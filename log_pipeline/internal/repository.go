package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/ochadipa/log_pipeline/models"
	"github.com/ochadipa/log_pipeline/processor"
	"github.com/ochadipa/log_pipeline/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type LogRepository interface {
	StreamLogs(ctx context.Context,log *pb.LogRequest) (*pb.LogResponse,error)
}

type logRepository struct {
	pool *processor.WorkerPool
}

func NewRepository( workerPool *processor.WorkerPool) LogRepository {
	return &logRepository{workerPool}
}


func (r *logRepository) StreamLogs(ctx context.Context, log *pb.LogRequest) (*pb.LogResponse, error) {
	fmt.Println("StreamLogs running")
	logEntry := &models.LogType{
		Service: log.GetServiceName(),
		Timestamp: log.GetTimestamp().AsTime(),
		Level: log.GetLevel(),
		Message: log.GetMessage(),
		Metadata: map[string]interface{}{
			"dummy":"dummy",
		},
	}

	job := processor.Job{
		Log: *logEntry,
	}

	err := r.pool.Submit(ctx, job)

	if err != nil {
		switch err {
			case context.DeadlineExceeded :
				return nil, status.Error(codes.DeadlineExceeded, "request timed out, server might be busy")
			case context.Canceled :
				return nil, status.Error(codes.Canceled, "request canceled")
			default :
				return nil, status.Errorf(codes.ResourceExhausted, "server is overloaded, please try again later: %v", err)

		}
	}

	return &pb.LogResponse{ Success: true}, nil
}
