package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version string `yaml:"version"`
}

const ConfigPath = ".seraphine/config.yaml"

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func WriteConfig(cfg *Config) error {
	err := os.MkdirAll(".seraphine", 0755)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath, data, 0644)
}
