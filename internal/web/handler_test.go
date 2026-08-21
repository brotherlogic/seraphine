package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/brotherlogic/seraphine/internal/dashboard"
	"github.com/brotherlogic/seraphine/internal/github"
)

type mockDashboardService struct {
	state *dashboard.DashboardState
	err   error
}

func (m *mockDashboardService) GetDashboardState(ctx context.Context) (*dashboard.DashboardState, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.state, nil
}

func (m *mockDashboardService) RunWorker(ctx context.Context, interval time.Duration) {}

func createTestStaticFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data:    []byte("<!DOCTYPE html><html><body><h1>Seraphine Dashboard</h1></body></html>"),
			Mode:    0644,
			ModTime: time.Now(),
		},
		"assets/app.js": &fstest.MapFile{
			Data:    []byte("console.log('seraphine');"),
			Mode:    0644,
			ModTime: time.Now(),
		},
		"assets/style.css": &fstest.MapFile{
			Data:    []byte("body { background: #000; }"),
			Mode:    0644,
			ModTime: time.Now(),
		},
		"docs/index.html": &fstest.MapFile{
			Data:    []byte("<h1>Docs</h1>"),
			Mode:    0644,
			ModTime: time.Now(),
		},
	}
}

func TestGetDashboard_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedState := &dashboard.DashboardState{
		PullRequests: []dashboard.PRSummary{
			{
				Repo:            "brotherlogic/seraphine",
				PRNumber:        178,
				Title:           "Test PR",
				Author:          "alice",
				CommitCount:     3,
				CommentCount:    2,
				CheckStatus:     github.CheckStatusSuccess,
				HasDevcontainer: true,
				ContainerID:     "container-178",
				ContainerState:  dashboard.ContainerStateReady,
			},
		},
		Freshness: dashboard.FreshnessMetadata{
			LastSuccessfulSync: now,
			LastAttemptedSync:  now,
			IsStale:            false,
		},
	}

	mockSvc := &mockDashboardService{state: expectedState}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var actualState dashboard.DashboardState
	if err := json.NewDecoder(resp.Body).Decode(&actualState); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if len(actualState.PullRequests) != 1 {
		t.Fatalf("expected 1 PR, got %d", len(actualState.PullRequests))
	}
	pr := actualState.PullRequests[0]
	if pr.Repo != "brotherlogic/seraphine" || pr.PRNumber != 178 || pr.Title != "Test PR" || pr.Author != "alice" {
		t.Errorf("unexpected PR data: %+v", pr)
	}
	if pr.CheckStatus != github.CheckStatusSuccess || pr.ContainerState != dashboard.ContainerStateReady {
		t.Errorf("unexpected PR status: %+v", pr)
	}
}

func TestGetDashboard_Error(t *testing.T) {
	mockSvc := &mockDashboardService{err: errors.New("upstream service failure")}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

func TestGetDashboard_NilService(t *testing.T) {
	staticFS := createTestStaticFS()
	handler := NewHandler(nil, staticFS)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error when service is nil, got %d", resp.StatusCode)
	}
}

func TestGetDashboard_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodPost, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 Method Not Allowed for POST /api/dashboard, got %d", resp.StatusCode)
	}
}

func TestHealthz_Success(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if string(body) != "OK\n" && string(body) != "OK" && string(body) != "ok" && string(body) != "{\"status\":\"ok\"}" {
		t.Logf("got body %q", string(body))
	}
}

func TestHealthz_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 Method Not Allowed for POST /healthz, got %d", resp.StatusCode)
	}
}

func TestStaticFiles_ServeIndexAndAssets(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	// Test root path "/" -> serves index.html
	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET / expected 200 OK, got %d", resp.StatusCode)
		}
		if string(body) != "<!DOCTYPE html><html><body><h1>Seraphine Dashboard</h1></body></html>" {
			t.Errorf("GET / unexpected body: %s", string(body))
		}
	}

	// Test direct asset path "/assets/app.js"
	{
		req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /assets/app.js expected 200 OK, got %d", resp.StatusCode)
		}
		if string(body) != "console.log('seraphine');" {
			t.Errorf("GET /assets/app.js unexpected body: %s", string(body))
		}
	}

	// Test direct asset path "/assets/style.css"
	{
		req := httptest.NewRequest(http.MethodGet, "/assets/style.css", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /assets/style.css expected 200 OK, got %d", resp.StatusCode)
		}
		if string(body) != "body { background: #000; }" {
			t.Errorf("GET /assets/style.css unexpected body: %s", string(body))
		}
	}
}

func TestStaticFiles_SPAFallback(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	// Non-existent client-side route should fallback to index.html
	testRoutes := []string{
		"/dashboard",
		"/pull-requests/178",
		"/settings/filters",
	}

	for _, route := range testRoutes {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s fallback expected 200 OK, got %d", route, resp.StatusCode)
		}
		if string(body) != "<!DOCTYPE html><html><body><h1>Seraphine Dashboard</h1></body></html>" {
			t.Errorf("GET %s fallback unexpected body: %s", route, string(body))
		}
	}
}

func TestStaticFiles_DirectoryIndex(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodGet, "/docs/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /docs/ expected 200 OK, got %d", resp.StatusCode)
	}
	if string(body) != "<h1>Docs</h1>" {
		t.Errorf("GET /docs/ unexpected body: %s", string(body))
	}
}

func TestStaticFiles_UnmatchedAPIRouteNotFound(t *testing.T) {
	mockSvc := &mockDashboardService{}
	staticFS := createTestStaticFS()
	handler := NewHandler(mockSvc, staticFS)

	req := httptest.NewRequest(http.MethodGet, "/api/unknown-endpoint", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found for unmatched /api/ route, got %d", resp.StatusCode)
	}
}

func TestStaticFiles_NilFS(t *testing.T) {
	mockSvc := &mockDashboardService{}
	handler := NewHandler(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found when staticFS is nil, got %d", resp.StatusCode)
	}
}

func TestStaticFiles_MissingIndexFallback(t *testing.T) {
	mockSvc := &mockDashboardService{}
	emptyFS := fstest.MapFS{
		"assets/app.js": &fstest.MapFile{
			Data: []byte("console.log('no index');"),
		},
	}
	handler := NewHandler(mockSvc, emptyFS)

	req := httptest.NewRequest(http.MethodGet, "/non-existent-page", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found when index.html is missing in fallback, got %d", resp.StatusCode)
	}
}
