package cmd

import (
	"testing"
)

func TestServerCmdFlags(t *testing.T) {
	for _, flagName := range []string{"http-port", "devcontainer-manager-address", "devcontainer_manager_address"} {
		flag := serverCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Fatalf("Expected --%s flag to be registered on serverCmd", flagName)
		}
		if flag.DefValue != "" {
			t.Errorf("Expected default value for --%s to be empty string, got %s", flagName, flag.DefValue)
		}
	}
}
