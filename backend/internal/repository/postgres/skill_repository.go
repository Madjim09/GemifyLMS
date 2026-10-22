package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type SkillRepository struct {
	db *sql.DB
}

func NewSkillRepository(db *sql.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) GetAll(ctx interface{}) ([]*domain.Skill, error) {
	c := toContext(ctx)
	return r.querySkills(c, `SELECT id, name, description, COALESCE(category,''), parent_id, required_xp, icon_url, "order", created_at FROM skills ORDER BY "order"`)
}

func (r *SkillRepository) GetByID(ctx interface{}, id uuid.UUID) (*domain.Skill, error) {
	c := toContext(ctx)
	skills, err := r.querySkills(c, `SELECT id, name, description, COALESCE(category,''), parent_id, required_xp, icon_url, "order", created_at FROM skills WHERE id=$1`, id)
	if err != nil || len(skills) == 0 {
		return nil, err
	}
	return skills[0], nil
}

func (r *SkillRepository) GetTree(ctx interface{}) ([]*domain.Skill, error) {
	c := toContext(ctx)
	skills, err := r.querySkills(c, `SELECT id, name, description, COALESCE(category,''), parent_id, required_xp, icon_url, "order", created_at FROM skills ORDER BY "order"`)
	if err != nil {
		return nil, err
	}

	skillMap := make(map[uuid.UUID]*domain.Skill)
	for _, s := range skills {
		skillMap[s.ID] = s
	}

	var roots []*domain.Skill
	for _, s := range skills {
		if s.ParentID == nil {
			roots = append(roots, s)
		} else if parent, ok := skillMap[*s.ParentID]; ok {
			parent.Children = append(parent.Children, s)
		}
	}
	return roots, nil
}

func (r *SkillRepository) GetUserSkills(ctx interface{}, userID uuid.UUID) ([]*domain.UserSkill, error) {
	c := toContext(ctx)
	rows, err := r.db.QueryContext(c,
		`SELECT us.id, us.user_id, us.skill_id, us.status, us.progress, us.updated_at,
		s.id, s.name, s.description, COALESCE(s.category,''), s.parent_id, s.required_xp, s.icon_url, s."order", s.created_at
		FROM user_skills us
		JOIN skills s ON us.skill_id = s.id
		WHERE us.user_id = $1
		ORDER BY s."order"`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.UserSkill
	for rows.Next() {
		us := &domain.UserSkill{Skill: &domain.Skill{}}
		if err := rows.Scan(
			&us.ID, &us.UserID, &us.SkillID, &us.Status, &us.Progress, &us.UpdatedAt,
			&us.Skill.ID, &us.Skill.Name, &us.Skill.Description, &us.Skill.Category,
			&us.Skill.ParentID, &us.Skill.RequiredXP, &us.Skill.IconURL, &us.Skill.Order, &us.Skill.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, us)
	}
	return result, rows.Err()
}

func (r *SkillRepository) UpsertUserSkill(ctx interface{}, us *domain.UserSkill) error {
	c := toContext(ctx)
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	us.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(c,
		`INSERT INTO user_skills (id, user_id, skill_id, status, progress, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, skill_id) DO UPDATE SET status=$4, progress=$5, updated_at=$6`,
		us.ID, us.UserID, us.SkillID, us.Status, us.Progress, us.UpdatedAt,
	)
	return err
}

func (r *SkillRepository) querySkills(ctx context.Context, query string, args ...interface{}) ([]*domain.Skill, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*domain.Skill
	for rows.Next() {
		s := &domain.Skill{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Category, &s.ParentID, &s.RequiredXP, &s.IconURL, &s.Order, &s.CreatedAt); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, rows.Err()
}
