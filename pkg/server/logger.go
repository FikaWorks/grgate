package server

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
)

func Logger(r *http.Request, status, size int, dur time.Duration) {
	hlog.FromRequest(r).Info().
		Str("host", r.Host).
		Int("status", status).
		Int("size", size).
		Dur("duration_ms", dur).
		Msg("request")
}
