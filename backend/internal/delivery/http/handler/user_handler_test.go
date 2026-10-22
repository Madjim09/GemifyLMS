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

func TestGetUser_Success(t *testing.T) {
	deps := newTestDeps(t)

	userID := uuid.New()
	deps.userRepo.Create(context.Background(), &domain.User{
		ID: userID, Name: "Frank", Email: "frank@example.com",
	})

	r := gin.New()
	r.GET("/users/:id", deps.handler.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id", deps.handler.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/users/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users/:id", deps.handler.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestListUsers_Success(t *testing.T) {
	deps := newTestDeps(t)

	for i := 0; i < 3; i++ {
		deps.userRepo.Create(context.Background(), &domain.User{
			ID: uuid.New(), Email: fmt.Sprintf("list%d@example.com", i),
		})
	}

	r := gin.New()
	r.GET("/users", deps.handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestListUsers_WithPagination(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users", deps.handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/users?offset=0&limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestListUsers_InvalidPagination(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users", deps.handler.ListUsers)

	// Invalid values should be ignored and defaults applied
	req := httptest.NewRequest(http.MethodGet, "/users?offset=abc&limit=xyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestListUsers_LimitOutOfRange(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.GET("/users", deps.handler.ListUsers)

	// limit=0 and limit=101 are out of range, default should apply
	for _, limit := range []string{"0", "101", "-5"} {
		req := httptest.NewRequest(http.MethodGet, "/users?limit="+limit, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("limit=%s: status = %d, want %d", limit, w.Code, http.StatusOK)
		}
	}
}

func TestListUsers_RepoError(t *testing.T) {
	deps := newTestDeps(t)
	deps.userRepo.ListErr = fmt.Errorf("db error")

	r := gin.New()
	r.GET("/users", deps.handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
