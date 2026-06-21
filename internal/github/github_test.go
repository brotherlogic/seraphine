package github

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestListRepositoryInvitations_Success(t *testing.T) {
	mockJSON := `[
		{
			"id": 12345,
			"repository": {
				"id": 67890,
				"name": "test-repo",
				"full_name": "owner/test-repo",
				"owner": {
					"login": "owner"
				}
			}
		}
	]`

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET request, got %s", req.Method)
			}
			if req.URL.Path != "/user/repository_invitations" {
				t.Errorf("expected path /user/repository_invitations, got %s", req.URL.Path)
			}
			if req.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("expected Authorization Bearer test-token, got %s", req.Header.Get("Authorization"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockJSON)),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	invitations, err := client.ListRepositoryInvitations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(invitations) != 1 {
		t.Fatalf("expected 1 invitation, got %d", len(invitations))
	}
	if invitations[0].ID != 12345 {
		t.Errorf("expected ID 12345, got %d", invitations[0].ID)
	}
	if invitations[0].Repository.Name != "test-repo" {
		t.Errorf("expected repo name test-repo, got %s", invitations[0].Repository.Name)
	}
}

func TestAcceptRepositoryInvitation_Success(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH request, got %s", req.Method)
			}
			if req.URL.Path != "/user/repository_invitations/12345" {
				t.Errorf("expected path /user/repository_invitations/12345, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	err := client.AcceptRepositoryInvitation(context.Background(), 12345)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateRuleset_Success(t *testing.T) {
	ruleset := &RulesetRequest{
		Name:        "Seraphine Default Branch Protection",
		Target:      "branch",
		Enforcement: "active",
		Conditions: Conditions{
			RefName: RefName{
				Include: []string{"~DEFAULT_BRANCH"},
				Exclude: []string{},
			},
		},
		Rules: []Rule{
			{
				Type: "pull_request",
				Parameters: &RuleParameters{
					RequiredApprovingReviewCount: 1,
				},
			},
		},
	}

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", req.Method)
			}
			if req.URL.Path != "/repos/owner/repo/rulesets" {
				t.Errorf("expected path /repos/owner/repo/rulesets, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewBufferString("{}")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	err := client.CreateRuleset(context.Background(), "owner", "repo", ruleset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateIssue_Success(t *testing.T) {
	mockResponseJSON := `{"number": 42}`

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", req.Method)
			}
			if req.URL.Path != "/repos/brotherlogic/devcontainer-manager/issues" {
				t.Errorf("expected path /repos/brotherlogic/devcontainer-manager/issues, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseJSON)),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	resp, err := client.CreateIssue(context.Background(), "brotherlogic", "devcontainer-manager", "Build out an API", "Add a grpc api", []string{"seraphine-needs-requirements"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Number != 42 {
		t.Errorf("expected issue number 42, got %d", resp.Number)
	}
}
