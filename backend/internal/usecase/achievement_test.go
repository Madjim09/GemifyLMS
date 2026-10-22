package usecase_test

import (
	"context"
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/google/uuid"
)

func TestAchievementUseCase_GetAll(t *testing.T) {
	achievementRepo := testutil.NewMockAchievementRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAchievementUseCase(achievementRepo, userRepo)

	achievementRepo.Achievements = []*domain.Achievement{
		{ID: uuid.New(), Name: "First Blood"},
		{ID: uuid.New(), Name: "Veteran"},
	}

	list, err := uc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}
}

func TestAchievementUseCase_GetAll_Empty(t *testing.T) {
	achievementRepo := testutil.NewMockAchievementRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAchievementUseCase(achievementRepo, userRepo)

	list, err := uc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestAchievementUseCase_GetUserAchievements(t *testing.T) {
	achievementRepo := testutil.NewMockAchievementRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAchievementUseCase(achievementRepo, userRepo)

	userID := uuid.New()
	ach := &domain.Achievement{ID: uuid.New(), Name: "Star"}
	achievementRepo.UserAchievements[userID] = []*domain.UserAchievement{
		{ID: uuid.New(), UserID: userID, AchievementID: ach.ID},
	}

	list, err := uc.GetUserAchievements(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("len = %d, want 1", len(list))
	}
}

func TestAchievementUseCase_GetUserAchievements_Empty(t *testing.T) {
	achievementRepo := testutil.NewMockAchievementRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewAchievementUseCase(achievementRepo, userRepo)

	list, err := uc.GetUserAchievements(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}
