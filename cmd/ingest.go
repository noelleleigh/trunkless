package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/db"
	"github.com/vilmibm/trunkless/ingest"
)

func init() {
	ingestCmd.Flags().StringP("cutupdir", "d", "", "directory to files produced by cutup cmd")
	ingestCmd.MarkFlagRequired("cutupdir")
	rootCmd.AddCommand(ingestCmd)
}

var ingestCmd = &cobra.Command{
	Use:   "ingest corpusname",
	Short: "ingest already cut-up corpora from disk into database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cutupDir, _ := cmd.Flags().GetString("cutupdir")
		corpus := args[0]

		conn, err := db.Connect()
		if err != nil {
			return err
		}

		opts := ingest.IngestOpts{
			Conn:     conn,
			CutupDir: cutupDir,
			Corpus:   corpus,
		}

		return ingest.Ingest(opts)
	},
}
