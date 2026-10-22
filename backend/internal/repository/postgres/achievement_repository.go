package postgres

import (
	"database/sql"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type AchievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) GetAll(ctx interface{}) ([]*domain.Achievement, error) {
	c := toContext(ctx)
	rows, err := r.db.QueryContext(c,
		`SELECT id, name, description, icon_url, condition_type, condition_value, xp_reward, created_at FROM achievements ORDER BY condition_value`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*domain.Achievement
	for rows.Next() {
		a := &domain.Achievement{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.IconURL, &a.ConditionType, &a.ConditionValue, &a.XPReward, &a.CreatedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, rows.Err()
}

func (r *AchievementRepository) GetByID(ctx interface{}, id uuid.UUID) (*domain.Achievement, error) {
	c := toContext(ctx)
	a := &domain.Achievement{}
	err := r.db.QueryRowContext(c,
		`SELECT id, name, description, icon_url, condition_type, condition_value, xp_reward, created_at FROM achievements WHERE id=$1`, id,
	).Scan(&a.ID, &a.Name, &a.Description, &a.IconURL, &a.ConditionType, &a.ConditionValue, &a.XPReward, &a.CreatedAt)
	return a, err
}

func (r *AchievementRepository) GetUserAchievements(ctx interface{}, userID uuid.UUID) ([]*domain.UserAchievement, error) {
	c := toContext(ctx)
	rows, err := r.db.QueryContext(c,
		`SELECT ua.id, ua.user_id, ua.achievement_id, ua.earned_at,
		a.id, a.name, a.description, a.icon_url, a.condition_type, a.condition_value, a.xp_reward, a.created_at
		FROM user_achievements ua
		JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.earned_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.UserAchievement
	for rows.Next() {
		ua := &domain.UserAchievement{Achievement: &domain.Achievement{}}
		if err := rows.Scan(
			&ua.ID, &ua.UserID, &ua.AchievementID, &ua.EarnedAt,
			&ua.Achievement.ID, &ua.Achievement.Name, &ua.Achievement.Description,
			&ua.Achievement.IconURL, &ua.Achievement.ConditionType, &ua.Achievement.ConditionValue,
			&ua.Achievement.XPReward, &ua.Achievement.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ua)
	}
	return result, rows.Err()
}

func (r *AchievementRepository) AwardToUser(ctx interface{}, userID, achievementID uuid.UUID) (*domain.UserAchievement, error) {
	c := toContext(ctx)
	ua := &domain.UserAchievement{
		ID:            uuid.New(),
		UserID:        userID,
		AchievementID: achievementID,
		EarnedAt:      time.Now(),
	}
	_, err := r.db.ExecContext(c,
		`INSERT INTO user_achievements (id, user_id, achievement_id, earned_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		ua.ID, ua.UserID, ua.AchievementID, ua.EarnedAt,
	)
	return ua, err
}

func (r *AchievementRepository) HasUserAchievement(ctx interface{}, userID, achievementID uuid.UUID) (bool, error) {
	c := toContext(ctx)
	var count int
	err := r.db.QueryRowContext(c,
		`SELECT COUNT(*) FROM user_achievements WHERE user_id=$1 AND achievement_id=$2`,
		userID, achievementID,
	).Scan(&count)
	return count > 0, err
}
