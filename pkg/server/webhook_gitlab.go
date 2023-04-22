package server

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"

	"github.com/fikaworks/grgate/pkg/utils"
)

var (
	gitlabEvents []gitlab.EventType = []gitlab.EventType{
		gitlab.EventTypeRelease,
		gitlab.EventTypePipeline,
	}
)

// GitlabHandler handle Gitlab webhook requests
func (h *WebhookHandler) GitlabHandler(c echo.Context) error {
	r := c.Request()
	defer func() {
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			log.Error().Err(err).Msg("Could discard request body")
		}
		if err := r.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Could not close request body")
		}
	}()

	signature := r.Header.Get("X-Gitlab-Token")
	if signature != h.WebhookSecret {
		log.Error().Msg("Token validation failed")
		return c.NoContent(http.StatusForbidden)
	}

	event := r.Header.Get("X-Gitlab-Event")
	if strings.TrimSpace(event) == "" {
		log.Error().Msg("Request is missing the X-Gitlab-Event header")
		return c.NoContent(http.StatusBadRequest)
	}

	eventType := gitlab.EventType(event)
	if !isGitlabEventSubscribed(eventType, gitlabEvents) {
		log.Error().Msgf("Event type %s is not supported", eventType)
		return c.NoContent(http.StatusBadRequest)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		log.Error().Msgf("Error reading request body from event type %s", eventType)
		return c.NoContent(http.StatusBadRequest)
	}

	parsedBody, err := gitlab.ParseWebhook(eventType, payload)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing request body from event type %s",
			eventType)
		return c.NoContent(http.StatusBadRequest)
	}

	log.Debug().Msgf("Received webhook event %s", event)

	switch eventType {
	case gitlab.EventTypeRelease:
		h.processGitlabReleaseEvent(*parsedBody.(*gitlab.ReleaseEvent))
	case gitlab.EventTypePipeline:
		h.processGitlabPipelineEvent(*parsedBody.(*gitlab.PipelineEvent))
	default:
		log.Info().Msgf("Event type %s is not supported", eventType)
	}

	return c.NoContent(http.StatusOK)
}

func (h *WebhookHandler) processGitlabReleaseEvent(event gitlab.ReleaseEvent) {
	owner := utils.GetRepositoryOrganization(event.Project.PathWithNamespace)
	repository := utils.GetRepositoryName(event.Project.PathWithNamespace)
	h.processEvent(owner, repository)
}

func (h *WebhookHandler) processGitlabPipelineEvent(event gitlab.PipelineEvent) {
	owner := utils.GetRepositoryOrganization(event.Project.PathWithNamespace)
	repository := utils.GetRepositoryName(event.Project.PathWithNamespace)
	h.processEvent(owner, repository)
}

func isGitlabEventSubscribed(event gitlab.EventType, events []gitlab.EventType) bool {
	for _, e := range events {
		if event == e {
			return true
		}
	}
	return false
}
