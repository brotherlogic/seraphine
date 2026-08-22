package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	ghwebhook_pb "github.com/brotherlogic/ghwebhook/proto/ghwebhook/v1"
	pstore_client "github.com/brotherlogic/pstore/client"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/dashboard"
	"github.com/brotherlogic/seraphine/internal/github"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	err = runSync(ctx, pClient, ghClient, nil)
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

type mockRegistrationClient struct {
	ghwebhook_pb.RegistrationServiceClient
	registered []string
}

func (m *mockRegistrationClient) Register(ctx context.Context, in *ghwebhook_pb.RegistrationRequest, opts ...grpc.CallOption) (*ghwebhook_pb.RegistrationResponse, error) {
	m.registered = append(m.registered, in.GetRepoFullName())
	return &ghwebhook_pb.RegistrationResponse{Success: true}, nil
}

func TestWebhookServerRegistration(t *testing.T) {
	grpcServer := grpc.NewServer()
	webhookServer := NewWebhookServer(nil, nil, nil)
	ghwebhook_pb.RegisterWebhookHandlerServer(grpcServer, webhookServer)

	info := grpcServer.GetServiceInfo()
	if _, ok := info["ghwebhook.v1.WebhookHandler"]; !ok {
		t.Errorf("Expected ghwebhook.v1.WebhookHandler to be registered, got services: %v", info)
	}
}

func TestSyncWorkerWebhookRegistration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pClient := pstore_client.GetTestClient()
	initialState := &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/repo1", "brotherlogic/repo2"},
	}
	if err := config.WriteServerState(ctx, pClient, initialState); err != nil {
		t.Fatalf("Failed to write initial state: %v", err)
	}

	mockHTTP := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == "GET" && req.URL.Path == "/user/repository_invitations" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
				}, nil
			}
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/rulesets") {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
			}, nil
		},
	}

	ghClient := github.NewClient("fake-token", mockHTTP)
	mockReg := &mockRegistrationClient{}

	err := runSync(ctx, pClient, ghClient, mockReg)
	if err != nil {
		t.Fatalf("runSync failed: %v", err)
	}

	if len(mockReg.registered) != 2 {
		t.Fatalf("Expected 2 registered repos, got %d: %v", len(mockReg.registered), mockReg.registered)
	}
}

func TestDevcontainerManagerAddressAndWiring(t *testing.T) {
	// Test default address
	t.Setenv("DEVCONTAINER_MANAGER_ADDRESS", "")
	addr := getDevcontainerAddress()
	expectedDefault := "devcontainer-manager.devcontainer-manager.svc.cluster.local:8080"
	if addr != expectedDefault {
		t.Errorf("Expected default address %s, got %s", expectedDefault, addr)
	}

	// Test custom env address
	customAddr := "custom-devcontainer:9090"
	t.Setenv("DEVCONTAINER_MANAGER_ADDRESS", customAddr)
	if getDevcontainerAddress() != customAddr {
		t.Errorf("Expected custom address %s, got %s", customAddr, getDevcontainerAddress())
	}

	// Test explicit argument overrides env address
	explicitAddr := "explicit-devcontainer:9999"
	if getDevcontainerAddress(explicitAddr) != explicitAddr {
		t.Errorf("Expected explicit address %s, got %s", explicitAddr, getDevcontainerAddress(explicitAddr))
	}

	// Test empty explicit argument falls back to env address
	if getDevcontainerAddress("") != customAddr {
		t.Errorf("Expected fallback to env address %s, got %s", customAddr, getDevcontainerAddress(""))
	}

	// Test wiring into WebhookServer
	mockDevClient := &mockDevcontainerClient{}
	ws := NewWebhookServer(nil, mockDevClient, nil)
	if ws.devcontainerClient != mockDevClient {
		t.Errorf("Expected devcontainerClient to be wired into WebhookServer")
	}
}

func TestSeraphineServer_DashboardServiceExposure(t *testing.T) {
	pClient := pstore_client.GetTestClient()
	dashService := dashboard.NewService(nil, nil, pClient)

	srv := NewSeraphineServer(dashService)
	if srv == nil {
		t.Fatalf("Expected non-nil SeraphineServer")
	}

	if srv.GetDashboardService() != dashService {
		t.Errorf("Expected GetDashboardService() to return dashService instance")
	}

	if srv.DashboardService != dashService {
		t.Errorf("Expected DashboardService field to match dashService instance")
	}
}

type mockDashboardService struct {
	syncCount int
	runCalled bool
}

func (m *mockDashboardService) GetDashboardState(ctx context.Context) (*dashboard.DashboardState, error) {
	return &dashboard.DashboardState{}, nil
}

func (m *mockDashboardService) RunWorker(ctx context.Context, interval time.Duration) {
	m.runCalled = true
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	m.syncCount++
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.syncCount++
		}
	}
}

func TestDashboardWorkerLifecycle(t *testing.T) {
	mockDash := &mockDashboardService{}
	ctx, cancel := context.WithCancel(context.Background())

	workerDone := make(chan struct{})
	go func() {
		mockDash.RunWorker(ctx, 10*time.Millisecond)
		close(workerDone)
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-workerDone:
		// Succeeded
	case <-time.After(1 * time.Second):
		t.Fatalf("Worker did not exit within timeout after context cancellation")
	}

	if !mockDash.runCalled {
		t.Errorf("Expected RunWorker to have been called")
	}
	if mockDash.syncCount < 2 {
		t.Errorf("Expected at least 2 sync executions, got %d", mockDash.syncCount)
	}
}

func getFreePort(t *testing.T) int {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port
}

func TestGetHTTPPort(t *testing.T) {
	// Test default
	t.Setenv("HTTP_PORT", "")
	if port := getHTTPPort(); port != ":8080" {
		t.Errorf("Expected default :8080, got %s", port)
	}

	// Test environment variable without colon
	t.Setenv("HTTP_PORT", "8888")
	if port := getHTTPPort(); port != ":8888" {
		t.Errorf("Expected :8888 from env, got %s", port)
	}

	// Test environment variable with colon
	t.Setenv("HTTP_PORT", ":9999")
	if port := getHTTPPort(); port != ":9999" {
		t.Errorf("Expected :9999 from env, got %s", port)
	}

	// Test explicit argument overrides env
	if port := getHTTPPort(":7777"); port != ":7777" {
		t.Errorf("Expected explicit :7777 to override env, got %s", port)
	}

	// Test explicit argument without colon
	if port := getHTTPPort("7777"); port != ":7777" {
		t.Errorf("Expected :7777 with added colon, got %s", port)
	}
}

func TestRunWithContext_ConcurrentLifecycleAndGracefulShutdown(t *testing.T) {
	grpcPort := getFreePort(t)
	httpPort := getFreePort(t)

	grpcAddr := fmt.Sprintf("127.0.0.1:%d", grpcPort)
	httpAddr := fmt.Sprintf("127.0.0.1:%d", httpPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- RunWithContext(ctx, grpcAddr, httpAddr)
	}()

	// 1. Verify HTTP server is serving requests
	httpClient := &http.Client{Timeout: 500 * time.Millisecond}
	httpHealthURL := fmt.Sprintf("http://%s/healthz", httpAddr)

	var lastHTTPErr error
	httpStarted := false
	for i := 0; i < 50; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err := httpClient.Get(httpHealthURL)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK && string(body) == "OK\n" {
				httpStarted = true
				break
			}
		}
		lastHTTPErr = err
	}
	if !httpStarted {
		t.Fatalf("HTTP server failed to respond on %s: %v", httpHealthURL, lastHTTPErr)
	}

	// 2. Verify gRPC server is accepting RPCs
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial gRPC server at %s: %v", grpcAddr, err)
	}
	defer grpcConn.Close()

	grpcClient := pb.NewSeraphineServiceClient(grpcConn)
	_, grpcErr := grpcClient.GetProjectState(context.Background(), &pb.GetProjectStateRequest{})
	// Expected to return codes.Unimplemented because method is unimplemented, but proving gRPC server is connected and responding
	if status.Code(grpcErr) != codes.Unimplemented {
		t.Fatalf("Expected Unimplemented code from gRPC server, got error: %v", grpcErr)
	}

	// 3. Graceful shutdown
	cancel()

	select {
	case err := <-serverErrChan:
		if err != nil {
			t.Errorf("RunWithContext returned unexpected error on shutdown: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("RunWithContext did not terminate within timeout after context cancellation")
	}

	// 4. Verify HTTP server is shut down
	_, err = httpClient.Get(httpHealthURL)
	if err == nil {
		t.Errorf("Expected HTTP connection to fail after shutdown, but got successful response")
	}
}

func TestRunWithContext_InvalidPort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := RunWithContext(ctx, "999.999.999.999:9009", "127.0.0.1:8080")
	if err == nil {
		t.Errorf("Expected error for invalid gRPC address, got nil")
	}
}



