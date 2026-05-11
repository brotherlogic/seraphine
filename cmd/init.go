package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/brotherlogic/seraphine/internal/workflow"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Seraphine project",
	Run: func(cmd *cobra.Command, args []string) {
		err := workflow.RunSync(context.Background(), "seraphine.brotherlogic-backend.com:9009")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
