package server

import (
	"net/http"

	"github.com/google/go-github/v34/github"
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/ggate/pkg/platforms"
	"github.com/fikaworks/ggate/pkg/workers"
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

// GithubHandler handle Github webhook requests
func (h *WebhookHandler) GithubHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(h.WebhookSecret))
	if err != nil {
		log.Error().Err(err).Msg("Error validating request body")
		return
	}

	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Error().Err(err).Msg("Could not parse webhook")
		return
	}

	switch e := event.(type) {
	case *github.StatusEvent:
		log.Debug().Msgf("Received webhook event %s for %s/%s",
			"StatusEvent", *e.Repo.Owner.Login, *e.Repo.Name)
		h.processStatusEvent(*event.(*github.StatusEvent))
	case *github.CheckSuiteEvent:
		log.Debug().Msgf("Received webhook event %s for %s/%s",
			"CheckSuiteEvent", *e.Repo.Owner.Login, *e.Repo.Name)
		h.processCheckSuiteEvent(*event.(*github.CheckSuiteEvent))
	case *github.CheckRunEvent:
		log.Debug().Msgf("Received webhook event %s for %s/%s",
			"CheckRunEvent", *e.Repo.Owner.Login, *e.Repo.Name)
		h.processCheckRunEvent(*event.(*github.CheckRunEvent))

	default:
		log.Info().Msgf("Unknown event type %s", github.WebHookType(r))
		return
	}
}

func (h *WebhookHandler) processStatusEvent(event github.StatusEvent) {
	if event.State != nil && *event.State == "success" {
		log.Debug().Msgf("Creating new job for %s/%s", *event.Repo.Owner.Login,
			*event.Repo.Name)

		job, err := workers.NewJob(h.Platform, *event.Repo.Owner.Login,
			*event.Repo.Name)
		if err != nil {
			log.Error().Err(err).Msgf("Could create job for %s/%s from event %s",
				*event.Repo.Owner.Login, *event.Repo.Name, *event.State)
		}

		h.JobQueue <- job
	}
}

func (h *WebhookHandler) processCheckSuiteEvent(event github.CheckSuiteEvent) {
	if event.Action != nil && *event.Action == "completed" {
		log.Debug().Msgf("Creating new job for %s/%s", *event.Repo.Owner.Login,
			*event.Repo.Name)

		job, err := workers.NewJob(h.Platform, *event.Repo.Owner.Login,
			*event.Repo.Name)
		if err != nil {
			log.Error().Err(err).Msgf("Could create job for %s/%s from event %s",
				*event.Repo.Owner.Login, *event.Repo.Name, *event.Action)
		}

		h.JobQueue <- job
	}
}

func (h *WebhookHandler) processCheckRunEvent(event github.CheckRunEvent) {
	if event.Action != nil && *event.Action == "completed" {
		log.Debug().Msgf("Creating new job for %s/%s", *event.Repo.Owner.Login,
			*event.Repo.Name)

		job, err := workers.NewJob(h.Platform, *event.Repo.Owner.Login,
			*event.Repo.Name)
		if err != nil {
			log.Error().Err(err).Msgf("Could create job for %s/%s from event %s",
				*event.Repo.Owner.Login, *event.Repo.Name, *event.Action)
		}

		h.JobQueue <- job
	}
}
