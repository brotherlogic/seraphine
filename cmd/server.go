package cmd

import (
	"fmt"
	"os"

	"github.com/brotherlogic/seraphine/internal/server"
	"github.com/spf13/cobra"
)

var (
	httpPort string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Seraphine gRPC and HTTP servers",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Run(":9009", httpPort)
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serverCmd.Flags().StringVar(&httpPort, "http-port", "", "HTTP port for web dashboard and health checks (defaults to $HTTP_PORT or :8080)")
	rootCmd.AddCommand(serverCmd)
}
