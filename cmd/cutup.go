package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/cutup"
)

func init() {
	rootCmd.AddCommand(cutupCmd)
}

var cutupCmd = &cobra.Command{
	Use:  "cutup [prefix]",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO arg for source file path
		// TODO arg for target path
		return cutup.CutupFiles()
	},
}
