package config

import (
	"os"

	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/protobuf/encoding/prototext"
)

const ConfigPath = ".seraphine/config.textpb"

func ReadConfig() (*pb.ProjectConfig, error) {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &pb.ProjectConfig{}, nil
		}
		return nil, err
	}
	var cfg pb.ProjectConfig
	err = prototext.Unmarshal(data, &cfg)
	return &cfg, err
}

func WriteConfig(cfg *pb.ProjectConfig) error {
	err := os.MkdirAll(".seraphine", 0755)
	if err != nil {
		return err
	}
	data, err := prototext.MarshalOptions{Multiline: true}.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath, data, 0644)
}
