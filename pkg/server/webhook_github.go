package server

import (
	"net/http"

	"github.com/google/go-github/v34/github"
	"github.com/rs/zerolog/log"
)

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

	switch event.(type) {
	case *github.StatusEvent:
		h.processGithubStatusEvent(event.(github.StatusEvent))
	case *github.CheckSuiteEvent:
		h.processGithubCheckSuiteEvent(event.(github.CheckSuiteEvent))
	case *github.CheckRunEvent:
		h.processGithubCheckRunEvent(event.(github.CheckRunEvent))
	case *github.ReleaseEvent:
		h.processGithubReleaseEvent(event.(github.ReleaseEvent))
	default:
		log.Info().Msgf("Event type %s is not supported", github.WebHookType(r))
	}
}

func (h *WebhookHandler) processGithubStatusEvent(event github.StatusEvent) {
	log.Trace().Msg("Received webhook event StatusEvent")
	if event.State != nil && *event.State == "success" {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}

func (h *WebhookHandler) processGithubCheckSuiteEvent(event github.CheckSuiteEvent) {
	log.Trace().Msg("Received webhook event CheckSuiteEvent")
	if event.Action != nil && *event.Action == "completed" {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}

func (h *WebhookHandler) processGithubCheckRunEvent(event github.CheckRunEvent) {
	log.Trace().Msg("Received webhook event CheckRunEvent")
	if event.Action != nil && *event.Action == "completed" {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}

func (h *WebhookHandler) processGithubReleaseEvent(event github.ReleaseEvent) {
	log.Trace().Msg("Received webhook event ReleaseEvent")
	if event.Action != nil && (*event.Action == "created" || *event.Action == "edited") {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}
