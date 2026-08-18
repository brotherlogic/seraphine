package proto_test

import (
	"testing"

	pb "github.com/brotherlogic/seraphine/proto"
)

func TestPRContainer(t *testing.T) {
	container := &pb.PRContainer{
		Repo:        "brotherlogic/seraphine",
		PrNumber:    149,
		ContainerId: "test-container-id",
	}

	if container.GetRepo() != "brotherlogic/seraphine" {
		t.Errorf("expected repo brotherlogic/seraphine, got %v", container.GetRepo())
	}
	if container.GetPrNumber() != 149 {
		t.Errorf("expected pr_number 149, got %v", container.GetPrNumber())
	}
	if container.GetContainerId() != "test-container-id" {
		t.Errorf("expected container_id test-container-id, got %v", container.GetContainerId())
	}

	state := &pb.ServerState{
		PrContainers: []*pb.PRContainer{container},
	}

	if len(state.GetPrContainers()) != 1 {
		t.Fatalf("expected 1 pr container, got %v", len(state.GetPrContainers()))
	}
	if state.GetPrContainers()[0].GetContainerId() != "test-container-id" {
		t.Errorf("expected container id test-container-id, got %v", state.GetPrContainers()[0].GetContainerId())
	}
}
