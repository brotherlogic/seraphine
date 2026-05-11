package workflow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brotherlogic/seraphine/internal/client"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/reconciler"
)

func RunSync(ctx context.Context, serverAddr string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	projectName := filepath.Base(dir)

	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	resp, err := client.GetProjectState(ctx, serverAddr, projectName, cfg.Version)
	if err != nil {
		return fmt.Errorf("error getting project state from server: %w", err)
	}

	if resp.Version == cfg.Version {
		fmt.Println("Project is already up to date.")
		return nil
	}

	err = reconciler.Reconcile(resp.Files)
	if err != nil {
		return fmt.Errorf("error reconciling project state: %w", err)
	}

	cfg.Version = resp.Version
	err = config.WriteConfig(cfg)
	if err != nil {
		return fmt.Errorf("error writing config: %w", err)
	}

	fmt.Printf("Successfully synced project to version %s\n", cfg.Version)
	return nil
}
