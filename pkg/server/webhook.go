package server

import (
	"net/http"

	"github.com/google/go-github/v34/github"
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/ggate/pkg/platforms"
	"github.com/fikaworks/ggate/pkg/workers"
)

type WebhookHandler struct {
  WebhookSecret string
  Platform platforms.Platform
}

func NewWebhookHandler(platform platforms.Platform, webhookSecret string) *WebhookHandler {
  return &WebhookHandler{
    WebhookSecret: webhookSecret,
    Platform: platform,
  }
}

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
	case *github.CheckRunEvent:
    // https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads#check_run
    if e.Action != nil && *e.Action == "completed" {
      job, err := workers.NewJob(h.Platform, *e.Repo.Owner.Login, *e.Repo.Name)
      if err != nil {
		    log.Error().Err(err).Msgf("Could create job for %s/%s from event %s",
           *e.Repo.Owner.Login, *e.Repo.Name, *e.Action)
      }
      job.Process()
    }
	default:
		log.Info().Msgf("Unknown event type %s", github.WebHookType(r))
		return
	}
}
