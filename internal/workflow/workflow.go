package workflow

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/brotherlogic/seraphine/internal/client"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/reconciler"
	pb "github.com/brotherlogic/seraphine/proto"
)

func RunInit(ctx context.Context, serverAddr string) error {
	repoURL, err := getGitRemoteURL()
	if err != nil {
		return fmt.Errorf("error getting git remote URL: %w", err)
	}

	projectName, err := extractProjectName(repoURL)
	if err != nil {
		return fmt.Errorf("error extracting project name: %w", err)
	}

	err = validateRuleset(projectName)
	if err != nil {
		fmt.Printf("Ruleset validation failed: %v. Filing issue...\n", err)
		err = fileRulesetIssue(projectName)
		if err != nil {
			fmt.Printf("Failed to file ruleset issue: %v\n", err)
		}
	}

	fmt.Printf("Registering project %s (%s)...\n", projectName, repoURL)
	resp, err := client.RegisterProject(ctx, serverAddr, projectName, repoURL)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			fmt.Println("Project is already registered.")
		} else {
			return fmt.Errorf("error registering project: %w", err)
		}
	}

	cfg := &pb.ProjectConfig{
		ProjectName: projectName,
		Version:     "0",
	}
	if resp != nil && resp.Version != "" {
		cfg.Version = resp.Version
	}

	err = config.WriteConfig(cfg)
	if err != nil {
		return fmt.Errorf("error writing config: %w", err)
	}

	fmt.Printf("Successfully initialized project at version %s\n", cfg.Version)

	// Issue #11 will handle the automated upgrade, but for now we just finish init
	return nil
}

func getGitRemoteURL() (string, error) {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func extractProjectName(repoURL string) (string, error) {
	// Handle https://github.com/username/repo.git or git@github.com:username/repo.git
	trimmed := strings.TrimSuffix(repoURL, ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		// Try SSH format git@github.com:username/repo
		parts = strings.Split(trimmed, ":")
	}

	if len(parts) < 2 {
		return "", fmt.Errorf("could not parse repo URL: %s", repoURL)
	}

	repo := parts[len(parts)-1]
	userPart := parts[len(parts)-2]
	
	// Handle userPart if it contains github.com:username
	if strings.Contains(userPart, ":") {
		uParts := strings.Split(userPart, ":")
		userPart = uParts[len(uParts)-1]
	}

	return fmt.Sprintf("%s/%s", userPart, repo), nil
}

func RunSync(ctx context.Context, serverAddr string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	projectName := cfg.ProjectName
	if projectName == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %w", err)
		}
		projectName = filepath.Base(dir)
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

	repoURL, err := getGitRemoteURL()
	if err != nil {
		return fmt.Errorf("error getting git remote URL: %w", err)
	}
	ownerRepo, err := extractProjectName(repoURL)
	if err != nil {
		return fmt.Errorf("error extracting project name: %w", err)
	}

	err = reconciler.ReconcileGithubSettings(ctx, ownerRepo, resp.GithubSettings)
	if err != nil {
		return fmt.Errorf("error reconciling github settings: %w", err)
	}

	cfg.Version = resp.Version
	err = config.WriteConfig(cfg)
	if err != nil {
		return fmt.Errorf("error writing config: %w", err)
	}

	err = gitCommitAndPush(cfg.Version)
	if err != nil {
		return fmt.Errorf("error pushing upgrade branch: %w", err)
	}

	fmt.Printf("Successfully synced project to version %s\n", cfg.Version)
	return nil
}

func gitCommitAndPush(version string) error {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		fmt.Println("No file changes to commit.")
		return nil
	}

	branchName := fmt.Sprintf("seraphine/upgrade-%s", version)

	cmd := exec.Command("git", "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("git", "checkout", branchName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout branch: %w", err)
		}
	}

	cmd = exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	commitMsg := fmt.Sprintf("chore: automated seraphine upgrade to %s", version)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	cmd = exec.Command("git", "push", "-u", "origin", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	return nil
}

func validateRuleset(ownerRepo string) error {
	cmd := exec.Command("gh", "api", fmt.Sprintf("/repos/%s/rulesets", ownerRepo))
	return cmd.Run()
}

func fileRulesetIssue(ownerRepo string) error {
	cmd := exec.Command("gh", "issue", "create", "--repo", ownerRepo, "--title", "Add branch protection ruleset", "--body", "Please add a branch protection ruleset to ensure repository security and compliance.")
	return cmd.Run()
}
