package usecase_test

import (
	"os"
	"testing"
	"time"

	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateToken_WithCustomSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-custom-secret-12345")
	defer os.Unsetenv("JWT_SECRET")

	userID := uuid.New()
	token, err := usecase.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	got, err := usecase.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if got != userID {
		t.Errorf("ValidateToken = %v, want %v", got, userID)
	}
}

func TestValidateToken_InvalidSubject(t *testing.T) {
	secret := "gamification-jwt-secret-change-in-production"
	// Create a token with non-UUID sub
	claims := jwt.MapClaims{
		"sub": 12345, // number instead of string
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := tok.SignedString([]byte(secret))

	_, err := usecase.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for non-string sub, got nil")
	}
}

func TestValidateToken_InvalidUUIDInSub(t *testing.T) {
	secret := "gamification-jwt-secret-change-in-production"
	claims := jwt.MapClaims{
		"sub": "not-a-valid-uuid",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := tok.SignedString([]byte(secret))

	_, err := usecase.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for invalid UUID in sub, got nil")
	}
}
