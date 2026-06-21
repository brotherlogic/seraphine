package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/brotherlogic/seraphine/proto"
	pstore_client "github.com/brotherlogic/pstore/client"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/github"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strings"
	"time"
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

func runSync(ctx context.Context, pClient pstore_client.PStoreClient, ghClient github.Client) error {
	state, err := config.ReadServerState(ctx, pClient)
	if err != nil {
		return fmt.Errorf("failed to read server state: %w", err)
	}

	invitations, err := ghClient.ListRepositoryInvitations(ctx)
	if err != nil {
		return fmt.Errorf("failed to list invitations: %w", err)
	}

	stateModified := false
	for _, inv := range invitations {
		err := ghClient.AcceptRepositoryInvitation(ctx, inv.ID)
		if err != nil {
			log.Printf("failed to accept invitation %d: %v", inv.ID, err)
			continue
		}

		repoFullName := inv.Repository.FullName
		// Check if it's already enrolled
		found := false
		for _, enrolled := range state.EnrolledRepositories {
			if enrolled == repoFullName {
				found = true
				break
			}
		}

		if !found {
			state.EnrolledRepositories = append(state.EnrolledRepositories, repoFullName)
			stateModified = true

			// Create issue in brotherlogic/devcontainer-manager
			_, err := ghClient.CreateIssue(ctx, "brotherlogic", "devcontainer-manager", fmt.Sprintf("Add %s to devcontainer manager", repoFullName), fmt.Sprintf("Automatically enrolled repository %s", repoFullName), []string{"seraphine-auto"})
			if err != nil {
				log.Printf("failed to create issue for %s: %v", repoFullName, err)
			}
		}
	}

	if stateModified {
		err = config.WriteServerState(ctx, pClient, state)
		if err != nil {
			return fmt.Errorf("failed to write server state: %w", err)
		}
	}

	// Validate ruleset configuration for all enrolled repositories
	for _, repoFullName := range state.EnrolledRepositories {
		parts := strings.Split(repoFullName, "/")
		if len(parts) != 2 {
			log.Printf("invalid repository full name: %s", repoFullName)
			continue
		}
		owner, repo := parts[0], parts[1]

		ruleset := &github.RulesetRequest{
			Name:        "Seraphine Enforced Rules",
			Target:      "branch",
			Enforcement: "active",
			Conditions: github.Conditions{
				RefName: github.RefName{
					Include: []string{"~DEFAULT_BRANCH"},
					Exclude: []string{},
				},
			},
			Rules: []github.Rule{
				{
					Type: "pull_request",
					Parameters: &github.RuleParameters{
						RequiredApprovingReviewCount:   1,
						DismissStaleReviewsOnPush:      true,
						RequireCodeOwnerReview:         false,
						RequireLastPushApproval:        true,
						RequiredReviewThreadResolution: true,
					},
				},
			},
		}

		err := ghClient.CreateRuleset(ctx, owner, repo, ruleset)
		if err != nil {
			log.Printf("failed to create ruleset for %s: %v", repoFullName, err)
		}
	}

	return nil
}

func RunWorkerLoop(ctx context.Context, pClient pstore_client.PStoreClient, ghClient github.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run once initially
	if err := runSync(ctx, pClient, ghClient); err != nil {
		log.Printf("sync error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := runSync(ctx, pClient, ghClient); err != nil {
				log.Printf("sync error: %v", err)
			}
		}
	}
}

func Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSeraphineServiceServer(grpcServer, &seraphineServer{})

	fmt.Printf("Starting Seraphine gRPC server on %s...\n", port)

	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Printf("GH_TOKEN is not set, skipping background worker")
	} else {
		pClient, err := pstore_client.GetClient()
		if err != nil {
			return fmt.Errorf("failed to get pstore client: %w", err)
		}
		ghClient := github.NewClient(token, nil)
		go RunWorkerLoop(context.Background(), pClient, ghClient, 1*time.Hour)
	}

	return grpcServer.Serve(lis)
}
