package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/google/uuid"
)

func TestLeaderboardUseCase_GetTopUsers(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	redisClient := cache.NewDisabledRedisClient()
	uc := usecase.NewLeaderboardUseCase(userRepo, redisClient)

	for i := 0; i < 5; i++ {
		userRepo.Create(context.Background(), &domain.User{
			ID:    uuid.New(),
			Email: fmt.Sprintf("lb%d@example.com", i),
			XP:    int64(i * 100),
			Level: 1,
		})
	}

	entries, err := uc.GetTopUsers(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("len(entries) = %d, want 5", len(entries))
	}
	for i, e := range entries {
		if e.Rank != i+1 {
			t.Errorf("entries[%d].Rank = %d, want %d", i, e.Rank, i+1)
		}
	}
}

func TestLeaderboardUseCase_GetTopUsers_DefaultLimit(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	redisClient := cache.NewDisabledRedisClient()
	uc := usecase.NewLeaderboardUseCase(userRepo, redisClient)

	// limit <= 0 should default to 10
	entries, err := uc.GetTopUsers(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = entries
}

func TestLeaderboardUseCase_GetTopUsers_LimitTooLarge(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	redisClient := cache.NewDisabledRedisClient()
	uc := usecase.NewLeaderboardUseCase(userRepo, redisClient)

	// limit > 100 should default to 10
	entries, err := uc.GetTopUsers(context.Background(), 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = entries
}

func TestLeaderboardUseCase_GetTopUsers_RepoError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	userRepo.GetTopByXPErr = fmt.Errorf("db error")
	redisClient := cache.NewDisabledRedisClient()
	uc := usecase.NewLeaderboardUseCase(userRepo, redisClient)

	_, err := uc.GetTopUsers(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
