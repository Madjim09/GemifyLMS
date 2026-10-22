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

func TestGetSkillTree_Success(t *testing.T) {
	deps := newTestDeps(t)
	deps.skillRepo.Skills = []*domain.Skill{
		{ID: uuid.New(), Name: "Go"},
		{ID: uuid.New(), Name: "Docker"},
	}

	r := gin.New()
	r.GET("/skills", deps.handler.GetSkillTree)

	req := httptest.NewRequest(http.MethodGet, "/skills", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetSkillTree_Error(t *testing.T) {
	deps := newTestDeps(t)
	deps.skillRepo.GetTreeErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/skills", deps.handler.GetSkillTree)

	req := httptest.NewRequest(http.MethodGet, "/skills", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetUserSkills_Success(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	deps.userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "skills@test.com", XP: 100})
	deps.skillRepo.Skills = []*domain.Skill{
		{ID: uuid.New(), Name: "Go", RequiredXP: 0},
	}

	r := gin.New()
	r.GET("/users/:id/skills", deps.handler.GetUserSkills)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/skills", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetUserSkills_InvalidID(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id/skills", deps.handler.GetUserSkills)

	req := httptest.NewRequest(http.MethodGet, "/users/not-valid/skills", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetUserSkills_UserNotFound(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id/skills", deps.handler.GetUserSkills)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/skills", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
