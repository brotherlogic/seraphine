package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seraphine",
	Short: "Seraphine is an agentic workflow manager",
	Long:  `Seraphine is an agentic workflow manager written in Go.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
