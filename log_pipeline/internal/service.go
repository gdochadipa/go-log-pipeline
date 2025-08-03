package internal

import (
	"context"
	"io"
	"log"

	"github.com/ochadipa/log_pipeline/proto/pb"
)


type ILogService interface {
	StreamLogs(ctx context.Context, stream pb.LogAggregator_StreamLogsServer) error
}

type logService struct {
	r LogRepository
}

func NewLogService(r LogRepository) ILogService {
	return &logService{r}
}

func (service *logService) StreamLogs(ctx context.Context, stream pb.LogAggregator_StreamLogsServer) error {
	// service.r.StreamLogs(ctx, )
	for {

		req, err := stream.Recv();
		if err == io.EOF {
			return stream.SendAndClose(&pb.LogResponse{Success: true})
		}
		if err != nil {
			return err
		}
		log.Printf("Received log from %s: %s", req.ServiceName, req.Message)
		// do submit logs here
	}
}
