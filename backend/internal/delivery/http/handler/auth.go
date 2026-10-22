package handler

import (
	"net/http"
	"strings"

	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "Проверьте данные: имя (мин. 2 символа), email и пароль (мин. 6 символов)"})
		return
	}

	out, err := h.auth.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: out})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "Укажите email и пароль"})
		return
	}

	out, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: out})
}

// AuthMiddleware validates JWT and sets userID in context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response{Error: "Требуется авторизация"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := usecase.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response{Error: "Недействительный токен"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserIDFromContext extracts userID set by AuthMiddleware.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
