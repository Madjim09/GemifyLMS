package postgres

import (
	"database/sql"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Save(ctx interface{}, event *domain.EventRecord) error {
	c := toContext(ctx)
	_, err := r.db.ExecContext(c,
		`INSERT INTO event_records (id, user_id, event_type, object_id, xp_awarded, raw_payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ID, event.UserID, event.EventType, event.ObjectID, event.XPAwarded, event.RawPayload, event.CreatedAt,
	)
	return err
}

func (r *EventRepository) GetUserEventCount(ctx interface{}, userID uuid.UUID, eventType domain.EventType) (int, error) {
	c := toContext(ctx)
	var count int
	err := r.db.QueryRowContext(c,
		`SELECT COUNT(*) FROM event_records WHERE user_id=$1 AND event_type=$2`,
		userID, eventType,
	).Scan(&count)
	return count, err
}

func (r *EventRepository) GetUserEventCountSince(ctx interface{}, userID uuid.UUID, eventType domain.EventType, since time.Time) (int, error) {
	c := toContext(ctx)
	var count int
	err := r.db.QueryRowContext(c,
		`SELECT COUNT(*) FROM event_records WHERE user_id=$1 AND event_type=$2 AND created_at >= $3`,
		userID, eventType, since,
	).Scan(&count)
	return count, err
}
