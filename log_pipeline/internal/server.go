package internal

import (
	"fmt"
	"net"

	"github.com/ochadipa/log_pipeline/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


type grpcServer struct {
	pb.UnimplementedLogAggregatorServer
	service ILogService
}
func ListenGRPC(s ILogService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()

	pb.RegisterLogAggregatorServer(serv, &grpcServer{
		UnimplementedLogAggregatorServer: pb.UnimplementedLogAggregatorServer{},
		service: s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}
