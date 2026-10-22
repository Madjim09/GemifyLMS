package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListQuizzes(c *gin.Context) {
	quizzes, err := h.quiz.ListQuizzes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: quizzes})
}

func (h *Handler) GetQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "некорректный ID теста"})
		return
	}

	quiz, err := h.quiz.GetQuizForClient(c.Request.Context(), quizID)
	if err != nil {
		c.JSON(http.StatusNotFound, response{Error: "тест не найден"})
		return
	}
	c.JSON(http.StatusOK, response{Data: quiz})
}

func (h *Handler) SubmitQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "некорректный ID теста"})
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response{Error: "требуется авторизация"})
		return
	}

	var body struct {
		Answers map[string]string `json:"answers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "укажите ответы на вопросы"})
		return
	}

	out, err := h.quiz.SubmitQuiz(c.Request.Context(), userID, quizID, body.Answers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: out})
}

func (h *Handler) GetUserQuizResults(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "некорректный ID пользователя"})
		return
	}

	results, err := h.quiz.GetUserResults(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: results})
}
