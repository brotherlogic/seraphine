package cmd

import (
	"fmt"
	"os"

	"github.com/brotherlogic/seraphine/internal/server"
	"github.com/spf13/cobra"
)

var (
	httpPort                   string
	devcontainerManagerAddress string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Seraphine gRPC and HTTP servers",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Run(":9009", httpPort, devcontainerManagerAddress)
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serverCmd.Flags().StringVar(&httpPort, "http-port", "", "HTTP port for web dashboard and health checks (defaults to $HTTP_PORT or :8080)")
	serverCmd.Flags().StringVar(&devcontainerManagerAddress, "devcontainer-manager-address", "", "devcontainer manager gRPC address (defaults to $DEVCONTAINER_MANAGER_ADDRESS or devcontainer-manager.devcontainer-manager.svc.cluster.local:8080)")
	serverCmd.Flags().StringVar(&devcontainerManagerAddress, "devcontainer_manager_address", "", "devcontainer manager gRPC address (defaults to $DEVCONTAINER_MANAGER_ADDRESS or devcontainer-manager.devcontainer-manager.svc.cluster.local:8080)")
	rootCmd.AddCommand(serverCmd)
}
