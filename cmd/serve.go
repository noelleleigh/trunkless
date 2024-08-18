package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vilmibm/trunkless/web"
)

func init() {
	serveCmd.Flags().IntP("port", "p", 8080, "port to listen on")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")

		opts := web.ServeOpts{
			Port: port,
		}

		return web.Serve(opts)
	},
}
