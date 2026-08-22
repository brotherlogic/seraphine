package cmd

import (
	"testing"
)

func TestServerCmdFlags(t *testing.T) {
	flag := serverCmd.Flags().Lookup("http-port")
	if flag == nil {
		t.Fatalf("Expected --http-port flag to be registered on serverCmd")
	}
	if flag.DefValue != "" {
		t.Errorf("Expected default value to be empty string, got %s", flag.DefValue)
	}
}
