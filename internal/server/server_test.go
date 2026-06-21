package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	pstore_client "github.com/brotherlogic/pstore/client"
	pb "github.com/brotherlogic/seraphine/proto"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/github"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.doFunc != nil {
		return m.doFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
	}, nil
}

func TestSyncWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pClient := pstore_client.GetTestClient()

	// Initial state: no enrolled repos
	initialState := &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/some-repo"},
	}
	err := config.WriteServerState(ctx, pClient, initialState)
	if err != nil {
		t.Fatalf("Failed to write initial state: %v", err)
	}

	invitations := []*github.RepositoryInvitation{
		{
			ID: 12345,
			Repository: github.Repository{
				Name:     "new-repo",
				FullName: "brotherlogic/new-repo",
				Owner: github.Owner{
					Login: "brotherlogic",
				},
			},
		},
	}
	invitationsJSON, _ := json.Marshal(invitations)

	issueResp := &github.IssueResponse{Number: 42}
	issueJSON, _ := json.Marshal(issueResp)

	mockHTTP := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == "GET" && req.URL.Path == "/user/repository_invitations" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(invitationsJSON)),
				}, nil
			}
			if req.Method == "PATCH" && req.URL.Path == "/user/repository_invitations/12345" {
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			}
			if req.Method == "POST" && req.URL.Path == "/repos/brotherlogic/new-repo/rulesets" {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			}
			if req.Method == "POST" && req.URL.Path == "/repos/brotherlogic/devcontainer-manager/issues" {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader(issueJSON)),
				}, nil
			}
			if req.Method == "POST" && req.URL.Path == "/repos/brotherlogic/some-repo/rulesets" {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			}

			t.Logf("Unexpected request: %s %s", req.Method, req.URL.Path)
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
			}, nil
		},
	}

	ghClient := github.NewClient("fake-token", mockHTTP)

	// Run the sync process once
	err = runSync(ctx, pClient, ghClient)
	if err != nil {
		t.Fatalf("runSync failed: %v", err)
	}

	// Verify state was updated
	state, err := config.ReadServerState(ctx, pClient)
	if err != nil {
		t.Fatalf("Failed to read server state: %v", err)
	}

	found := false
	for _, repo := range state.EnrolledRepositories {
		if repo == "brotherlogic/new-repo" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected brotherlogic/new-repo to be enrolled, got: %v", state.EnrolledRepositories)
	}
}
