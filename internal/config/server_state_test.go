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
