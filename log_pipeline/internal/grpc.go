package internal

import (
	"context"
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

	// grpc start new server
	serv := grpc.NewServer()

	pb.RegisterLogAggregatorServer(serv, &grpcServer{
		UnimplementedLogAggregatorServer: pb.UnimplementedLogAggregatorServer{},
		service: s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}


// this another option i think for listening grpc combine with gorutine
func ListenGRPC2(ctx context.Context, s ILogService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// grpc start new server
	serv := grpc.NewServer()

	pb.RegisterLogAggregatorServer(serv, &grpcServer{
		UnimplementedLogAggregatorServer: pb.UnimplementedLogAggregatorServer{},
		service: s,
	})
	reflection.Register(serv)

	go func() {
		fmt.Printf("gRPC server listening on: %s", lis.Addr().String())
		if err := serv.Serve(lis); err != nil {
			fmt.Printf("Failed to serve gRPC server: %v", err)
		}
	}()

	// if channel ctx.Done() never filled
	<- ctx.Done();
	fmt.Println("Shutting down gRPC server...")
    serv.GracefulStop()
    fmt.Println("gRPC server stopped.")
    return nil;
}
