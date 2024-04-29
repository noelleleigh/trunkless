package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/cutup"
)

func init() {
	rootCmd.Flags().StringP("cutupdir", "d", "", "directory in which to write phrase files")
	rootCmd.Flags().StringP("srcdir", "s", "", "directory of files to cut up")
	rootCmd.Flags().IntP("workers", "w", 10, "number of workers to use when cutting up")

	rootCmd.MarkFlagRequired("cutupdir")
	rootCmd.MarkFlagRequired("srcdir")

	rootCmd.AddCommand(cutupCmd)
}

var cutupCmd = &cobra.Command{
	Use:  "cutup",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cutupdir, _ := cmd.Flags().GetString("cutupdir")
		srcdir, _ := cmd.Flags().GetString("srcdir")
		workers, _ := cmd.Flags().GetInt("workers")
		opts := cutup.CutupOpts{
			CutupDir:   cutupdir,
			SrcDir:     srcdir,
			NumWorkers: workers,
		}
		return cutup.Cutup(opts)
	},
}
