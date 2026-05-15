package client

import (
	"context"
	"time"

	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetProjectState(ctx context.Context, serverAddr string, projectName string, currentVersion string) (*pb.GetProjectStateResponse, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewSeraphineServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return client.GetProjectState(ctx, &pb.GetProjectStateRequest{
		ProjectName:    projectName,
		CurrentVersion: currentVersion,
	})
}

func RegisterProject(ctx context.Context, serverAddr string, projectName string, repositoryURL string) (*pb.RegisterProjectResponse, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewSeraphineServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return client.RegisterProject(ctx, &pb.RegisterProjectRequest{
		ProjectName:   projectName,
		RepositoryUrl: repositoryURL,
	})
}
