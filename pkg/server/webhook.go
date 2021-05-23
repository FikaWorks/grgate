package server

import (
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/workers"
)

// WebhookHandler hold webhook configuration
type WebhookHandler struct {
	JobQueue      chan *workers.Job
	Platform      platforms.Platform
	WebhookSecret string
}

// NewWebhookHandler returns an instance of WebhookHandler
func NewWebhookHandler(platform platforms.Platform, webhookSecret string, jobQueue chan *workers.Job) *WebhookHandler {
	return &WebhookHandler{
		Platform:      platform,
		WebhookSecret: webhookSecret,
		JobQueue:      jobQueue,
	}
}

func (h *WebhookHandler) processEvent(owner, repository string) {
	log.Debug().Msgf("Creating new job for %s/%s", owner, repository)

	job, err := workers.NewJob(h.Platform, owner, repository)
	if err != nil {
		log.Error().Err(err).Msgf("Could create job for %s/%s", owner, repository)
		return
	}

	h.JobQueue <- job
}
