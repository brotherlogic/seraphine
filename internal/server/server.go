package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type seraphineServer struct {
	pb.UnimplementedSeraphineServiceServer
}

func (s *seraphineServer) GetProjectState(ctx context.Context, req *pb.GetProjectStateRequest) (*pb.GetProjectStateResponse, error) {
	// TODO: implement business logic
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectState not implemented")
}

func (s *seraphineServer) RegisterProject(ctx context.Context, req *pb.RegisterProjectRequest) (*pb.RegisterProjectResponse, error) {
	// TODO: implement business logic
	return nil, status.Errorf(codes.Unimplemented, "method RegisterProject not implemented")
}

func Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSeraphineServiceServer(grpcServer, &seraphineServer{})

	fmt.Printf("Starting Seraphine gRPC server on %s...\n", port)
	return grpcServer.Serve(lis)
}
