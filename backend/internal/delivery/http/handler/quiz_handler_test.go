package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corporate-training/gamification/internal/delivery/http/handler"
	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestListQuizzes_Success(t *testing.T) {
	deps := newTestDeps(t)
	deps.quizRepo.Quizzes[uuid.New()] = &domain.Quiz{ID: uuid.New(), Title: "Q1"}
	deps.quizRepo.Quizzes[uuid.New()] = &domain.Quiz{ID: uuid.New(), Title: "Q2"}

	r := gin.New()
	r.GET("/quizzes", deps.handler.ListQuizzes)

	req := httptest.NewRequest(http.MethodGet, "/quizzes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestListQuizzes_Error(t *testing.T) {
	deps := newTestDeps(t)
	deps.quizRepo.GetAllErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/quizzes", deps.handler.ListQuizzes)

	req := httptest.NewRequest(http.MethodGet, "/quizzes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetQuiz_Success(t *testing.T) {
	deps := newTestDeps(t)

	quizID := uuid.New()
	deps.quizRepo.Quizzes[quizID] = &domain.Quiz{ID: quizID, Title: "TestQuiz"}

	r := gin.New()
	r.GET("/quizzes/:id", deps.handler.GetQuiz)

	req := httptest.NewRequest(http.MethodGet, "/quizzes/"+quizID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetQuiz_InvalidID(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/quizzes/:id", deps.handler.GetQuiz)

	req := httptest.NewRequest(http.MethodGet, "/quizzes/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetQuiz_NotFound(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/quizzes/:id", deps.handler.GetQuiz)

	req := httptest.NewRequest(http.MethodGet, "/quizzes/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestSubmitQuiz_Success(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	deps.userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "submit@t.com", XP: 0, Level: 1})

	qID := uuid.New()
	optID := uuid.New()
	quizID := uuid.New()
	deps.quizRepo.Quizzes[quizID] = &domain.Quiz{
		ID:            quizID,
		Title:         "Test",
		XPReward:      100,
		PassThreshold: 70,
		Questions: []*domain.QuizQuestion{
			{
				ID: qID,
				Options: []*domain.QuizOption{
					{ID: optID, IsCorrect: true},
					{ID: uuid.New(), IsCorrect: false},
				},
			},
		},
	}

	token, _ := usecase.GenerateToken(userID)
	body := fmt.Sprintf(`{"answers":{"%s":"%s"}}`, qID.String(), optID.String())

	r := gin.New()
	r.POST("/quizzes/:id/submit", handler.AuthMiddleware(), deps.handler.SubmitQuiz)

	req := httptest.NewRequest(http.MethodPost, "/quizzes/"+quizID.String()+"/submit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestSubmitQuiz_InvalidQuizID(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	token, _ := usecase.GenerateToken(userID)

	r := gin.New()
	r.POST("/quizzes/:id/submit", handler.AuthMiddleware(), deps.handler.SubmitQuiz)

	req := httptest.NewRequest(http.MethodPost, "/quizzes/bad-id/submit", strings.NewReader(`{"answers":{}}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSubmitQuiz_NoAuth(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/quizzes/:id/submit", handler.AuthMiddleware(), deps.handler.SubmitQuiz)

	req := httptest.NewRequest(http.MethodPost, "/quizzes/"+uuid.New().String()+"/submit", strings.NewReader(`{"answers":{}}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSubmitQuiz_BadBody(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	token, _ := usecase.GenerateToken(userID)
	quizID := uuid.New()

	r := gin.New()
	r.POST("/quizzes/:id/submit", handler.AuthMiddleware(), deps.handler.SubmitQuiz)

	req := httptest.NewRequest(http.MethodPost, "/quizzes/"+quizID.String()+"/submit", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSubmitQuiz_QuizNotFound(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	deps.userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "notfound@t.com"})
	token, _ := usecase.GenerateToken(userID)

	r := gin.New()
	r.POST("/quizzes/:id/submit", handler.AuthMiddleware(), deps.handler.SubmitQuiz)

	body := `{"answers":{}}`
	req := httptest.NewRequest(http.MethodPost, "/quizzes/"+uuid.New().String()+"/submit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetUserQuizResults_Success(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	deps.quizRepo.UserResults[userID] = []*domain.UserQuizResult{
		{UserID: userID, Score: 80},
	}

	r := gin.New()
	r.GET("/users/:id/quiz-results", deps.handler.GetUserQuizResults)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/quiz-results", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetUserQuizResults_InvalidID(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id/quiz-results", deps.handler.GetUserQuizResults)

	req := httptest.NewRequest(http.MethodGet, "/users/bad-uuid/quiz-results", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetUserQuizResults_RepoError(t *testing.T) {
	deps := newTestDeps(t)
	deps.quizRepo.GetUserResultsErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/users/:id/quiz-results", deps.handler.GetUserQuizResults)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/quiz-results", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
