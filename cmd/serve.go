package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/web"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		return web.Serve()
	},
}
