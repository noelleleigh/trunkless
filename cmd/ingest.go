package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/ingest"
)

func init() {
	rootCmd.AddCommand(ingestCmd)
}

var ingestCmd = &cobra.Command{
	Use: "ingest",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ingest.Ingest(strings.Join(args[1:], " "), os.Stdin)
	},
}
