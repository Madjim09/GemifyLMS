package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	user, err := h.user.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, response{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, response{Data: user})
}

func (h *Handler) ListUsers(c *gin.Context) {
	offset := 0
	limit := 20
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	users, total, err := h.user.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   users,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}
