package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/cutup"
)

func init() {
	rootCmd.AddCommand(cutupCmd)
}

var cutupCmd = &cobra.Command{
	Use: "cutup",
	Run: func(cmd *cobra.Command, args []string) {
		cutup.Cutup(os.Stdin)
	},
}
