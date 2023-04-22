package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/workers"
)

const (
	// shutdownTimeout is the number of seconds to wait before shutting
	// down the server
	shutdownTimeout time.Duration = 5 * time.Second

	// requestTimeout is the number of seconds to wait before dropping
	// connection
	requestTimeout time.Duration = 3 * time.Second
)

// Config hold configuration to run a server
type Config struct {
	ListenAddr    string
	Logger        zerolog.Logger
	MetricsAddr   string
	Platform      platforms.Platform
	ProbeAddr     string
	WebhookSecret string
	Workers       int
}

// Server hold a server instance
type Server struct {
	CancelWorker  chan struct{}
	Config        *Config
	MainServer    *http.Server
	MetricsServer *http.Server
	ProbeServer   *http.Server
	WorkerPool    *workers.WorkerPool
}

// NewServer returns an instance of server
func NewServer(config *Config) *Server {
	// Probe server
	probeServer := NewProbeServer(config.ProbeAddr)

	// Prometheus metrics middleware
	metricsServer, promHandlerFunc := NewMetricsServer(config.MetricsAddr)

	// Main server
	e := echo.New()
	e.HideBanner = true

	// Prometheus metrics
	e.Use(promHandlerFunc)

	// logging
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogError:      true,
		HandleError:   true,
		LogValuesFunc: Logger,
	}))

	cancelWorker := make(chan struct{})
	workerPool := workers.NewWorkerPool(config.Workers, cancelWorker)

	webhook := NewWebhookHandler(config.Platform, config.WebhookSecret,
		workerPool.JobQueue)

	e.POST("/github/webhook", webhook.GithubHandler)
	e.POST("/gitlab/webhook", webhook.GitlabHandler)

	mainServer := &http.Server{
		Addr:              config.ListenAddr,
		Handler:           e,
		ReadHeaderTimeout: requestTimeout,
	}

	return &Server{
		Config:        config,
		MainServer:    mainServer,
		MetricsServer: metricsServer,
		ProbeServer:   probeServer,
		WorkerPool:    workerPool,
		CancelWorker:  cancelWorker,
	}
}

func (s *Server) startWorkerPool() {
	s.WorkerPool.Start()
}

// Start all components required to run a server
func (s *Server) Start() {
	// handle graceful shutdown
	quit := make(chan os.Signal, 1)
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
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.MainServer.Shutdown(ctx); err != nil {
		log.Error().Err(err)
	}
}

func (s *Server) serveHTTP() {
	log.Info().
		Msgf("Server started at %s", s.Config.ListenAddr)

	err := s.MainServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Failed starting HTTP server")
	}
}

func (s *Server) serveProbe() {
	log.Info().
		Msgf("Probe server running at %s", s.Config.ProbeAddr)

	err := s.ProbeServer.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("Starting probe listener failed")
	}
}

func (s *Server) serveMetrics() {
	log.Info().
		Msgf("Serving Prometheus metrics on port %s", s.Config.MetricsAddr)

	err := s.MetricsServer.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("Starting Prometheus listener failed")
	}
}
