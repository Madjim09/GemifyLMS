package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx interface{}, user *domain.User) error {
	c := toContext(ctx)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	_, err := r.db.ExecContext(c,
		`INSERT INTO users (id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		user.ID, user.ExternalID, user.Name, user.Email, user.PasswordHash,
		user.XP, user.Level, user.TestsPassed, user.AvatarURL, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx interface{}, id uuid.UUID) (*domain.User, error) {
	c := toContext(ctx)
	user := &domain.User{}
	err := r.db.QueryRowContext(c,
		`SELECT id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at
		FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.ExternalID, &user.Name, &user.Email, &user.PasswordHash,
		&user.XP, &user.Level, &user.TestsPassed, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByExternalID(ctx interface{}, externalID string) (*domain.User, error) {
	c := toContext(ctx)
	user := &domain.User{}
	err := r.db.QueryRowContext(c,
		`SELECT id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at
		FROM users WHERE external_id = $1`, externalID,
	).Scan(&user.ID, &user.ExternalID, &user.Name, &user.Email, &user.PasswordHash,
		&user.XP, &user.Level, &user.TestsPassed, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByEmail(ctx interface{}, email string) (*domain.User, error) {
	c := toContext(ctx)
	user := &domain.User{}
	err := r.db.QueryRowContext(c,
		`SELECT id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at
		FROM users WHERE email = $1 AND password_hash != '' ORDER BY created_at DESC LIMIT 1`, email,
	).Scan(&user.ID, &user.ExternalID, &user.Name, &user.Email, &user.PasswordHash,
		&user.XP, &user.Level, &user.TestsPassed, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *UserRepository) Update(ctx interface{}, user *domain.User) error {
	c := toContext(ctx)
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(c,
		`UPDATE users SET name=$2, email=$3, avatar_url=$4, updated_at=$5 WHERE id=$1`,
		user.ID, user.Name, user.Email, user.AvatarURL, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) AddXP(ctx interface{}, userID uuid.UUID, amount int64) (*domain.User, error) {
	c := toContext(ctx)
	user := &domain.User{}
	err := r.db.QueryRowContext(c,
		`UPDATE users SET xp = xp + $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at`,
		userID, amount,
	).Scan(&user.ID, &user.ExternalID, &user.Name, &user.Email, &user.PasswordHash,
		&user.XP, &user.Level, &user.TestsPassed, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.Level = domain.CalculateLevel(user.XP)
	_, err = r.db.ExecContext(c, `UPDATE users SET level=$2 WHERE id=$1`, userID, user.Level)
	return user, err
}

func (r *UserRepository) IncrementTestsPassed(ctx interface{}, userID uuid.UUID) error {
	c := toContext(ctx)
	_, err := r.db.ExecContext(c,
		`UPDATE users SET tests_passed = tests_passed + 1, updated_at = NOW() WHERE id = $1`, userID)
	return err
}

func (r *UserRepository) GetTopByXP(ctx interface{}, limit int) ([]*domain.User, error) {
	c := toContext(ctx)
	rows, err := r.db.QueryContext(c,
		`SELECT id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at
		FROM users ORDER BY xp DESC LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.ExternalID, &u.Name, &u.Email, &u.PasswordHash,
			&u.XP, &u.Level, &u.TestsPassed, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) List(ctx interface{}, offset, limit int) ([]*domain.User, int64, error) {
	c := toContext(ctx)
	var total int64
	r.db.QueryRowContext(c, `SELECT COUNT(*) FROM users`).Scan(&total)

	rows, err := r.db.QueryContext(c,
		`SELECT id, external_id, name, email, password_hash, xp, level, tests_passed, avatar_url, created_at, updated_at
		FROM users ORDER BY xp DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.ExternalID, &u.Name, &u.Email, &u.PasswordHash,
			&u.XP, &u.Level, &u.TestsPassed, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func toContext(ctx interface{}) context.Context {
	if c, ok := ctx.(context.Context); ok {
		return c
	}
	return context.Background()
}
