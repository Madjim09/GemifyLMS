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

func newGamificationUC(
	userRepo *testutil.MockUserRepository,
	achievementRepo *testutil.MockAchievementRepository,
	skillRepo *testutil.MockSkillRepository,
	eventRepo *testutil.MockEventRepository,
) *usecase.GamificationUseCase {
	return usecase.NewGamificationUseCase(
		userRepo, achievementRepo, skillRepo, eventRepo,
		cache.NewDisabledRedisClient(),
	)
}

func makeXAPIEvent(mbox, verbID, objectID string) *domain.XAPIEvent {
	return &domain.XAPIEvent{
		Actor:  domain.XAPIActor{Mbox: mbox, Name: "Test User"},
		Verb:   domain.XAPIVerb{ID: verbID},
		Object: domain.XAPIObject{ID: objectID},
	}
}

func TestGamificationUseCase_ProcessEvent_NewUser(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	event := makeXAPIEvent("mailto:new@example.com", "http://adlnet.gov/expapi/verbs/passed", "obj-1")
	input := usecase.ProcessEventInput{
		XAPIEvent:  event,
		EventType:  domain.EventTestPassed,
		ExternalID: "new@example.com",
		UserName:   "New User",
		UserEmail:  "new@example.com",
		ObjectID:   "obj-1",
	}

	out, err := uc.ProcessEvent(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.User == nil {
		t.Fatal("expected user in output")
	}
	if out.XPAwarded != domain.XPRewards[domain.EventTestPassed] {
		t.Errorf("XPAwarded = %d, want %d", out.XPAwarded, domain.XPRewards[domain.EventTestPassed])
	}
}

func TestGamificationUseCase_ProcessEvent_ExistingUser(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	existing := &domain.User{
		ID:         uuid.New(),
		ExternalID: "existing-ext-id",
		Email:      "existing@example.com",
		Name:       "Existing",
		XP:         400,
		Level:      1,
	}
	userRepo.Create(context.Background(), existing)

	event := makeXAPIEvent("mailto:existing@example.com", "http://adlnet.gov/expapi/verbs/passed", "obj-2")
	input := usecase.ProcessEventInput{
		XAPIEvent:  event,
		EventType:  domain.EventTestPassed,
		ExternalID: "existing-ext-id",
		UserName:   "Existing",
		UserEmail:  "existing@example.com",
		ObjectID:   "obj-2",
	}

	out, err := uc.ProcessEvent(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.User.ID != existing.ID {
		t.Errorf("expected same user ID")
	}
}

func TestGamificationUseCase_ProcessEvent_LevelUp(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	// XP = 450, EventTestPassed gives 100 XP → total 550 → level 2 (threshold 500)
	existing := &domain.User{
		ID:         uuid.New(),
		ExternalID: "levelup-user",
		Email:      "levelup@example.com",
		XP:         450,
		Level:      1,
	}
	userRepo.Create(context.Background(), existing)

	input := usecase.ProcessEventInput{
		XAPIEvent:  makeXAPIEvent("mailto:levelup@example.com", "http://adlnet.gov/expapi/verbs/passed", "obj"),
		EventType:  domain.EventTestPassed,
		ExternalID: "levelup-user",
		UserName:   "LevelUp",
		UserEmail:  "levelup@example.com",
		ObjectID:   "obj",
	}

	out, err := uc.ProcessEvent(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.LeveledUp {
		t.Error("expected LeveledUp = true")
	}
	if out.PreviousLevel != 1 {
		t.Errorf("PreviousLevel = %d, want 1", out.PreviousLevel)
	}
}

func TestGamificationUseCase_ProcessEvent_CreateUserError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	userRepo.CreateErr = fmt.Errorf("db error")
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	input := usecase.ProcessEventInput{
		XAPIEvent:  makeXAPIEvent("mailto:fail@example.com", "http://adlnet.gov/expapi/verbs/passed", "obj"),
		EventType:  domain.EventTestPassed,
		ExternalID: "fail-user",
	}

	_, err := uc.ProcessEvent(context.Background(), input)
	if err == nil {
		t.Fatal("expected error when create fails, got nil")
	}
}

func TestGamificationUseCase_ProcessEvent_SaveEventError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	eventRepo.SaveErr = fmt.Errorf("db error")
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	input := usecase.ProcessEventInput{
		XAPIEvent:  makeXAPIEvent("mailto:savefail@example.com", "http://adlnet.gov/expapi/verbs/passed", "obj"),
		EventType:  domain.EventTestPassed,
		ExternalID: "save-fail-user",
		UserEmail:  "savefail@example.com",
	}

	_, err := uc.ProcessEvent(context.Background(), input)
	if err == nil {
		t.Fatal("expected error when save fails, got nil")
	}
}

func TestGamificationUseCase_ProcessEvent_AddXPError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	existing := &domain.User{
		ID:         uuid.New(),
		ExternalID: "xp-fail-user",
		Email:      "xpfail@example.com",
		XP:         0,
		Level:      1,
	}
	userRepo.Create(context.Background(), existing)
	userRepo.AddXPErr = fmt.Errorf("db error")

	input := usecase.ProcessEventInput{
		XAPIEvent:  makeXAPIEvent("mailto:xpfail@example.com", "http://adlnet.gov/expapi/verbs/passed", "obj"),
		EventType:  domain.EventTestPassed,
		ExternalID: "xp-fail-user",
	}

	_, err := uc.ProcessEvent(context.Background(), input)
	if err == nil {
		t.Fatal("expected error when AddXP fails, got nil")
	}
}

// --- CheckAchievements tests ---

func TestCheckAchievements_TotalXP(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, Name: "Rich", ConditionType: domain.ConditionTotalXP, ConditionValue: 100, XPReward: 0},
	}

	user := &domain.User{ID: uuid.New(), Email: "xp@t.com", XP: 500}
	userRepo.Create(context.Background(), user)

	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if len(result) != 1 {
		t.Errorf("expected 1 achievement, got %d", len(result))
	}
}

func TestCheckAchievements_WithXPReward(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, Name: "Bonus", ConditionType: domain.ConditionTotalXP, ConditionValue: 50, XPReward: 200},
	}

	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "bonus@t.com", XP: 500}
	userRepo.Create(context.Background(), user)

	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if len(result) != 1 {
		t.Errorf("expected 1 achievement, got %d", len(result))
	}
}

func TestCheckAchievements_TestPerfect(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	userID := uuid.New()
	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionTestPerfect, ConditionValue: 1},
	}

	user := &domain.User{ID: userID, Email: "perfect@t.com", XP: 0}
	userRepo.Create(context.Background(), user)
	// Seed the count map directly
	eventRepo.CountMap[userID.String()+":"+string(domain.EventTestPerfect)] = 1

	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPerfect)
	if len(result) != 1 {
		t.Errorf("expected 1 achievement, got %d", len(result))
	}
}

func TestCheckAchievements_ModulesPerWeek(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	userID := uuid.New()
	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionModulesPerWeek, ConditionValue: 2},
	}

	user := &domain.User{ID: userID, Email: "weekly@t.com", XP: 0}
	userRepo.Create(context.Background(), user)
	eventRepo.CountMap[userID.String()+":"+string(domain.EventModuleCompleted)] = 3

	result := uc.CheckAchievements(context.Background(), user, domain.EventModuleCompleted)
	if len(result) != 1 {
		t.Errorf("expected 1 achievement, got %d", len(result))
	}
}

func TestCheckAchievements_VideoCount(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	userID := uuid.New()
	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionVideoCount, ConditionValue: 5},
	}

	user := &domain.User{ID: userID, Email: "video@t.com", XP: 0}
	userRepo.Create(context.Background(), user)
	eventRepo.CountMap[userID.String()+":"+string(domain.EventVideoWatched)] = 5

	result := uc.CheckAchievements(context.Background(), user, domain.EventVideoWatched)
	if len(result) != 1 {
		t.Errorf("expected 1 achievement, got %d", len(result))
	}
}

func TestCheckAchievements_TestCount(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	userID := uuid.New()
	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionTestCount, ConditionValue: 3},
	}

	user := &domain.User{ID: userID, Email: "tests@t.com", XP: 0}
	userRepo.Create(context.Background(), user)
	eventRepo.CountMap[userID.String()+":"+string(domain.EventTestPassed)] = 3

	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if len(result) != 1 {
		t.Errorf("expected 1 achievement, got %d", len(result))
	}
}

func TestCheckAchievements_AlreadyHas(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionTotalXP, ConditionValue: 0},
	}

	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "has@t.com", XP: 999}
	userRepo.Create(context.Background(), user)
	achievementRepo.HasResult[userID.String()+":"+achID.String()] = true

	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if len(result) != 0 {
		t.Errorf("expected 0 achievements (already has), got %d", len(result))
	}
}

func TestCheckAchievements_GetAllError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	achievementRepo.GetAllErr = fmt.Errorf("db error")
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	user := &domain.User{ID: uuid.New(), XP: 999}
	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if result != nil {
		t.Errorf("expected nil result on GetAll error, got %v", result)
	}
}

func TestCheckAchievements_AwardError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	achievementRepo.AwardToUserErr = fmt.Errorf("db error")
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionTotalXP, ConditionValue: 0},
	}
	user := &domain.User{ID: uuid.New(), XP: 100}
	userRepo.Create(context.Background(), user)

	// AwardToUser fails → achievement not in result but no panic
	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if len(result) != 0 {
		t.Errorf("expected 0 achievements when award fails, got %d", len(result))
	}
}

func TestCheckAchievements_ConditionNotMet(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newGamificationUC(userRepo, achievementRepo, skillRepo, eventRepo)

	achID := uuid.New()
	achievementRepo.Achievements = []*domain.Achievement{
		{ID: achID, ConditionType: domain.ConditionTotalXP, ConditionValue: 10000},
	}
	user := &domain.User{ID: uuid.New(), XP: 50}
	userRepo.Create(context.Background(), user)

	result := uc.CheckAchievements(context.Background(), user, domain.EventTestPassed)
	if len(result) != 0 {
		t.Errorf("expected 0 achievements (condition not met), got %d", len(result))
	}
}

// --- MapXAPIVerbToEventType tests ---

func TestMapXAPIVerbToEventType(t *testing.T) {
	tests := []struct {
		verbID    string
		eventType domain.EventType
		ok        bool
	}{
		{"http://adlnet.gov/expapi/verbs/passed", domain.EventTestPassed, true},
		{"http://adlnet.gov/expapi/verbs/completed", domain.EventModuleCompleted, true},
		{"http://adlnet.gov/expapi/verbs/experienced", domain.EventVideoWatched, true},
		{"http://adlnet.gov/expapi/verbs/watched", domain.EventVideoWatched, true},
		{"http://adlnet.gov/expapi/verbs/unknown", "", false},
		{"", "", false},
	}
	for _, tc := range tests {
		got, ok := usecase.MapXAPIVerbToEventType(tc.verbID)
		if ok != tc.ok {
			t.Errorf("MapXAPIVerbToEventType(%q) ok=%v, want %v", tc.verbID, ok, tc.ok)
			continue
		}
		if ok && got != tc.eventType {
			t.Errorf("MapXAPIVerbToEventType(%q) = %q, want %q", tc.verbID, got, tc.eventType)
		}
	}
}
