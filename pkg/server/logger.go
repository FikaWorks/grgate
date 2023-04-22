package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func Logger(_ echo.Context, v middleware.RequestLoggerValues) error {
	if v.Error == nil {
		log.Info().
			Str("URI", v.URI).
			Int("status", v.Status).
			Msg("request")
	} else {
		log.Error().
			Err(v.Error).
			Str("URI", v.URI).
			Int("status", v.Status).
			Msg("request error")
	}
	return nil
}
