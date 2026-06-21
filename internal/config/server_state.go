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
