package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConditionType string

const (
	ConditionTestPerfect    ConditionType = "test_perfect"
	ConditionModulesPerWeek ConditionType = "modules_per_week"
	ConditionTotalXP        ConditionType = "total_xp"
	ConditionVideoCount     ConditionType = "video_count"
	ConditionTestCount      ConditionType = "test_count"
	ConditionStreakDays      ConditionType = "streak_days"
)

type Achievement struct {
	ID             uuid.UUID     `json:"id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	IconURL        string        `json:"icon_url"`
	ConditionType  ConditionType `json:"condition_type"`
	ConditionValue int           `json:"condition_value"`
	XPReward       int64         `json:"xp_reward"`
	CreatedAt      time.Time     `json:"created_at"`
}

type UserAchievement struct {
	ID            uuid.UUID    `json:"id"`
	UserID        uuid.UUID    `json:"user_id"`
	AchievementID uuid.UUID    `json:"achievement_id"`
	Achievement   *Achievement `json:"achievement,omitempty"`
	EarnedAt      time.Time    `json:"earned_at"`
}

type AchievementRepository interface {
	GetAll(ctx interface{}) ([]*Achievement, error)
	GetByID(ctx interface{}, id uuid.UUID) (*Achievement, error)
	GetUserAchievements(ctx interface{}, userID uuid.UUID) ([]*UserAchievement, error)
	AwardToUser(ctx interface{}, userID, achievementID uuid.UUID) (*UserAchievement, error)
	HasUserAchievement(ctx interface{}, userID, achievementID uuid.UUID) (bool, error)
}
