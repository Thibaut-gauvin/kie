package main

import (
	"fmt"
	"strconv"

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

	if err := validateListenPort(opts.ListenPort); err != nil {
		return err
	}

	return kie.StartServer(opts)
}

func validateListenPort(listenPort string) error {
	portNumber, err := strconv.Atoi(listenPort)
	if err != nil {
		return fmt.Errorf("%s is not a valid port number", listenPort)
	}
	if portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("%s is not a valid port number", listenPort)
	}
	return nil
}
