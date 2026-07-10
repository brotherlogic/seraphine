package workflow

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

// FetchWorkflows fetches ISSUES.md and workflow files from the raw github repo
func FetchWorkflows(ctx context.Context, destDir string) error {
	baseURL := "https://raw.githubusercontent.com/brotherlogic/seraphine/main"
	
	// Fetch ISSUES.md
	issuesURL := fmt.Sprintf("%s/ISSUES.md", baseURL)
	issuesContent, err := fetchURL(ctx, issuesURL)
	if err != nil {
		return fmt.Errorf("failed to fetch ISSUES.md: %w", err)
	}

	// Write ISSUES.md
	issuesPath := filepath.Join(destDir, "ISSUES.md")
	err = os.WriteFile(issuesPath, issuesContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write ISSUES.md: %w", err)
	}

	// Parse ISSUES.md for workflow files
	re := regexp.MustCompile(`\.agent/workflows/([a-zA-Z0-9_-]+\.md)`)
	matches := re.FindAllStringSubmatch(string(issuesContent), -1)

	// Create workflows directory
	workflowsDir := filepath.Join(destDir, ".agent", "workflows")
	err = os.MkdirAll(workflowsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	// Fetch each workflow file
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		filename := match[1]
		workflowURL := fmt.Sprintf("%s/.agent/workflows/%s", baseURL, filename)
		workflowContent, err := fetchURL(ctx, workflowURL)
		if err != nil {
			return fmt.Errorf("failed to fetch workflow %s: %w", filename, err)
		}

		workflowPath := filepath.Join(workflowsDir, filename)
		err = os.WriteFile(workflowPath, workflowContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write workflow %s: %w", filename, err)
		}
	}

	return nil
}

func fetchURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
