package server

import (
	"context"
	"log"
	"strings"

	pstore_client "github.com/brotherlogic/pstore/client"
	ghwebhook_pb "github.com/brotherlogic/ghwebhook/proto/ghwebhook/v1"
	manager_pb "github.com/brotherlogic/devcontainer-manager/proto"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/github"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
)

// DevcontainerClient interface for devcontainer-manager interactions.
type DevcontainerClient interface {
	Up(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error)
	Down(ctx context.Context, in *manager_pb.DownRequest, opts ...grpc.CallOption) (*manager_pb.DownResponse, error)
	List(ctx context.Context, in *manager_pb.ListRequest, opts ...grpc.CallOption) (*manager_pb.ListResponse, error)
}

// WebhookServer implements ghwebhook_pb.WebhookHandlerServer to handle incoming webhook events.
type WebhookServer struct {
	ghwebhook_pb.UnimplementedWebhookHandlerServer
	ghClient           github.Client
	devcontainerClient DevcontainerClient
	pstoreClient       pstore_client.PStoreClient
}

// NewWebhookServer creates a new instance of WebhookServer with given clients.
func NewWebhookServer(ghClient github.Client, devcontainerClient DevcontainerClient, pstoreClient pstore_client.PStoreClient) *WebhookServer {
	return &WebhookServer{
		ghClient:           ghClient,
		devcontainerClient: devcontainerClient,
		pstoreClient:       pstoreClient,
	}
}

func (s *WebhookServer) ReceiveWebhook(ctx context.Context, req *ghwebhook_pb.WebhookEvent) (*ghwebhook_pb.WebhookResponse, error) {
	if req == nil || req.GetHeader() == nil {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}

	if req.GetHeader().GetEventType() != "pull_request" {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}

	pr := req.GetPullRequest()
	if pr == nil {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}

	action := pr.GetAction()
	if action != "opened" && action != "closed" {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}

	repoFullName := pr.GetRepository().GetFullName()
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}
	owner, repo := parts[0], parts[1]

	const reviewerEmail = "brotherlogicreviewer@gmail.com"
	if s.ghClient != nil {
		isCollab, err := s.ghClient.IsCollaborator(ctx, owner, repo, reviewerEmail)
		if err != nil {
			log.Printf("Failed to check collaborator status for %s in %s: %v", reviewerEmail, repoFullName, err)
			return &ghwebhook_pb.WebhookResponse{}, nil
		}
		if !isCollab {
			return &ghwebhook_pb.WebhookResponse{}, nil
		}
	}

	prNumber := pr.GetNumber()
	log.Printf("PR %s: repo=%s, pr_number=%d, title=%s, author=%s",
		action,
		repoFullName,
		prNumber,
		pr.GetTitle(),
		pr.GetUser().GetLogin(),
	)

	var state *pb.ServerState
	if s.pstoreClient != nil {
		st, err := config.ReadServerState(ctx, s.pstoreClient)
		if err != nil {
			log.Printf("Failed to read server state from pstore: %v", err)
			return &ghwebhook_pb.WebhookResponse{}, nil
		}
		state = st
	}

	if action == "opened" {
		if existing := config.FindPRContainer(state, repoFullName, prNumber); existing != nil {
			log.Printf("Container already exists for %s PR #%d: %s", repoFullName, prNumber, existing.GetContainerId())
			return &ghwebhook_pb.WebhookResponse{}, nil
		}

		if s.devcontainerClient != nil {
			upReq := &manager_pb.UpRequest{
				Repo: repoFullName,
				Identifier: &manager_pb.Identifier{
					Id: &manager_pb.Identifier_PrNumber{
						PrNumber: prNumber,
					},
				},
				Branch:  "",
				Prompt:  "",
				Harness: manager_pb.Harness_HARNESS_ANTIGRAVITY,
			}

			resp, err := s.devcontainerClient.Up(ctx, upReq)
			if err != nil {
				log.Printf("Failed to provision devcontainer via Up RPC: %v", err)
				return &ghwebhook_pb.WebhookResponse{}, nil
			}

			containerID := resp.GetConfig().GetId()
			if containerID != "" && s.pstoreClient != nil && state != nil {
				config.AddPRContainer(state, &pb.PRContainer{
					Repo:        repoFullName,
					PrNumber:    prNumber,
					ContainerId: containerID,
				})
				if err := config.WriteServerState(ctx, s.pstoreClient, state); err != nil {
					log.Printf("Failed to write updated server state to pstore: %v", err)
				}
			}
		}
	} else if action == "closed" {
		existing := config.FindPRContainer(state, repoFullName, prNumber)
		if existing == nil {
			log.Printf("No container tracked for %s PR #%d on closed action", repoFullName, prNumber)
			return &ghwebhook_pb.WebhookResponse{}, nil
		}

		containerID := existing.GetContainerId()
		if s.devcontainerClient != nil && containerID != "" {
			downReq := &manager_pb.DownRequest{
				Id: containerID,
			}
			_, err := s.devcontainerClient.Down(ctx, downReq)
			if err != nil {
				log.Printf("Failed to teardown devcontainer %s via Down RPC: %v", containerID, err)
			}
		}

		if s.pstoreClient != nil && state != nil {
			config.RemovePRContainer(state, repoFullName, prNumber)
			if err := config.WriteServerState(ctx, s.pstoreClient, state); err != nil {
				log.Printf("Failed to write updated server state to pstore after remove: %v", err)
			}
		}
	}

	return &ghwebhook_pb.WebhookResponse{}, nil
}
