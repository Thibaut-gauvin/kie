package main

import (
	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start kie server.",
	Run: func(_ *cobra.Command, _ []string) {
		if err := doServe(); err != nil {
			logger.Fatalf("Error starting kie server. message: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func doServe() error {
	logger.Debugf("Lorem ipsum")
	logger.Infof("Lorem ipsum")

	return nil
}
