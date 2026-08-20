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

func TestIsCollaborator_True(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET request, got %s", req.Method)
			}
			if req.URL.Path != "/repos/owner/repo/collaborators/testuser" {
				t.Errorf("expected path /repos/owner/repo/collaborators/testuser, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	isCollab, err := client.IsCollaborator(context.Background(), "owner", "repo", "testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isCollab {
		t.Errorf("expected true, got false")
	}
}

func TestIsCollaborator_False(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET request, got %s", req.Method)
			}
			if req.URL.Path != "/repos/owner/repo/collaborators/noncollaborator" {
				t.Errorf("expected path /repos/owner/repo/collaborators/noncollaborator, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	isCollab, err := client.IsCollaborator(context.Background(), "owner", "repo", "noncollaborator")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isCollab {
		t.Errorf("expected false, got true")
	}
}

func TestIsCollaborator_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("internal error")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	_, err := client.IsCollaborator(context.Background(), "owner", "repo", "testuser")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestListOpenPullRequests_Success(t *testing.T) {
	mockJSON := `[
		{
			"id": 101,
			"number": 12,
			"title": "Add feature X",
			"state": "open",
			"user": {
				"login": "alice"
			},
			"head": {
				"sha": "abc1234",
				"ref": "feature-x"
			},
			"base": {
				"sha": "def5678",
				"ref": "main"
			},
			"html_url": "https://github.com/owner/repo/pull/12"
		}
	]`

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET request, got %s", req.Method)
			}
			if req.URL.Path != "/repos/owner/repo/pulls" {
				t.Errorf("expected path /repos/owner/repo/pulls, got %s", req.URL.Path)
			}
			if req.URL.Query().Get("state") != "open" {
				t.Errorf("expected query param state=open, got %s", req.URL.Query().Get("state"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockJSON)),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	prs, err := client.ListOpenPullRequests(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(prs) != 1 {
		t.Fatalf("expected 1 pull request, got %d", len(prs))
	}
	if prs[0].ID != 101 {
		t.Errorf("expected ID 101, got %d", prs[0].ID)
	}
	if prs[0].Number != 12 {
		t.Errorf("expected Number 12, got %d", prs[0].Number)
	}
	if prs[0].Title != "Add feature X" {
		t.Errorf("expected Title 'Add feature X', got %s", prs[0].Title)
	}
	if prs[0].State != "open" {
		t.Errorf("expected State 'open', got %s", prs[0].State)
	}
	if prs[0].User.Login != "alice" {
		t.Errorf("expected user alice, got %s", prs[0].User.Login)
	}
	if prs[0].Head.SHA != "abc1234" || prs[0].Head.Ref != "feature-x" {
		t.Errorf("expected head abc1234/feature-x, got %v", prs[0].Head)
	}
	if prs[0].Base.Ref != "main" {
		t.Errorf("expected base ref main, got %s", prs[0].Base.Ref)
	}
	if prs[0].HTMLURL != "https://github.com/owner/repo/pull/12" {
		t.Errorf("expected html_url https://github.com/owner/repo/pull/12, got %s", prs[0].HTMLURL)
	}
}

func TestListOpenPullRequests_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("error")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	_, err := client.ListOpenPullRequests(context.Background(), "owner", "repo")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetPullRequestDetails_Success(t *testing.T) {
	mockJSON := `{
		"id": 102,
		"number": 42,
		"title": "Bug fix",
		"state": "open",
		"user": {
			"login": "bob"
		},
		"head": {
			"sha": "deadbeef",
			"ref": "fix-bug"
		},
		"base": {
			"sha": "feedface",
			"ref": "main"
		},
		"html_url": "https://github.com/owner/repo/pull/42",
		"commits": 3,
		"comments": 5,
		"review_comments": 2
	}`

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET request, got %s", req.Method)
			}
			if req.URL.Path != "/repos/owner/repo/pulls/42" {
				t.Errorf("expected path /repos/owner/repo/pulls/42, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockJSON)),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	detail, err := client.GetPullRequestDetails(context.Background(), "owner", "repo", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if detail.Number != 42 {
		t.Errorf("expected Number 42, got %d", detail.Number)
	}
	if detail.Commits != 3 {
		t.Errorf("expected Commits 3, got %d", detail.Commits)
	}
	if detail.Comments != 5 {
		t.Errorf("expected Comments 5, got %d", detail.Comments)
	}
	if detail.ReviewComments != 2 {
		t.Errorf("expected ReviewComments 2, got %d", detail.ReviewComments)
	}
	if detail.User.Login != "bob" {
		t.Errorf("expected user bob, got %s", detail.User.Login)
	}
}

func TestGetPullRequestDetails_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("not found")),
			}, nil
		},
	}

	client := NewClient("test-token", mockClient)
	_, err := client.GetPullRequestDetails(context.Background(), "owner", "repo", 42)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetCommitCheckStatus(t *testing.T) {
	tests := []struct {
		name           string
		responseJSON   string
		statusCode     int
		expectedStatus CheckStatus
		expectErr      bool
	}{
		{
			name:           "empty check runs",
			responseJSON:   `{"total_count": 0, "check_runs": []}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusUnknown,
			expectErr:      false,
		},
		{
			name: "all successful",
			responseJSON: `{
				"total_count": 2,
				"check_runs": [
					{"id": 1, "status": "completed", "conclusion": "success"},
					{"id": 2, "status": "completed", "conclusion": "success"}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusSuccess,
			expectErr:      false,
		},
		{
			name: "in progress pending",
			responseJSON: `{
				"total_count": 2,
				"check_runs": [
					{"id": 1, "status": "completed", "conclusion": "success"},
					{"id": 2, "status": "in_progress", "conclusion": null}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusPending,
			expectErr:      false,
		},
		{
			name: "queued pending",
			responseJSON: `{
				"total_count": 1,
				"check_runs": [
					{"id": 1, "status": "queued", "conclusion": null}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusPending,
			expectErr:      false,
		},
		{
			name: "failure priority over pending",
			responseJSON: `{
				"total_count": 3,
				"check_runs": [
					{"id": 1, "status": "completed", "conclusion": "success"},
					{"id": 2, "status": "in_progress", "conclusion": null},
					{"id": 3, "status": "completed", "conclusion": "failure"}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusFailure,
			expectErr:      false,
		},
		{
			name: "timed out failure",
			responseJSON: `{
				"total_count": 1,
				"check_runs": [
					{"id": 1, "status": "completed", "conclusion": "timed_out"}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusFailure,
			expectErr:      false,
		},
		{
			name: "cancelled failure",
			responseJSON: `{
				"total_count": 1,
				"check_runs": [
					{"id": 1, "status": "completed", "conclusion": "cancelled"}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusFailure,
			expectErr:      false,
		},
		{
			name: "action required failure",
			responseJSON: `{
				"total_count": 1,
				"check_runs": [
					{"id": 1, "status": "completed", "conclusion": "action_required"}
				]
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: CheckStatusFailure,
			expectErr:      false,
		},
		{
			name:           "http error",
			responseJSON:   `internal error`,
			statusCode:     http.StatusInternalServerError,
			expectedStatus: CheckStatusUnknown,
			expectErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet {
						t.Errorf("expected GET request, got %s", req.Method)
					}
					if req.URL.Path != "/repos/owner/repo/commits/sha123/check-runs" {
						t.Errorf("expected path /repos/owner/repo/commits/sha123/check-runs, got %s", req.URL.Path)
					}
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(bytes.NewBufferString(tc.responseJSON)),
					}, nil
				},
			}

			client := NewClient("test-token", mockClient)
			status, err := client.GetCommitCheckStatus(context.Background(), "owner", "repo", "sha123")
			if (err != nil) != tc.expectErr {
				t.Fatalf("expected error %v, got %v", tc.expectErr, err)
			}
			if status != tc.expectedStatus {
				t.Errorf("expected status %s, got %s", tc.expectedStatus, status)
			}
		})
	}
}
