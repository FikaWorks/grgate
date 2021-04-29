package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/ggate/pkg/config"
)

var (
  cfgFile string
  globalConfig *viper.Viper

  // Version of GGate
  Version string
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "ggate",
	Short: "Publish draft/unpublished releases if all status check succeed",
	Long: `GGate is git release gate utility which autopublish draft/unpublished
releases based on commit status (aka checks). It can be triggered automatically
using Git webhook or directly from the CLI.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
    "config file (default is $HOME/.ggate.yaml)")

	rootCmd.PersistentFlags().Int64("github.appID", 0, "Github App ID")
	rootCmd.PersistentFlags().Int64("github.installationID", 0,
    "Github Installation ID")
	rootCmd.PersistentFlags().String("github.privateKeyPath", "",
    "Github private key path")
	rootCmd.PersistentFlags().String("github.webhookSecret", "",
    "Github webhook secret")
	rootCmd.PersistentFlags().String("logLevel", "info",
    "Log level: trace, debug, info, warn, error, fatal or panic")
	rootCmd.PersistentFlags().String("logFormat", "pretty",
    "Log format: json or pretty")
}

func initConfig() {
  // read global config and override it with flags value
  globalConfig, _ := config.NewGlobalConfig(cfgFile)

  globalConfig.BindPFlags(rootCmd.PersistentFlags())
  globalConfig.Unmarshal(&config.Main)

	// logs
  logLevel, _ := zerolog.ParseLevel(config.Main.LogLevel)
  zerolog.SetGlobalLevel(logLevel)

  if config.Main.LogFormat == "pretty" {
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
  }

	zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

  // inform about which config file is being used
  configFile := globalConfig.ConfigFileUsed()
  if configFile != "" {
    log.Info().Msgf("Using config file: %s", configFile)
  }
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
