package server

import (
	"context"
	"testing"

	ghwebhook_pb "github.com/brotherlogic/ghwebhook/proto/ghwebhook/v1"
	"github.com/brotherlogic/seraphine/internal/github"
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

func TestReceiveWebhook_OpenedPROfEligibleRepo(t *testing.T) {
	collabCalled := false
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			collabCalled = true
			if owner == "brotherlogic" && repo == "seraphine" && user == "brotherlogicreviewer@gmail.com" {
				return true, nil
			}
			return false, nil
		},
	}

	srv := NewWebhookServer(mockClient)

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
}

func TestReceiveWebhook_NonOpenedPRAction(t *testing.T) {
	collabCalled := false
	mockClient := &mockGHClient{
		isCollaboratorFunc: func(ctx context.Context, owner, repo, user string) (bool, error) {
			collabCalled = true
			return true, nil
		},
	}

	srv := NewWebhookServer(mockClient)

	req := &ghwebhook_pb.WebhookEvent{
		Header: &ghwebhook_pb.EventHeader{
			EventType: "pull_request",
		},
		Payload: &ghwebhook_pb.WebhookEvent_PullRequest{
			PullRequest: &ghwebhook_pb.PullRequestEvent{
				Action: "closed",
				Number: 141,
				Title:  "Test PR",
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
	if collabCalled {
		t.Errorf("IsCollaborator should not be called for non-opened PR actions")
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

	srv := NewWebhookServer(mockClient)

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

	srv := NewWebhookServer(mockClient)

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
