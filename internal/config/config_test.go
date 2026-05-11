package config

import (
	"os"
	"testing"
)

func TestConfigReadWrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "seraphine-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	cfg := &Config{Version: "12345"}
	err = WriteConfig(cfg)
	if err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	readCfg, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	if readCfg.Version != "12345" {
		t.Errorf("Version mismatch: expected 12345, got %s", readCfg.Version)
	}
}
