package workflow

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFetchWorkflows(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "workflow-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = FetchWorkflows(context.Background(), tempDir)
	if err != nil {
		t.Fatalf("FetchWorkflows returned error: %v", err)
	}

	// Verify ISSUES.md was fetched
	if _, err := os.Stat(filepath.Join(tempDir, "ISSUES.md")); os.IsNotExist(err) {
		t.Errorf("ISSUES.md was not fetched")
	}

	// Verify at least one workflow was fetched
	workflowsDir := filepath.Join(tempDir, ".agent", "workflows")
	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		t.Errorf("Workflows directory was not created")
	}
}
