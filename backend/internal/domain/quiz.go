package domain

import (
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID            uuid.UUID       `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Difficulty    string          `json:"difficulty"` // easy, medium, hard
	XPReward      int64           `json:"xp_reward"`
	BonusXP       int64           `json:"bonus_xp"`
	PassThreshold int             `json:"pass_threshold"`
	Questions     []*QuizQuestion `json:"questions,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type QuizQuestion struct {
	ID      uuid.UUID     `json:"id"`
	QuizID  uuid.UUID     `json:"quiz_id"`
	Text    string        `json:"text"`
	Order   int           `json:"order"`
	Options []*QuizOption `json:"options,omitempty"`
}

type QuizOption struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Text       string    `json:"text"`
	IsCorrect  bool      `json:"is_correct,omitempty"`
	Order      int       `json:"order"`
}

type UserQuizResult struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	QuizID      uuid.UUID `json:"quiz_id"`
	Quiz        *Quiz     `json:"quiz,omitempty"`
	Score       int       `json:"score"`
	XPEarned    int64     `json:"xp_earned"`
	IsPerfect   bool      `json:"is_perfect"`
	CompletedAt time.Time `json:"completed_at"`
}

type QuizRepository interface {
	GetAll(ctx interface{}) ([]*Quiz, error)
	GetByIDWithQuestions(ctx interface{}, id uuid.UUID) (*Quiz, error)
	SaveResult(ctx interface{}, result *UserQuizResult) error
	GetUserResults(ctx interface{}, userID uuid.UUID) ([]*UserQuizResult, error)
}
