package handler

import (
	"net/http"
	"strings"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleWebhookEvent(c *gin.Context) {
	var event domain.XAPIEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid JSON: " + err.Error()})
		return
	}

	if event.Actor.Mbox == "" && (event.Actor.Account == nil || event.Actor.Account.Name == "") {
		c.JSON(http.StatusBadRequest, response{Error: "actor identifier required"})
		return
	}

	externalID := event.Actor.Mbox
	if externalID == "" && event.Actor.Account != nil {
		externalID = event.Actor.Account.Name
	}
	externalID = strings.TrimPrefix(externalID, "mailto:")

	eventType, ok := usecase.MapXAPIVerbToEventType(event.Verb.ID)
	if !ok {
		eventType = domain.EventType(event.Verb.ID)
	}

	input := usecase.ProcessEventInput{
		XAPIEvent:  &event,
		EventType:  eventType,
		ExternalID: externalID,
		UserName:   event.Actor.Name,
		UserEmail:  strings.TrimPrefix(event.Actor.Mbox, "mailto:"),
		ObjectID:   event.Object.ID,
	}

	result, err := h.gamification.ProcessEvent(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{
		Message: "event processed",
		Data: gin.H{
			"user":             result.User,
			"xp_awarded":       result.XPAwarded,
			"leveled_up":       result.LeveledUp,
			"previous_level":   result.PreviousLevel,
			"new_achievements": result.NewAchievements,
		},
	})
}
