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
	"github.com/fikaworks/ggate/pkg/workers"
)

var (
	version string
)

// ServerConfig hold configuration to run a server
type ServerConfig struct {
  ListenAddr string
  Logger zerolog.Logger
  MetricsAddr string
  Platform platforms.Platform
  ProbeAddr string
  WebhookSecret string
  Workers int
}

// Server hold a server instance
type Server struct {
  CancelWorker chan struct{}
  Config *ServerConfig
  Srv *http.Server
  WorkerPool *workers.WorkerPool
}

// NewServer returns an instance of server
func NewServer(config *ServerConfig) *Server {
	// router
	r := mux.NewRouter()
	c := alice.New(hlog.NewHandler(config.Logger), hlog.AccessHandler(Logger))

  cancelWorker := make(chan struct{})
  workerPool := workers.NewWorkerPool(config.Workers, cancelWorker)

  webhook := NewWebhookHandler(config.Platform, config.WebhookSecret,
    workerPool.JobQueue)
  r.HandleFunc("/github/webhook", webhook.GithubHandler).Methods("POST")

	srv := &http.Server{
    Addr: config.ListenAddr,
    Handler: c.Then(PromRequestHandler(r)),
  }

  return &Server{
    Config: config,
    Srv: srv,
    WorkerPool: workerPool,
    CancelWorker: cancelWorker,
  }
}

func (s *Server) startWorkerPool() {
  s.WorkerPool.Start()
}

// Start all components required to run a server
func (s *Server) Start() {
  // handle graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go s.startWorkerPool()
	go s.serveMetrics()
	go s.serveHTTP()
	go s.serveProbe()

	<-quit

	log.Info().Msg("Shutting down worker pool...")
  close(s.CancelWorker)

	log.Info().Msg("Shutting down server...")

	// Gracefully shutdown connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Srv.Shutdown(ctx)
}

func (s *Server) serveHTTP() {
	log.Info().
    Msgf("Server started at %s", s.Config.ListenAddr)

  err := s.Srv.ListenAndServe();
	if err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Failed starting HTTP server")
	}
}

func (s *Server) serveProbe() {
	log.Info().
    Msgf("Probe server running at %s", s.Config.ProbeAddr)

  err := http.ListenAndServe(s.Config.ProbeAddr, healthcheck.NewHandler())
	if err != nil {
		log.Error().Err(err).Msg("Starting probe listener failed")
  }
}

func (s *Server) serveMetrics() {
	log.Info().
    Msgf("Serving Prometheus metrics on port %s", s.Config.MetricsAddr)

	http.Handle("/metrics", promhttp.Handler())

  err := http.ListenAndServe(s.Config.MetricsAddr, nil)
	if err != nil {
		log.Error().Err(err).Msg("Starting Prometheus listener failed")
	}
}
