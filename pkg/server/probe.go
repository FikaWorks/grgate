package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewProbeServer(addr string) *http.Server {
	e := echo.New()

	e.GET("/live", livenessProbeHandler)
	e.GET("/ready", readinessProbeHandler)

	return &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadHeaderTimeout: requestTimeout,
	}
}

func livenessProbeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Ok")
}

func readinessProbeHandler(c echo.Context) error {
	// TODO: should check dependencies or third party provider availability
	return c.String(http.StatusOK, "Ok")
}
