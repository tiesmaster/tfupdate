package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Subcommand to update various things (Terraform, providers, modules)",
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
