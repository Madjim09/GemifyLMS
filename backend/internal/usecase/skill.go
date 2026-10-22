package usecase

import (
	"context"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type SkillUseCase struct {
	skillRepo domain.SkillRepository
	userRepo  domain.UserRepository
}

func NewSkillUseCase(skillRepo domain.SkillRepository, userRepo domain.UserRepository) *SkillUseCase {
	return &SkillUseCase{skillRepo: skillRepo, userRepo: userRepo}
}

func (uc *SkillUseCase) GetSkillTree(ctx context.Context) ([]*domain.Skill, error) {
	return uc.skillRepo.GetTree(ctx)
}

func (uc *SkillUseCase) GetUserSkills(ctx context.Context, userID uuid.UUID) ([]*domain.UserSkill, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userSkills, err := uc.skillRepo.GetUserSkills(ctx, userID)
	if err != nil {
		return nil, err
	}

	allSkills, err := uc.skillRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	userSkillMap := make(map[uuid.UUID]*domain.UserSkill)
	for _, us := range userSkills {
		userSkillMap[us.SkillID] = us
	}

	var result []*domain.UserSkill
	for _, skill := range allSkills {
		if us, exists := userSkillMap[skill.ID]; exists {
			result = append(result, us)
		} else {
			status := domain.SkillLocked
			if user.XP >= skill.RequiredXP {
				status = domain.SkillAvailable
			}
			result = append(result, &domain.UserSkill{
				UserID:   userID,
				SkillID:  skill.ID,
				Skill:    skill,
				Status:   status,
				Progress: 0,
			})
		}
	}

	return result, nil
}
