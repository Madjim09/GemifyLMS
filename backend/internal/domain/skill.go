package domain

import (
	"time"

	"github.com/google/uuid"
)

type SkillStatus string

const (
	SkillLocked     SkillStatus = "locked"
	SkillAvailable  SkillStatus = "available"
	SkillInProgress SkillStatus = "in_progress"
	SkillMastered   SkillStatus = "mastered"
)

type Skill struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Category    string     `json:"category,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	RequiredXP  int64      `json:"required_xp"`
	IconURL     string     `json:"icon_url,omitempty"`
	Order       int        `json:"order"`
	Children    []*Skill   `json:"children,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UserSkill struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	SkillID   uuid.UUID   `json:"skill_id"`
	Skill     *Skill      `json:"skill,omitempty"`
	Status    SkillStatus `json:"status"`
	Progress  int         `json:"progress"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type SkillRepository interface {
	GetAll(ctx interface{}) ([]*Skill, error)
	GetByID(ctx interface{}, id uuid.UUID) (*Skill, error)
	GetTree(ctx interface{}) ([]*Skill, error)
	GetUserSkills(ctx interface{}, userID uuid.UUID) ([]*UserSkill, error)
	UpsertUserSkill(ctx interface{}, userSkill *UserSkill) error
}
