package reconciler

import (
	"context"
	"fmt"
	"strings"
	"testing"

	pb "github.com/brotherlogic/seraphine/proto"
)

func TestReconcileGithubSettings_Empty(t *testing.T) {
	err := ReconcileGithubSettings(context.Background(), "owner/repo", nil)
	if err != nil {
		t.Fatalf("expected no error for empty settings, got %v", err)
	}
}

func TestReconcileGithubSettings_NoChanges(t *testing.T) {
	// Mock runCommand to return existing settings matching the desired settings
	oldRunCommand := runCommand
	defer func() { runCommand = oldRunCommand }()

	runCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		if args[0] == "api" && args[1] == "repos/owner/repo" {
			return []byte(`{"delete_branch_on_merge": true, "allow_squash_merge": false}`), nil
		}
		if args[0] == "api" && args[1] == "repos/owner/repo/actions/permissions/workflow" {
			return []byte(`{"default_workflow_permissions": "read"}`), nil
		}
		return nil, fmt.Errorf("unexpected command call: %v", args)
	}

	settings := []*pb.GithubSetting{
		{Key: "delete_branch_on_merge", Value: "true"},
		{Key: "allow_squash_merge", Value: "false"},
		{Key: "default_workflow_permissions", Value: "read"},
	}

	err := ReconcileGithubSettings(context.Background(), "owner/repo", settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReconcileGithubSettings_WithChanges(t *testing.T) {
	oldRunCommand := runCommand
	defer func() { runCommand = oldRunCommand }()

	var patchCalled bool
	var putCalled bool

	runCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		// Mock GET calls
		if args[0] == "api" && args[1] == "repos/owner/repo" {
			return []byte(`{"delete_branch_on_merge": false}`), nil
		}
		if args[0] == "api" && args[1] == "repos/owner/repo/actions/permissions/workflow" {
			return []byte(`{"default_workflow_permissions": "write"}`), nil
		}

		// Mock PATCH/PUT calls
		if args[0] == "api" && args[1] == "-X" && args[2] == "PATCH" && args[3] == "repos/owner/repo" && args[4] == "-F" && args[5] == "delete_branch_on_merge=true" {
			patchCalled = true
			return []byte(`{}`), nil
		}
		if args[0] == "api" && args[1] == "-X" && args[2] == "PUT" && args[3] == "repos/owner/repo/actions/permissions/workflow" && args[4] == "-F" && args[5] == "default_workflow_permissions=read" {
			putCalled = true
			return []byte(`{}`), nil
		}

		return nil, fmt.Errorf("unexpected command call: %v", args)
	}

	settings := []*pb.GithubSetting{
		{Key: "delete_branch_on_merge", Value: "true"},
		{Key: "default_workflow_permissions", Value: "read"},
	}

	err := ReconcileGithubSettings(context.Background(), "owner/repo", settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !patchCalled {
		t.Errorf("expected PATCH to delete_branch_on_merge to be called")
	}
	if !putCalled {
		t.Errorf("expected PUT to default_workflow_permissions to be called")
	}
}

func TestReconcileGithubSettings_Errors(t *testing.T) {
	oldRunCommand := runCommand
	defer func() { runCommand = oldRunCommand }()

	// Test 1: Unknown setting
	settings := []*pb.GithubSetting{
		{Key: "unknown_setting_key", Value: "true"},
	}
	err := ReconcileGithubSettings(context.Background(), "owner/repo", settings)
	if err == nil || !strings.Contains(err.Error(), "unknown GitHub setting") {
		t.Errorf("expected unknown setting error, got %v", err)
	}

	// Test 2: Invalid boolean
	settings = []*pb.GithubSetting{
		{Key: "delete_branch_on_merge", Value: "not-a-bool"},
	}
	// Setup GET response for the test
	runCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return []byte(`{"delete_branch_on_merge": false}`), nil
	}
	err = ReconcileGithubSettings(context.Background(), "owner/repo", settings)
	if err == nil || !strings.Contains(err.Error(), "invalid boolean value") {
		t.Errorf("expected invalid boolean error, got %v", err)
	}
}
