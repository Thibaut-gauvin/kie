package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel = "info"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use: "kie",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Short: "kie exports Prometheus metrics about image usage in your cluster",
	Long: `kubernetes-image-exporter help monitors images running on your cluster by exporting Prometheus metrics

Run kie --help for more information`,
}

// Execute runs the root cobra command.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init initialize Cobra root command.
func init() {
	cobra.OnInitialize(initConfig, initLogLevel)

	rootCmd.PersistentFlags().StringP("log-level", "l", defaultLogLevel,
		"Log level. Can be any standard log-level (\"info\", \"debug\", etc...)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("log_level", defaultLogLevel)
}

// initLogLevel parse log-level from cmd flag, instantiate slog logger & use it as default logger.
func initLogLevel() {
	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	logLevel := viper.GetString("log_level")
	logLvl, err := logger.ParseLevel(logLevel)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	internalLogger := logger.NewLogger(os.Stdout, &slog.HandlerOptions{Level: logLvl})
	slog.SetDefault(internalLogger)
}

// hydrateOptsFromViper copies all the viper values into our config struct.
// The mapping between viper identifiers and struct field names
// is ensured by `mapstructure` struct tags.
func hydrateOptsFromViper(opts interface{}) {
	_ = viper.Unmarshal(opts)
}

// bindPFlagsSnakeCase binds the flags with viper values. The identifier of the viper value
// is the name of the flag with dashes replaced by underscores. This is required so we can
// retrieve values from viper with the same behaviour with config coming from files
// (my_config: "value") or from flags (--my-config=value).
func bindPFlagsSnakeCase(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(strings.ReplaceAll(flag.Name, "-", "_"), flag)
	})
}
