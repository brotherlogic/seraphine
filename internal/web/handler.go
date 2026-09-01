package web

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/brotherlogic/seraphine/internal/dashboard"
)

// NewHandler creates and configures the HTTP handler for Seraphine's web endpoints and static assets.
func NewHandler(dashboardSvc dashboard.Service, staticFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	// GET /api/dashboard
	mux.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if dashboardSvc == nil {
			http.Error(w, "dashboard service not configured", http.StatusInternalServerError)
			return
		}

		state, err := dashboardSvc.GetDashboardState(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// GET /healthz
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	})

	// Static file server with SPA fallback routing
	fileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if staticFS == nil {
			http.NotFound(w, r)
			return
		}

		// Don't fallback to index.html for unmatched API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		reqPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if reqPath == "" || reqPath == "." {
			reqPath = "index.html"
		}

		// Check if requested file or directory index exists
		f, err := staticFS.Open(reqPath)
		if err == nil {
			stat, statErr := f.Stat()
			_ = f.Close()
			if statErr == nil {
				if !stat.IsDir() {
					http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
					return
				}
				idxPath := path.Join(reqPath, "index.html")
				if idxF, idxErr := staticFS.Open(idxPath); idxErr == nil {
					_ = idxF.Close()
					http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
					return
				}
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			// Other FS errors
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fallback to index.html for client-side routing
		if idxF, err := staticFS.Open("index.html"); err == nil {
			_ = idxF.Close()
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/"
			http.FileServer(http.FS(staticFS)).ServeHTTP(w, r2)
			return
		}

		http.NotFound(w, r)
	})

	mux.Handle("/", fileServer)

	return mux
}
