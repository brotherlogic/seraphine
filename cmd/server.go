package cmd

import (
	"fmt"
	"os"

	"github.com/brotherlogic/seraphine/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Seraphine gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Run(":9009")
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
