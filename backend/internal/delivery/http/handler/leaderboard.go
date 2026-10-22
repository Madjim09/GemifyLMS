package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLeaderboard(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	entries, err := h.leaderboard.GetTopUsers(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Data: entries})
}
