package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/server"
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
		platform, err := newPlatform()
		if err != nil {
			return
		}

		srv := server.NewServer(&server.Config{
			ListenAddr:    config.Main.Server.ListenAddress,
			Logger:        log.Logger,
			MetricsAddr:   config.Main.Server.MetricsAddress,
			Platform:      platform,
			ProbeAddr:     config.Main.Server.ProbeAddress,
			WebhookSecret: config.Main.Server.WebhookSecret,
			Workers:       config.Main.Workers,
		})
		srv.Start()

		return
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	flags := serveCmd.PersistentFlags()

	flags.String("server.listenAddress", "0.0.0.0:8080",
		"The address to listen on for HTTP requests")
	flags.String("server.metricsAddress", "0.0.0.0:9101",
		"The address to listen on for Prometheus metrics requests")
	flags.String("server.probeAddress", "0.0.0.0:8086",
		"The address to listen on for probe requests")
	flags.IntP("workers", "w", 5, "Number of workers to run")
}
