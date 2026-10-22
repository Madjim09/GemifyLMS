package usecase

import (
	"context"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error) {
	return uc.userRepo.List(ctx, offset, limit)
}
