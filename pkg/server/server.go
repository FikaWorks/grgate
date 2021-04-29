package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/ggate/pkg/platforms"
)

var (
	version string
)

type ServerConfig struct {
  ListenAddr string
  MetricsAddr string
  ProbeAddr string
  Platform platforms.Platform
  WebhookSecret string
}

type Server struct {
  Config *ServerConfig
  Srv *http.Server
}

func init() {
	// prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)
	version = os.Getenv("VERSION")
}

func NewServer(config *ServerConfig) *Server {
	// init log
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("version", version).
		Logger()

	// router
	r := mux.NewRouter()
	c := alice.New(hlog.NewHandler(log), hlog.AccessHandler(Logger))

  webhook := NewWebhookHandler(config.Platform, config.WebhookSecret)
  r.HandleFunc("/github/webhook", webhook.GithubHandler).Methods("POST")

	srv := &http.Server{
    Addr: config.ListenAddr,
    Handler: c.Then(PromRequestHandler(r)),
  }

  return &Server{
    Config: config,
    Srv: srv,
  }
}

func (s *Server) Start() {
  // handle graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go s.serveMetrics()
	go s.serveHTTP()
	go s.serveProbe()

	<-quit

	log.Info().Msg("Shutting down server...")

	// Gracefully shutdown connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Srv.Shutdown(ctx)
}

func (s *Server) serveHTTP() {
	log.Info().Msgf("Server started at %s", s.Config.ListenAddr)
	err := s.Srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Failed starting HTTP server")
	}
}

func (s *Server) serveProbe() {
	log.Info().Msgf("Probe server running at %s", s.Config.ProbeAddr)
	http.ListenAndServe(s.Config.ProbeAddr, healthcheck.NewHandler())
}

func (s *Server) serveMetrics() {
	log.Info().Msgf("Serving Prometheus metrics on port %s", s.Config.MetricsAddr)

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(s.Config.MetricsAddr, nil); err != nil {
		log.Error().Err(err).Msg("Starting Prometheus listener failed")
	}
}
