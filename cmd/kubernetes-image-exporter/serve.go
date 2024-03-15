package main

import (
	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/Thibaut-gauvin/kie/pkg/kie"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start kie server.",
	Run: func(cmd *cobra.Command, _ []string) {
		bindPFlagsSnakeCase(cmd.Flags())

		opts := kie.ServeOpts{}
		hydrateOptsFromViper(&opts)

		if err := doServe(opts); err != nil {
			logger.Fatalf("Error starting kie server. message: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("listen-port", "p", "9145", "The server listening port.")
}

func doServe(opts kie.ServeOpts) error {
	logger.Debugf("%+v", opts)

	return kie.StartServer(opts)
}
