package config

import (
	"context"

	pstore_client "github.com/brotherlogic/pstore/client"
	pstore_pb "github.com/brotherlogic/pstore/proto"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const ServerStateKey = "seraphine/server_state"

// ReadServerState reads the server state from pstore
func ReadServerState(ctx context.Context, client pstore_client.PStoreClient) (*pb.ServerState, error) {
	res, err := client.Read(ctx, &pstore_pb.ReadRequest{Key: ServerStateKey})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &pb.ServerState{}, nil
		}
		return nil, err
	}

	state := &pb.ServerState{}
	err = proto.Unmarshal(res.GetValue().GetValue(), state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// WriteServerState writes the server state to pstore
func WriteServerState(ctx context.Context, client pstore_client.PStoreClient, state *pb.ServerState) error {
	anyVal, err := anypb.New(state)
	if err != nil {
		return err
	}

	_, err = client.Write(ctx, &pstore_pb.WriteRequest{
		Key:   ServerStateKey,
		Value: anyVal,
	})
	return err
}

// FindPRContainer finds a PRContainer record matching the repo and pr_number.
// Returns nil if state is nil or record is not found.
func FindPRContainer(state *pb.ServerState, repo string, prNumber int32) *pb.PRContainer {
	if state == nil {
		return nil
	}
	for _, c := range state.GetPrContainers() {
		if c.GetRepo() == repo && c.GetPrNumber() == prNumber {
			return c
		}
	}
	return nil
}

// AddPRContainer adds or updates a PRContainer record in ServerState.
func AddPRContainer(state *pb.ServerState, container *pb.PRContainer) {
	if state == nil || container == nil {
		return
	}
	for i, c := range state.GetPrContainers() {
		if c.GetRepo() == container.GetRepo() && c.GetPrNumber() == container.GetPrNumber() {
			state.PrContainers[i] = container
			return
		}
	}
	state.PrContainers = append(state.PrContainers, container)
}

// RemovePRContainer removes a PRContainer record matching the repo and pr_number from ServerState.
func RemovePRContainer(state *pb.ServerState, repo string, prNumber int32) {
	if state == nil {
		return
	}
	n := 0
	for _, c := range state.GetPrContainers() {
		if c.GetRepo() == repo && c.GetPrNumber() == prNumber {
			continue
		}
		state.PrContainers[n] = c
		n++
	}
	state.PrContainers = state.PrContainers[:n]
}

