package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/google/uuid"
)

func TestSkillUseCase_GetSkillTree(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	skillRepo.Skills = []*domain.Skill{
		{ID: uuid.New(), Name: "Go"},
		{ID: uuid.New(), Name: "Docker"},
	}

	tree, err := uc.GetSkillTree(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tree) != 2 {
		t.Errorf("len(tree) = %d, want 2", len(tree))
	}
}

func TestSkillUseCase_GetSkillTree_Error(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	skillRepo.GetTreeErr = fmt.Errorf("db error")
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	_, err := uc.GetSkillTree(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSkillUseCase_GetUserSkills_UserNotFound(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	_, err := uc.GetUserSkills(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
}

func TestSkillUseCase_GetUserSkills_GetUserSkillsError(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	skillRepo.GetUserSkillsErr = fmt.Errorf("db error")
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "s@s.com", XP: 0})

	_, err := uc.GetUserSkills(context.Background(), userID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSkillUseCase_GetUserSkills_GetAllError(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	skillRepo.GetAllErr = fmt.Errorf("db error")
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "s2@s.com", XP: 0})

	_, err := uc.GetUserSkills(context.Background(), userID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSkillUseCase_GetUserSkills_LockedAndAvailable(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "skills@s.com", XP: 600})

	skillRepo.Skills = []*domain.Skill{
		{ID: uuid.New(), Name: "Go", RequiredXP: 0},
		{ID: uuid.New(), Name: "Docker", RequiredXP: 500},
		{ID: uuid.New(), Name: "Kubernetes", RequiredXP: 1000},
	}

	skills, err := uc.GetUserSkills(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 3 {
		t.Fatalf("expected 3 skill entries, got %d", len(skills))
	}

	available := 0
	locked := 0
	for _, s := range skills {
		switch s.Status {
		case domain.SkillAvailable:
			available++
		case domain.SkillLocked:
			locked++
		}
	}
	if available != 2 {
		t.Errorf("available = %d, want 2", available)
	}
	if locked != 1 {
		t.Errorf("locked = %d, want 1", locked)
	}
}

func TestSkillUseCase_GetUserSkills_WithExistingUserSkill(t *testing.T) {
	skillRepo := testutil.NewMockSkillRepository()
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewSkillUseCase(skillRepo, userRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "existing@s.com", XP: 100})

	skillID := uuid.New()
	skillRepo.Skills = []*domain.Skill{
		{ID: skillID, Name: "Go", RequiredXP: 0},
	}
	skillRepo.UserSkills[userID] = []*domain.UserSkill{
		{UserID: userID, SkillID: skillID, Status: domain.SkillMastered, Progress: 100},
	}

	skills, err := uc.GetUserSkills(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Status != domain.SkillMastered {
		t.Errorf("status = %q, want %q", skills[0].Status, domain.SkillMastered)
	}
}
