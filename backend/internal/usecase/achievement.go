package usecase

import (
	"context"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type AchievementUseCase struct {
	achievementRepo domain.AchievementRepository
	userRepo        domain.UserRepository
}

func NewAchievementUseCase(achievementRepo domain.AchievementRepository, userRepo domain.UserRepository) *AchievementUseCase {
	return &AchievementUseCase{achievementRepo: achievementRepo, userRepo: userRepo}
}

func (uc *AchievementUseCase) GetAll(ctx context.Context) ([]*domain.Achievement, error) {
	return uc.achievementRepo.GetAll(ctx)
}

func (uc *AchievementUseCase) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]*domain.UserAchievement, error) {
	return uc.achievementRepo.GetUserAchievements(ctx, userID)
}
