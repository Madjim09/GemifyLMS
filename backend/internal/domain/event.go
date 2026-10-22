package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTestPassed      EventType = "test_passed"
	EventTestPerfect     EventType = "test_perfect"
	EventVideoWatched    EventType = "video_watched"
	EventModuleCompleted EventType = "module_completed"
	EventCourseCompleted EventType = "course_completed"
	EventLoginStreak     EventType = "login_streak"
)

var XPRewards = map[EventType]int64{
	EventTestPassed:      100,
	EventVideoWatched:    50,
	EventModuleCompleted: 200,
	EventCourseCompleted: 500,
	EventLoginStreak:     25,
}

// xAPI-compatible event structure
type XAPIEvent struct {
	Actor  XAPIActor   `json:"actor"`
	Verb   XAPIVerb    `json:"verb"`
	Object XAPIObject  `json:"object"`
	Result *XAPIResult `json:"result,omitempty"`
}

type XAPIActor struct {
	Name    string       `json:"name"`
	Mbox    string       `json:"mbox"`    // mailto:user@example.com
	Account *XAPIAccount `json:"account,omitempty"`
}

type XAPIAccount struct {
	HomePage string `json:"homePage"`
	Name     string `json:"name"`
}

type XAPIVerb struct {
	ID      string            `json:"id"`
	Display map[string]string `json:"display"`
}

type XAPIObject struct {
	ID         string          `json:"id"`
	ObjectType string          `json:"objectType"`
	Definition *XAPIDefinition `json:"definition,omitempty"`
}

type XAPIDefinition struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description,omitempty"`
	Type        string            `json:"type,omitempty"`
}

type XAPIResult struct {
	Score    *XAPIScore `json:"score,omitempty"`
	Success  *bool      `json:"success,omitempty"`
	Duration string     `json:"duration,omitempty"`
}

type XAPIScore struct {
	Scaled float64 `json:"scaled"`
	Raw    float64 `json:"raw"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

type EventRecord struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	EventType  EventType `json:"event_type"`
	ObjectID   string    `json:"object_id"`
	XPAwarded  int64     `json:"xp_awarded"`
	RawPayload []byte    `json:"raw_payload,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type EventRepository interface {
	Save(ctx interface{}, event *EventRecord) error
	GetUserEventCount(ctx interface{}, userID uuid.UUID, eventType EventType) (int, error)
	GetUserEventCountSince(ctx interface{}, userID uuid.UUID, eventType EventType, since time.Time) (int, error)
}
