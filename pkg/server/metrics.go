package server

import (
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
)

func NewMetricsServer(addr string) (*http.Server, func(next echo.HandlerFunc) echo.HandlerFunc) {
	e := echo.New()

	prom := prometheus.NewPrometheus("echo", nil)
	prom.SetMetricsPath(e)

	return &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadHeaderTimeout: requestTimeout,
	}, prom.HandlerFunc
}
