package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/ingest"
)

func init() {
	rootCmd.AddCommand(ingestCmd)
}

var ingestCmd = &cobra.Command{
	Use:  "ingest corpusname",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "gutenberg":
			return ingest.IngestGut()
		}
		return fmt.Errorf("corpus unknown: %s", args[0])
	},
}
