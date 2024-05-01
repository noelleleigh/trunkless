package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/cutup"
)

func init() {
	cutupCmd.Flags().StringP("cutupdir", "d", "", "directory in which to write phrase files")
	cutupCmd.Flags().StringP("srcdir", "s", "", "directory of files to cut up")
	cutupCmd.Flags().IntP("workers", "w", 10, "number of workers to use when cutting up")
	cutupCmd.Flags().StringP("flavor", "f", "", "set of adapters to use when cutting up")

	cutupCmd.MarkFlagRequired("cutupdir")
	cutupCmd.MarkFlagRequired("srcdir")

	rootCmd.AddCommand(cutupCmd)
}

var validFlavors = []string{"gutenberg"}

var cutupCmd = &cobra.Command{
	Use:  "cutup",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cutupdir, _ := cmd.Flags().GetString("cutupdir")
		srcdir, _ := cmd.Flags().GetString("srcdir")
		workers, _ := cmd.Flags().GetInt("workers")
		flavor, _ := cmd.Flags().GetString("flavor")

		if flavor != "" {
			valid := false
			for _, f := range validFlavors {
				if flavor == f {
					valid = true
				}
			}
			if !valid {
				return fmt.Errorf("invalid flavor '%s'; valid flavors: %v", flavor, validFlavors)
			}
		}
		opts := cutup.CutupOpts{
			CutupDir:   cutupdir,
			SrcDir:     srcdir,
			NumWorkers: workers,
			Flavor:     flavor,
		}
		return cutup.Cutup(opts)
	},
}
