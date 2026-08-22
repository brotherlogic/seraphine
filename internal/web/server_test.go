package web

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/brotherlogic/seraphine/internal/dashboard"
)

type mockServerDashboardService struct{}

func (m *mockServerDashboardService) GetDashboardState(ctx context.Context) (*dashboard.DashboardState, error) {
	return &dashboard.DashboardState{
		PullRequests: []dashboard.PRSummary{},
	}, nil
}

func (m *mockServerDashboardService) RunWorker(ctx context.Context, interval time.Duration) {}

func getFreePort(t *testing.T) int {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port
}

func TestRunHTTPServer_LifecycleAndShutdown(t *testing.T) {
	port := getFreePort(t)
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dashSvc := &mockServerDashboardService{}
	serverErrChan := make(chan error, 1)

	go func() {
		serverErrChan <- RunHTTPServer(ctx, addr, dashSvc)
	}()

	// Wait for server to start accepting connections
	client := &http.Client{Timeout: 500 * time.Millisecond}
	reqURL := fmt.Sprintf("http://%s/healthz", addr)

	var lastErr error
	started := false
	for i := 0; i < 50; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err := client.Get(reqURL)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK && string(body) == "OK\n" {
				started = true
				break
			}
		}
		lastErr = err
	}

	if !started {
		t.Fatalf("HTTP server failed to start at %s: %v", reqURL, lastErr)
	}

	// Test API endpoint
	apiURL := fmt.Sprintf("http://%s/api/dashboard", addr)
	resp, err := client.Get(apiURL)
	if err != nil {
		t.Fatalf("Failed to GET /api/dashboard: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Trigger graceful shutdown
	client.CloseIdleConnections()
	cancel()

	select {
	case err := <-serverErrChan:
		if err != nil {
			t.Errorf("RunHTTPServer returned unexpected error on shutdown: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("RunHTTPServer did not terminate within timeout after context cancellation")
	}

	// Verify server is no longer accepting requests
	_, err = client.Get(reqURL)
	if err == nil {
		t.Errorf("Expected connection to fail after shutdown, but got successful response")
	}
}

func TestRunHTTPServer_StaticAssetServing(t *testing.T) {
	port := getFreePort(t)
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dashSvc := &mockServerDashboardService{}
	serverErrChan := make(chan error, 1)

	go func() {
		serverErrChan <- RunHTTPServer(ctx, addr, dashSvc)
	}()

	// Wait for server to start
	client := &http.Client{Timeout: 500 * time.Millisecond}
	reqURL := fmt.Sprintf("http://%s/healthz", addr)

	started := false
	for i := 0; i < 50; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err := client.Get(reqURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				started = true
				break
			}
		}
	}
	if !started {
		t.Fatalf("HTTP server failed to start at %s", reqURL)
	}

	// Test GET /
	resp, err := client.Get(fmt.Sprintf("http://%s/", addr))
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for GET /, got %d", resp.StatusCode)
	}
	if !strings.Contains(string(body), "html") && !strings.Contains(string(body), "DOCTYPE") && !strings.Contains(string(body), "Seraphine") {
		t.Errorf("Expected HTML response from GET /, got: %s", string(body))
	}

	// Test GET /index.html
	resp, err = client.Get(fmt.Sprintf("http://%s/index.html", addr))
	if err != nil {
		t.Fatalf("GET /index.html failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for GET /index.html, got %d", resp.StatusCode)
	}

	// Test SPA client-side fallback route (e.g. GET /dashboard/settings)
	resp, err = client.Get(fmt.Sprintf("http://%s/dashboard/settings", addr))
	if err != nil {
		t.Fatalf("GET /dashboard/settings failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for SPA route, got %d", resp.StatusCode)
	}

	client.CloseIdleConnections()
	cancel()
	<-serverErrChan
}

func TestRunHTTPServer_InvalidAddress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use an invalid address
	err := RunHTTPServer(ctx, "999.999.999.999:80", nil)
	if err == nil {
		t.Errorf("Expected error for invalid address, got nil")
	}
}


