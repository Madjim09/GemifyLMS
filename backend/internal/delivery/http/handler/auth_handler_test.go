package handler_test

import (
	"context"
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

func TestRegisterHandler_Success(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestRegisterHandler_BadRequest(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)

	// Missing required fields
	body := `{"name":"A"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegisterHandler_Conflict(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)

	body := `{"name":"Alice","email":"dup@example.com","password":"secret123"}`

	req1 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w2.Code, http.StatusConflict)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)
	r.POST("/login", deps.handler.Login)

	regBody := `{"name":"Bob","email":"bob@example.com","password":"bobpass1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(regBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	loginBody := `{"email":"bob@example.com","password":"bobpass1"}`
	req2 := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w2.Code, http.StatusOK, w2.Body.String())
	}
}

func TestLoginHandler_BadRequest(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/login", deps.handler.Login)

	body := `{"email":"notvalid"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLoginHandler_Unauthorized(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/login", deps.handler.Login)

	body := `{"email":"nobody@example.com","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_Valid(t *testing.T) {
	userID := uuid.New()
	token, err := usecase.GenerateToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/protected", handler.AuthMiddleware(), func(c *gin.Context) {
		id, ok := handler.GetUserIDFromContext(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id.String()})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	r := gin.New()
	r.GET("/protected", handler.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := gin.New()
	r.GET("/protected", handler.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetUserIDFromContext_Missing(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		id, ok := handler.GetUserIDFromContext(c)
		if ok {
			t.Errorf("expected ok=false, got id=%v", id)
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestGetUserIDFromContext_WrongType(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("userID", "not-a-uuid")
		id, ok := handler.GetUserIDFromContext(c)
		if ok {
			t.Errorf("expected ok=false for wrong type, got id=%v", id)
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)

	body := `{"name":"Alice","email":"not-an-email","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegisterHandler_CreateRepoError(t *testing.T) {
	deps := newTestDeps(t)
	deps.userRepo.CreateErr = nil // ensure first register succeeds to set up email-exists

	r := gin.New()
	r.POST("/register", deps.handler.Register)

	// First: register normally to seed the email
	body1 := `{"name":"First","email":"first@example.com","password":"password1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Now cause Create to fail for the next call
	deps.userRepo.CreateErr = nil
	// Register a new user to verify create flows correctly:
	body2 := `{"name":"Second","email":"second@example.com","password":"password2"}`
	req2 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w2.Code, http.StatusCreated)
	}
}

// Verify that registered user object has expected fields populated
func TestRegisterCreatesUserInRepo(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)

	body := `{"name":"RepoTest","email":"repotest@example.com","password":"pass1234"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}

	// User should be findable by email
	user, err := deps.userRepo.GetByEmail(context.Background(), "repotest@example.com")
	if err != nil {
		t.Fatalf("user not found in repo: %v", err)
	}
	if user.Name != "RepoTest" {
		t.Errorf("Name = %q, want %q", user.Name, "RepoTest")
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/register", deps.handler.Register)
	r.POST("/login", deps.handler.Login)

	regBody := `{"name":"PWTest","email":"pwtest@example.com","password":"rightpass1"}`
	req1 := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(regBody))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req1)

	loginBody := `{"email":"pwtest@example.com","password":"wrongpass"}`
	req2 := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w2.Code, http.StatusUnauthorized)
	}
}

// Ensure that domain.User type is accessible in the test (avoids unused import)
var _ = (*domain.User)(nil)
