package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/brotherlogic/seraphine/internal/workflow"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade an existing Seraphine project",
	Run: func(cmd *cobra.Command, args []string) {
		err := workflow.RunSync(context.Background(), "seraphine.brotherlogic-backend.com:9009")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
