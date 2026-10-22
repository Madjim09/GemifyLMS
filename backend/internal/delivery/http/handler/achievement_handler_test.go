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

func TestGetAchievements_Success(t *testing.T) {
	deps := newTestDeps(t)
	deps.achievementRepo.Achievements = []*domain.Achievement{
		{ID: uuid.New(), Name: "Pioneer"},
		{ID: uuid.New(), Name: "Expert"},
	}

	r := gin.New()
	r.GET("/achievements", deps.handler.GetAchievements)

	req := httptest.NewRequest(http.MethodGet, "/achievements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetAchievements_Error(t *testing.T) {
	deps := newTestDeps(t)
	deps.achievementRepo.GetAllErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/achievements", deps.handler.GetAchievements)

	req := httptest.NewRequest(http.MethodGet, "/achievements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetUserAchievements_Success(t *testing.T) {
	deps := newTestDeps(t)
	userID := uuid.New()
	achID := uuid.New()
	deps.achievementRepo.UserAchievements[userID] = []*domain.UserAchievement{
		{ID: uuid.New(), UserID: userID, AchievementID: achID},
	}

	r := gin.New()
	r.GET("/users/:id/achievements", deps.handler.GetUserAchievements)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/achievements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetUserAchievements_InvalidID(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id/achievements", deps.handler.GetUserAchievements)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid-id/achievements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetUserAchievements_Empty(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id/achievements", deps.handler.GetUserAchievements)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/achievements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetUserAchievements_RepoError(t *testing.T) {
	deps := newTestDeps(t)
	deps.achievementRepo.GetUserAchievementsErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/users/:id/achievements", deps.handler.GetUserAchievements)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/achievements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// Ensure context import used
var _ = context.Background
