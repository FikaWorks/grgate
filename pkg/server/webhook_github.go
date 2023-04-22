package server

import (
	"net/http"

	"github.com/google/go-github/v43/github"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *WebhookHandler) GithubHandler(c echo.Context) error {
	r := c.Request()
	payload, err := github.ValidatePayload(r, []byte(h.WebhookSecret))
	if err != nil {
		log.Error().Err(err).Msg("Error validating request body")
		return c.NoContent(http.StatusForbidden)
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Error().Err(err).Msg("Could not parse webhook")
		return c.NoContent(http.StatusBadRequest)
	}

	switch event := event.(type) {
	case *github.StatusEvent:
		h.processGithubStatusEvent(event)
	case *github.CheckSuiteEvent:
		h.processGithubCheckSuiteEvent(event)
	case *github.CheckRunEvent:
		h.processGithubCheckRunEvent(event)
	case *github.ReleaseEvent:
		h.processGithubReleaseEvent(event)
	default:
		log.Info().Msgf("Event type %s is not supported", github.WebHookType(r))
	}

	return c.NoContent(http.StatusOK)
}

func (h *WebhookHandler) processGithubStatusEvent(event *github.StatusEvent) {
	log.Debug().Msg("Received webhook event StatusEvent")
	if event.State != nil && *event.State == "success" {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}

func (h *WebhookHandler) processGithubCheckSuiteEvent(event *github.CheckSuiteEvent) {
	log.Debug().Msg("Received webhook event CheckSuiteEvent")
	if event.Action != nil && *event.Action == "completed" {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}

func (h *WebhookHandler) processGithubCheckRunEvent(event *github.CheckRunEvent) {
	log.Debug().Msg("Received webhook event CheckRunEvent")
	if event.Action != nil && *event.Action == "completed" {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}

func (h *WebhookHandler) processGithubReleaseEvent(event *github.ReleaseEvent) {
	log.Debug().Msg("Received webhook event ReleaseEvent")
	if event.Action != nil && (*event.Action == "created" || *event.Action == "edited") {
		h.processEvent(*event.Repo.Owner.Login, *event.Repo.Name)
	}
}
