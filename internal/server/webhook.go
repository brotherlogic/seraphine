package server

import (
	"context"
	"log"
	"strings"

	ghwebhook_pb "github.com/brotherlogic/ghwebhook/proto/ghwebhook/v1"
	"github.com/brotherlogic/seraphine/internal/github"
)

// WebhookServer implements ghwebhook_pb.WebhookHandlerServer to handle incoming webhook events.
type WebhookServer struct {
	ghwebhook_pb.UnimplementedWebhookHandlerServer
	ghClient github.Client
}

// NewWebhookServer creates a new instance of WebhookServer with the given GitHub client.
func NewWebhookServer(ghClient github.Client) *WebhookServer {
	return &WebhookServer{
		ghClient: ghClient,
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
	if pr == nil || pr.GetAction() != "opened" {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}

	repoFullName := pr.GetRepository().GetFullName()
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return &ghwebhook_pb.WebhookResponse{}, nil
	}
	owner, repo := parts[0], parts[1]

	const reviewerEmail = "brotherlogicreviewer@gmail.com"
	isCollab, err := s.ghClient.IsCollaborator(ctx, owner, repo, reviewerEmail)
	if err != nil {
		log.Printf("Failed to check collaborator status for %s in %s: %v", reviewerEmail, repoFullName, err)
		return &ghwebhook_pb.WebhookResponse{}, nil
	}

	if isCollab {
		log.Printf("PR opened: repo=%s, pr_number=%d, title=%s, author=%s",
			repoFullName,
			pr.GetNumber(),
			pr.GetTitle(),
			pr.GetUser().GetLogin(),
		)
	}

	return &ghwebhook_pb.WebhookResponse{}, nil
}
