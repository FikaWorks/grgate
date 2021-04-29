package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fikaworks/ggate/pkg/server"
	"github.com/fikaworks/ggate/pkg/config"
	"github.com/fikaworks/ggate/pkg/platforms"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run HTTP server to receive git webhook event",
	Long: `The serve command create 3 HTTP server with the following
functionnalities:
  - 0.0.0.0:8080 listen to git webhook
  - 0.0.0.0:9101 expose Prometheus metrics
  - 0.0.0.0:8086 expose health probe (liveness/readiness)
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
    platform, err := platforms.NewGithub(&platforms.GithubConfig{
      AppID: config.Main.Github.AppID,
      InstallationID: config.Main.Github.InstallationID,
      PrivateKeyPath: config.Main.Github.PrivateKeyPath,
    })
    if err != nil {
      return err
    }

    srv := server.NewServer(&server.ServerConfig{
      ListenAddr: config.Main.Server.ListenAddress,
      MetricsAddr: config.Main.Server.MetricsAddress,
      Platform: platform,
      ProbeAddr: config.Main.Server.ProbeAddress,
      WebhookSecret: config.Main.Github.WebhookSecret,
    })
    srv.Start()

    return
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().String("server.listenAddress",
    "0.0.0.0:8080", "The address to listen on for HTTP requests")
	serveCmd.PersistentFlags().String("server.metricsAddress",
    "0.0.0.0:9101", "The address to listen on for Prometheus metrics requests")
	serveCmd.PersistentFlags().String("server.probeAddress",
    "0.0.0.0:8086", "The address to listen on for probe requests")
	serveCmd.PersistentFlags().IntP("workers",
    "w", 5, "Number of workers to run")
}

