package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/Thibaut-gauvin/kie/pkg/kie"
	"github.com/spf13/cobra"
)

const (
	listenPortDefault      = "9145"
	kubeconfigDefault      = ""
	refreshIntervalDefault = 30
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
	serveCmd.Flags().StringP("listen-port", "p", listenPortDefault,
		"the listening port that kie server will use.")
	serveCmd.Flags().StringP("kubeconfig", "k", kubeconfigDefault,
		"path to kubeconfig used to authenticate with cluster when your running kie locally. "+
			"If not provided, use service-account from pod.")
	serveCmd.Flags().IntP("refresh-interval", "i", refreshIntervalDefault,
		"metrics values refresh interval in seconds.")
}

func doServe(opts kie.ServeOpts) error {
	logger.Debugf("%+v", opts)

	if err := validateListenPort(opts.ListenPort); err != nil {
		return err
	}
	if err := validateKubeconfigPath(opts.Kubeconfig); opts.Kubeconfig != "" && err != nil {
		return err
	}
	if err := validateRefreshInterval(opts.RefreshInterval); err != nil {
		return err
	}

	return kie.StartKie(opts)
}

// validateListenPort ensure that given ListenPort option is valid.
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

// validateKubeconfigPath ensure that given Kubeconfig option is valid.
func validateKubeconfigPath(kubeconfigPath string) error {
	if kubeconfigPath == "" {
		return nil
	}
	_, err := os.Stat(kubeconfigPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", kubeconfigPath)
	}
	if err != nil {
		return fmt.Errorf("%s is invalid, %w", kubeconfigPath, err)
	}
	return nil
}

// validateRefreshInterval ensure that given refreshInterval option is valid.
func validateRefreshInterval(refreshInterval int) error {
	minI := 1     // 1 second
	maxI := 86400 // 24 hours
	if minI > refreshInterval || maxI < refreshInterval {
		return fmt.Errorf("%d is not a valid refresh interval", refreshInterval)
	}
	return nil
}
