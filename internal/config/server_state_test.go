package config

import (
	"context"
	"testing"

	pstore_client "github.com/brotherlogic/pstore/client"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/protobuf/proto"
)

func TestServerStateCompileAndSerialization(t *testing.T) {
	state := &pb.ServerState{
		EnrolledRepositories: []string{"repo1", "repo2"},
	}
	if len(state.GetEnrolledRepositories()) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(state.GetEnrolledRepositories()))
	}

	// Serialize the state
	data, err := proto.Marshal(state)
	if err != nil {
		t.Fatalf("proto.Marshal failed: %v", err)
	}

	// Deserialize the state
	newState := &pb.ServerState{}
	err = proto.Unmarshal(data, newState)
	if err != nil {
		t.Fatalf("proto.Unmarshal failed: %v", err)
	}

	// Verify content
	repos := newState.GetEnrolledRepositories()
	if len(repos) != 2 || repos[0] != "repo1" || repos[1] != "repo2" {
		t.Errorf("Deserialized state mismatch: expected [repo1, repo2], got %v", repos)
	}
}

func TestReadWriteServerState(t *testing.T) {
	ctx := context.Background()
	client := pstore_client.GetTestClient()

	// Initial read should return empty/default state
	state, err := ReadServerState(ctx, client)
	if err != nil {
		t.Fatalf("Failed to read initial server state: %v", err)
	}
	if len(state.GetEnrolledRepositories()) != 0 {
		t.Errorf("Expected empty repositories, got %v", state.GetEnrolledRepositories())
	}

	// Write new state
	newState := &pb.ServerState{
		EnrolledRepositories: []string{"repo1", "repo2"},
	}
	err = WriteServerState(ctx, client, newState)
	if err != nil {
		t.Fatalf("Failed to write server state: %v", err)
	}

	// Read it back
	state2, err := ReadServerState(ctx, client)
	if err != nil {
		t.Fatalf("Failed to read back server state: %v", err)
	}
	if len(state2.GetEnrolledRepositories()) != 2 {
		t.Fatalf("Expected 2 repositories, got %d", len(state2.GetEnrolledRepositories()))
	}
	if state2.GetEnrolledRepositories()[0] != "repo1" || state2.GetEnrolledRepositories()[1] != "repo2" {
		t.Errorf("Mismatch in repositories: %v", state2.GetEnrolledRepositories())
	}
}

func TestPRContainerHelpers(t *testing.T) {
	state := &pb.ServerState{}

	// Initially finding container should return nil
	if c := FindPRContainer(state, "owner/repo", 10); c != nil {
		t.Errorf("Expected nil, got %v", c)
	}

	// Add a PR container
	AddPRContainer(state, &pb.PRContainer{
		Repo:        "owner/repo",
		PrNumber:    10,
		ContainerId: "container-123",
	})

	if len(state.GetPrContainers()) != 1 {
		t.Fatalf("Expected 1 PR container, got %d", len(state.GetPrContainers()))
	}

	// Find the added container
	c := FindPRContainer(state, "owner/repo", 10)
	if c == nil {
		t.Fatalf("Expected container, got nil")
	}
	if c.GetContainerId() != "container-123" {
		t.Errorf("Expected container_id 'container-123', got '%s'", c.GetContainerId())
	}

	// Finding non-existent repo or PR number should return nil
	if FindPRContainer(state, "owner/other", 10) != nil {
		t.Errorf("Expected nil for different repo")
	}
	if FindPRContainer(state, "owner/repo", 99) != nil {
		t.Errorf("Expected nil for different PR number")
	}

	// Updating existing container should overwrite instead of duplicate
	AddPRContainer(state, &pb.PRContainer{
		Repo:        "owner/repo",
		PrNumber:    10,
		ContainerId: "container-456",
	})

	if len(state.GetPrContainers()) != 1 {
		t.Fatalf("Expected 1 PR container after update, got %d", len(state.GetPrContainers()))
	}
	if c := FindPRContainer(state, "owner/repo", 10); c == nil || c.GetContainerId() != "container-456" {
		t.Errorf("Expected updated container_id 'container-456', got %v", c)
	}

	// Remove container
	RemovePRContainer(state, "owner/repo", 10)
	if len(state.GetPrContainers()) != 0 {
		t.Errorf("Expected 0 PR containers after removal, got %d", len(state.GetPrContainers()))
	}
	if FindPRContainer(state, "owner/repo", 10) != nil {
		t.Errorf("Expected nil after removal")
	}

	// Removing non-existent container should be a no-op
	RemovePRContainer(state, "owner/repo", 10)
	if len(state.GetPrContainers()) != 0 {
		t.Errorf("Expected 0 PR containers after removing non-existent record, got %d", len(state.GetPrContainers()))
	}
}

