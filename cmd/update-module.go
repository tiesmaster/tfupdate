package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateModuleCmd = &cobra.Command{
	Use:   "update-module MODULE_ID VERSION",
	Short: "Update module",
	RunE:  runUpdateModuleCommand,
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(updateModuleCmd)
}

func runUpdateModuleCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("Will update module '%s' to version '%s'\n", args[0], args[1])
	return nil
}
