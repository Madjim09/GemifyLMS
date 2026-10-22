package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	ExternalID   string    `json:"external_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	XP           int64     `json:"xp"`
	Level        int       `json:"level"`
	TestsPassed  int       `json:"tests_passed"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) RecalculateLevel() {
	u.Level = CalculateLevel(u.XP)
}

func CalculateLevel(xp int64) int {
	thresholds := []int64{0, 500, 1500, 3000, 5500, 9000, 14000, 21000, 30000, 42000}
	level := 1
	for i, threshold := range thresholds {
		if xp >= threshold {
			level = i + 1
		}
	}
	return level
}

func XPForNextLevel(currentLevel int) int64 {
	thresholds := []int64{0, 500, 1500, 3000, 5500, 9000, 14000, 21000, 30000, 42000}
	if currentLevel >= len(thresholds) {
		return -1
	}
	return thresholds[currentLevel]
}

type UserRepository interface {
	Create(ctx interface{}, user *User) error
	GetByID(ctx interface{}, id uuid.UUID) (*User, error)
	GetByExternalID(ctx interface{}, externalID string) (*User, error)
	GetByEmail(ctx interface{}, email string) (*User, error)
	Update(ctx interface{}, user *User) error
	AddXP(ctx interface{}, userID uuid.UUID, amount int64) (*User, error)
	IncrementTestsPassed(ctx interface{}, userID uuid.UUID) error
	GetTopByXP(ctx interface{}, limit int) ([]*User, error)
	List(ctx interface{}, offset, limit int) ([]*User, int64, error)
}
