package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/grgate/pkg/config"
)

var (
	cfgFile string
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "grgate",
	Short: "Publish draft/unpublished releases if all status check succeed",
	Long: `GRGate is a git release gate utility which autopublish
draft/unpublished releases based on commit status (aka checks). It can be
triggered automatically using Git webhook or directly from the CLI.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&cfgFile, "config", "c", "",
		"config file (default is /etc/grgate/config.yaml)")
	flags.Int64("github.appID", 0, "Github App ID")
	flags.Int64("github.installationID", 0, "Github Installation ID")
	flags.String("github.privateKeyPath", "", "Github private key path")
	flags.String("github.webhookSecret", "", "Github webhook secret")
	flags.String("gitlab.token", "", "Gitlab Token")
	flags.String("logLevel", "info", "Log level: trace, debug, info, warn,"+
		"error, fatal or panic")
	flags.String("logFormat", "pretty", "Log format: json or pretty")
	flags.String("platform", "github", "Platform to run against: github or gitlab"+
		"(default: github)")
}

func initConfig() {
	// read global config and override it with flags value
	globalConfig, err := config.NewGlobalConfig(cfgFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
		return
	}

	if err := globalConfig.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Print(err)
		os.Exit(1)
		return
	}

	if err := globalConfig.Unmarshal(&config.Main); err != nil {
		fmt.Print(err)
		os.Exit(1)
		return
	}

	// logs
	logLevel, err := zerolog.ParseLevel(config.Main.LogLevel)
	if err != nil {
		log.Error().Err(err)
		return
	}

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
