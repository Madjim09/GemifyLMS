package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestRegister_Success(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAuthUseCase(userRepo)

	out, err := uc.Register(context.Background(), "Alice", "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.User == nil || out.Token == "" {
		t.Fatal("expected user and token in output")
	}
	if out.User.Email != "alice@example.com" {
		t.Errorf("email = %q, want %q", out.User.Email, "alice@example.com")
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAuthUseCase(userRepo)

	_, err := uc.Register(context.Background(), "Alice", "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = uc.Register(context.Background(), "Alice2", "alice@example.com", "other123")
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestRegister_CreateError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	userRepo.CreateErr = fmt.Errorf("db error")
	uc := usecase.NewAuthUseCase(userRepo)

	_, err := uc.Register(context.Background(), "Bob", "bob@example.com", "pass123")
	if err == nil {
		t.Fatal("expected error when create fails, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAuthUseCase(userRepo)

	_, err := uc.Register(context.Background(), "Carol", "carol@example.com", "mypassword")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	out, err := uc.Login(context.Background(), "carol@example.com", "mypassword")
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
	if out.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAuthUseCase(userRepo)

	_, err := uc.Login(context.Background(), "nobody@example.com", "pass")
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAuthUseCase(userRepo)

	_, err := uc.Register(context.Background(), "Dave", "dave@example.com", "correctpass")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err = uc.Login(context.Background(), "dave@example.com", "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	userID := uuid.New()

	token, err := usecase.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	got, err := usecase.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if got != userID {
		t.Errorf("ValidateToken = %v, want %v", got, userID)
	}
}

func TestValidateToken_InvalidString(t *testing.T) {
	_, err := usecase.ValidateToken("not.a.valid.jwt")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "gamification-jwt-secret-change-in-production"

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := tok.SignedString([]byte(secret))

	_, err := usecase.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	_, err := usecase.ValidateToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.invalid")
	if err == nil {
		t.Fatal("expected error for wrong signing method, got nil")
	}
}
