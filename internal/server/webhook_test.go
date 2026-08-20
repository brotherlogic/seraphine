package server

import (
	"context"
	"errors"
	"testing"

	pstore_client "github.com/brotherlogic/pstore/client"
	ghwebhook_pb "github.com/brotherlogic/ghwebhook/proto/ghwebhook/v1"
	manager_pb "github.com/brotherlogic/devcontainer-manager/proto"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/github"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
)

type mockGHClient struct {
	isCollaboratorFunc func(ctx context.Context, owner, repo, user string) (bool, error)
}

func (m *mockGHClient) ListRepositoryInvitations(ctx context.Context) ([]*github.RepositoryInvitation, error) {
	return nil, nil
}

func (m *mockGHClient) AcceptRepositoryInvitation(ctx context.Context, invitationID int64) error {
	return nil
}

func (m *mockGHClient) CreateRuleset(ctx context.Context, owner, repo string, ruleset *github.RulesetRequest) error {
	return nil
}

func (m *mockGHClient) CreateIssue(ctx context.Context, owner, repo string, title, body string, labels []string) (*github.IssueResponse, error) {
	return nil, nil
}

func (m *mockGHClient) IsCollaborator(ctx context.Context, owner, repo, user string) (bool, error) {
	if m.isCollaboratorFunc != nil {
		return m.isCollaboratorFunc(ctx, owner, repo, user)
	}
	return false, nil
}

func (m *mockGHClient) ListOpenPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	return nil, nil
}

func (m *mockGHClient) GetPullRequestDetails(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error) {
	return nil, nil
}

func (m *mockGHClient) GetCommitCheckStatus(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error) {
	return github.CheckStatusUnknown, nil
}

type mockDevcontainerClient struct {
	upFunc   func(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error)
	downFunc func(ctx context.Context, in *manager_pb.DownRequest, opts ...grpc.CallOption) (*manager_pb.DownResponse, error)
}

func (m *mockDevcontainerClient) Up(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error) {
	if m.upFunc != nil {
		return m.upFunc(ctx, in, opts...)
	}
	return &manager_pb.UpResponse{}, nil
}

func (m *mockDevcontainerClient) Down(ctx context.Context, in *manager_pb.DownRequest, opts ...grpc.CallOption) (*manager_pb.DownResponse, error) {
	if m.downFunc != nil {
		return m.downFunc(ctx, in, opts...)
	}
	return &manager_pb.DownResponse{}, nil
}

func TestReceiveWebhook_OpenedPROfEligibleRepo(t *testing.T) {
	collabCalled := false
	upCalled := false

	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			collabCalled = true
			if owner == "brotherlogic" && repo == "seraphine" && user == "brotherlogicreviewer@gmail.com" {
				return true, nil
			}
			return false, nil
		},
	}

	mockDev := &mockDevcontainerClient{
		upFunc: func(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error) {
			upCalled = true
			if in.GetRepo() != "brotherlogic/seraphine" {
				t.Errorf("expected repo brotherlogic/seraphine, got %s", in.GetRepo())
			}
			if in.GetIdentifier().GetPrNumber() != 141 {
				t.Errorf("expected PR number 141, got %d", in.GetIdentifier().GetPrNumber())
			}
			if in.GetHarness() != manager_pb.Harness_HARNESS_ANTIGRAVITY {
				t.Errorf("expected harness HARNESS_ANTIGRAVITY, got %v", in.GetHarness())
			}
			return &manager_pb.UpResponse{
				Config: &manager_pb.DevcontainerConfig{
					Id: "container-141",
				},
			}, nil
		},
	}

	pClient := pstore_client.GetTestClient()
	srv := NewWebhookServer(mockClient, mockDev, pClient)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "opened",
				Number: 141,
				Title:  "Test PR",
				User: &ghwebhook_pb.User{
					Login: "testuser",
				},
				Repository: &ghwebhook_pb.Repository{
					FullName: "brotherlogic/seraphine",
				},
			},
		},
	}

	resp, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("ReceiveWebhook returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("ReceiveWebhook returned nil response")
	}
	if !collabCalled {
		t.Errorf("Expected IsCollaborator to be called, but it was not")
	}
	if !upCalled {
		t.Errorf("Expected Devcontainer.Up to be called, but it was not")
	}

	// Verify state saved to pstore
	state, err := config.ReadServerState(context.Background(), pClient)
	if err != nil {
		t.Fatalf("Failed to read server state: %v", err)
	}
	c := config.FindPRContainer(state, "brotherlogic/seraphine", 141)
	if c == nil || c.GetContainerId() != "container-141" {
		t.Errorf("Expected container container-141 in pstore, got: %v", c)
	}
}

func TestReceiveWebhook_OpenedPRIdempotency(t *testing.T) {
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			return true, nil
		},
	}

	upCalled := false
	mockDev := &mockDevcontainerClient{
		upFunc: func(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error) {
			upCalled = true
			return &manager_pb.UpResponse{}, nil
		},
	}

	pClient := pstore_client.GetTestClient()
	// Pre-populate state
	state := &pb.ServerState{
		PrContainers: []*pb.PRContainer{
			{Repo: "brotherlogic/seraphine", PrNumber: 141, ContainerId: "container-141"},
		},
	}
	_ = config.WriteServerState(context.Background(), pClient, state)

	srv := NewWebhookServer(mockClient, mockDev, pClient)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "opened",
				Number: 141,
				Repository: &ghwebhook_pb.Repository{
					FullName: "brotherlogic/seraphine",
				},
			},
		},
	}

	_, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if upCalled {
		t.Errorf("Devcontainer.Up should not be called for duplicate opened PR")
	}
}

func TestReceiveWebhook_ClosedPROfTrackedContainer(t *testing.T) {
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			return true, nil
		},
	}

	downCalled := false
	mockDev := &mockDevcontainerClient{
		downFunc: func(ctx context.Context, in *manager_pb.DownRequest, opts ...grpc.CallOption) (*manager_pb.DownResponse, error) {
			downCalled = true
			if in.GetId() != "container-141" {
				t.Errorf("expected container id container-141, got %s", in.GetId())
			}
			return &manager_pb.DownResponse{}, nil
		},
	}

	pClient := pstore_client.GetTestClient()
	state := &pb.ServerState{
		PrContainers: []*pb.PRContainer{
			{Repo: "brotherlogic/seraphine", PrNumber: 141, ContainerId: "container-141"},
		},
	}
	_ = config.WriteServerState(context.Background(), pClient, state)

	srv := NewWebhookServer(mockClient, mockDev, pClient)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "closed",
				Number: 141,
				Repository: &ghwebhook_pb.Repository{
					FullName: "brotherlogic/seraphine",
				},
			},
		},
	}

	_, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !downCalled {
		t.Errorf("Expected Devcontainer.Down to be called, but it was not")
	}

	// Verify state removed from pstore
	st, _ := config.ReadServerState(context.Background(), pClient)
	c := config.FindPRContainer(st, "brotherlogic/seraphine", 141)
	if c != nil {
		t.Errorf("Expected PRContainer to be removed, but found: %v", c)
	}
}

func TestReceiveWebhook_ClosedPROfUntrackedContainer(t *testing.T) {
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			return true, nil
		},
	}

	downCalled := false
	mockDev := &mockDevcontainerClient{
		downFunc: func(ctx context.Context, in *manager_pb.DownRequest, opts ...grpc.CallOption) (*manager_pb.DownResponse, error) {
			downCalled = true
			return &manager_pb.DownResponse{}, nil
		},
	}

	pClient := pstore_client.GetTestClient()
	srv := NewWebhookServer(mockClient, mockDev, pClient)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "closed",
				Number: 141,
				Repository: &ghwebhook_pb.Repository{
					FullName: "brotherlogic/seraphine",
				},
			},
		},
	}

	_, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if downCalled {
		t.Errorf("Devcontainer.Down should not be called for untracked container")
	}
}

func TestReceiveWebhook_UpstreamFailureCases(t *testing.T) {
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			return true, nil
		},
	}

	mockDev := &mockDevcontainerClient{
		upFunc: func(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error) {
			return nil, errors.New("upstream gRPC error")
		},
	}

	pClient := pstore_client.GetTestClient()
	srv := NewWebhookServer(mockClient, mockDev, pClient)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "opened",
				Number: 141,
				Repository: &ghwebhook_pb.Repository{
					FullName: "brotherlogic/seraphine",
				},
			},
		},
	}

	resp, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("ReceiveWebhook should handle upstream errors gracefully, got: %v", err)
	}
	if resp == nil {
		t.Fatalf("Expected non-nil response")
	}
}

func TestReceiveWebhook_NonPREvent(t *testing.T) {
	collabCalled := false
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			collabCalled = true
			return true, nil
		},
	}

	srv := NewWebhookServer(mockClient, nil, nil)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "issues",
		},
	}

	resp, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("ReceiveWebhook returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("ReceiveWebhook returned nil response")
	}
	if collabCalled {
		t.Errorf("IsCollaborator should not be called for non-PR events")
	}
}

func TestReceiveWebhook_NonEligibleRepo(t *testing.T) {
	collabCalled := false
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			collabCalled = true
			return false, nil
		},
	}

	srv := NewWebhookServer(mockClient, nil, nil)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "opened",
				Number: 141,
				Title:  "Test PR",
				User: &ghwebhook_pb.User{
					Login: "testuser",
				},
				Repository: &ghwebhook_pb.Repository{
					FullName: "external/repo",
				},
			},
		},
	}

	resp, err := srv.ReceiveWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("ReceiveWebhook returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("ReceiveWebhook returned nil response")
	}
	if !collabCalled {
		t.Errorf("Expected IsCollaborator to be called, but it was not")
	}
}

