package web

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/brotherlogic/seraphine/internal/dashboard"
	seraphineweb "github.com/brotherlogic/seraphine/web"
)

// RunHTTPServer runs the HTTP web server on the specified address with graceful shutdown support using embedded static assets.
func RunHTTPServer(ctx context.Context, addr string, dashboardSvc dashboard.Service) error {
	return RunHTTPServerWithFS(ctx, addr, dashboardSvc, seraphineweb.GetStaticFS())
}

// RunHTTPServerWithFS runs the HTTP web server on the specified address with custom static assets and graceful shutdown support.
func RunHTTPServerWithFS(ctx context.Context, addr string, dashboardSvc dashboard.Service, staticFS fs.FS) error {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	return ServeHTTPWithFS(ctx, lis, dashboardSvc, staticFS)
}

// ServeHTTP serves the HTTP web handler on an existing listener with graceful shutdown using embedded static assets.
func ServeHTTP(ctx context.Context, lis net.Listener, dashboardSvc dashboard.Service) error {
	return ServeHTTPWithFS(ctx, lis, dashboardSvc, seraphineweb.GetStaticFS())
}

// ServeHTTPWithFS serves the HTTP web handler on an existing listener with custom static assets and graceful shutdown.
func ServeHTTPWithFS(ctx context.Context, lis net.Listener, dashboardSvc dashboard.Service, staticFS fs.FS) error {
	srv := &http.Server{
		Handler: NewHandler(dashboardSvc, staticFS),
	}

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

