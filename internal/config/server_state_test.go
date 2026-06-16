package config

import (
	"testing"

	pb "github.com/brotherlogic/seraphine/proto"
)

func TestServerStateCompile(t *testing.T) {
	state := &pb.ServerState{
		EnrolledRepositories: []string{"repo1", "repo2"},
	}
	if len(state.GetEnrolledRepositories()) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(state.GetEnrolledRepositories()))
	}
}
