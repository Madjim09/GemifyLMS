package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
	"github.com/google/uuid"
)

type GamificationUseCase struct {
	userRepo        domain.UserRepository
	achievementRepo domain.AchievementRepository
	skillRepo       domain.SkillRepository
	eventRepo       domain.EventRepository
	cache           *cache.RedisClient
}

func NewGamificationUseCase(
	userRepo domain.UserRepository,
	achievementRepo domain.AchievementRepository,
	skillRepo domain.SkillRepository,
	eventRepo domain.EventRepository,
	cache *cache.RedisClient,
) *GamificationUseCase {
	return &GamificationUseCase{
		userRepo:        userRepo,
		achievementRepo: achievementRepo,
		skillRepo:       skillRepo,
		eventRepo:       eventRepo,
		cache:           cache,
	}
}

type ProcessEventInput struct {
	XAPIEvent  *domain.XAPIEvent
	EventType  domain.EventType
	ExternalID string
	UserName   string
	UserEmail  string
	ObjectID   string
}

type ProcessEventOutput struct {
	User            *domain.User
	XPAwarded       int64
	LeveledUp       bool
	PreviousLevel   int
	NewAchievements []*domain.Achievement
}

func (uc *GamificationUseCase) ProcessEvent(ctx context.Context, input ProcessEventInput) (*ProcessEventOutput, error) {
	user, err := uc.userRepo.GetByExternalID(ctx, input.ExternalID)
	if err != nil {
		user = &domain.User{
			ID:         uuid.New(),
			ExternalID: input.ExternalID,
			Name:       input.UserName,
			Email:      input.UserEmail,
			XP:         0,
			Level:      1,
		}
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	xpAmount := domain.XPRewards[input.EventType]

	rawPayload, _ := json.Marshal(input.XAPIEvent)
	eventRecord := &domain.EventRecord{
		ID:         uuid.New(),
		UserID:     user.ID,
		EventType:  input.EventType,
		ObjectID:   input.ObjectID,
		XPAwarded:  xpAmount,
		RawPayload: rawPayload,
		CreatedAt:  time.Now(),
	}

	if err := uc.eventRepo.Save(ctx, eventRecord); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}

	previousLevel := user.Level
	updatedUser, err := uc.userRepo.AddXP(ctx, user.ID, xpAmount)
	if err != nil {
		return nil, fmt.Errorf("add xp: %w", err)
	}

	uc.cache.InvalidateLeaderboard(ctx)

	newAchievements := uc.checkAndAwardAchievements(ctx, updatedUser, input.EventType)

	return &ProcessEventOutput{
		User:            updatedUser,
		XPAwarded:       xpAmount,
		LeveledUp:       updatedUser.Level > previousLevel,
		PreviousLevel:   previousLevel,
		NewAchievements: newAchievements,
	}, nil
}

func (uc *GamificationUseCase) CheckAchievements(ctx context.Context, user *domain.User, eventType domain.EventType) []*domain.Achievement {
	return uc.checkAndAwardAchievements(ctx, user, eventType)
}

func (uc *GamificationUseCase) checkAndAwardAchievements(ctx context.Context, user *domain.User, eventType domain.EventType) []*domain.Achievement {
	achievements, err := uc.achievementRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get achievements: %v", err)
		return nil
	}

	var newAchievements []*domain.Achievement
	weekAgo := time.Now().AddDate(0, 0, -7)

	for _, ach := range achievements {
		has, _ := uc.achievementRepo.HasUserAchievement(ctx, user.ID, ach.ID)
		if has {
			continue
		}

		earned := false
		switch ach.ConditionType {
		case domain.ConditionTotalXP:
			earned = user.XP >= int64(ach.ConditionValue)
		case domain.ConditionTestPerfect:
			count, _ := uc.eventRepo.GetUserEventCount(ctx, user.ID, domain.EventTestPerfect)
			earned = count >= ach.ConditionValue
		case domain.ConditionModulesPerWeek:
			count, _ := uc.eventRepo.GetUserEventCountSince(ctx, user.ID, domain.EventModuleCompleted, weekAgo)
			earned = count >= ach.ConditionValue
		case domain.ConditionVideoCount:
			count, _ := uc.eventRepo.GetUserEventCount(ctx, user.ID, domain.EventVideoWatched)
			earned = count >= ach.ConditionValue
		case domain.ConditionTestCount:
			count, _ := uc.eventRepo.GetUserEventCount(ctx, user.ID, domain.EventTestPassed)
			earned = count >= ach.ConditionValue
		}

		if earned {
			ua, err := uc.achievementRepo.AwardToUser(ctx, user.ID, ach.ID)
			if err != nil {
				log.Printf("Failed to award achievement %s: %v", ach.ID, err)
				continue
			}
			_ = ua
			if ach.XPReward > 0 {
				uc.userRepo.AddXP(ctx, user.ID, ach.XPReward)
			}
			newAchievements = append(newAchievements, ach)
		}
	}

	return newAchievements
}

func MapXAPIVerbToEventType(verbID string) (domain.EventType, bool) {
	verbMap := map[string]domain.EventType{
		"http://adlnet.gov/expapi/verbs/passed":      domain.EventTestPassed,
		"http://adlnet.gov/expapi/verbs/completed":   domain.EventModuleCompleted,
		"http://adlnet.gov/expapi/verbs/experienced": domain.EventVideoWatched,
		"http://adlnet.gov/expapi/verbs/watched":     domain.EventVideoWatched,
	}

	for k, v := range verbMap {
		if strings.EqualFold(k, verbID) {
			return v, true
		}
	}
	return "", false
}
