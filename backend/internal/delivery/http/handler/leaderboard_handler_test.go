package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetLeaderboard_Success(t *testing.T) {
	deps := newTestDeps(t)

	for i := 0; i < 3; i++ {
		deps.userRepo.Create(context.Background(), &domain.User{
			ID:    uuid.New(),
			Email: fmt.Sprintf("lb%d@example.com", i),
			XP:    int64(i * 100),
			Level: 1,
		})
	}

	r := gin.New()
	r.GET("/leaderboard", deps.handler.GetLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetLeaderboard_WithLimit(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/leaderboard", deps.handler.GetLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/leaderboard?limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetLeaderboard_InvalidLimit(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/leaderboard", deps.handler.GetLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/leaderboard?limit=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetLeaderboard_Error(t *testing.T) {
	deps := newTestDeps(t)
	deps.userRepo.GetTopByXPErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/leaderboard", deps.handler.GetLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
