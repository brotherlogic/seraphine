package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/brotherlogic/seraphine/internal/dashboard"
)

// RunHTTPServer runs the HTTP web server on the specified address with graceful shutdown support.
func RunHTTPServer(ctx context.Context, addr string, dashboardSvc dashboard.Service) error {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	return ServeHTTP(ctx, lis, dashboardSvc)
}

// ServeHTTP serves the HTTP web handler on an existing listener with graceful shutdown.
func ServeHTTP(ctx context.Context, lis net.Listener, dashboardSvc dashboard.Service) error {
	srv := &http.Server{
		Handler: NewHandler(dashboardSvc, nil),
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
