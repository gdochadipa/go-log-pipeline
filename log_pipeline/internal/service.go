package internal

import "context"


type ILogService interface {
	StreamLogs(ctx context.Context) error
}

type logService struct {
	r LogRepository
}

func NewLogService(r LogRepository) ILogService {
	return &logService{r}
}

func (service *logService) StreamLogs(ctx context.Context) error {
	return nil
}
