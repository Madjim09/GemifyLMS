package postgres

import (
	"database/sql"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) GetAll(ctx interface{}) ([]*domain.Quiz, error) {
	c := toContext(ctx)
	rows, err := r.db.QueryContext(c,
		`SELECT id, title, description, difficulty, xp_reward, bonus_xp, pass_threshold, created_at
		FROM quizzes ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*domain.Quiz
	for rows.Next() {
		q := &domain.Quiz{}
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.Difficulty,
			&q.XPReward, &q.BonusXP, &q.PassThreshold, &q.CreatedAt); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, rows.Err()
}

func (r *QuizRepository) GetByIDWithQuestions(ctx interface{}, id uuid.UUID) (*domain.Quiz, error) {
	c := toContext(ctx)

	quiz := &domain.Quiz{}
	err := r.db.QueryRowContext(c,
		`SELECT id, title, description, difficulty, xp_reward, bonus_xp, pass_threshold, created_at
		FROM quizzes WHERE id = $1`, id,
	).Scan(&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Difficulty,
		&quiz.XPReward, &quiz.BonusXP, &quiz.PassThreshold, &quiz.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	qRows, err := r.db.QueryContext(c,
		`SELECT id, quiz_id, text, "order" FROM quiz_questions WHERE quiz_id = $1 ORDER BY "order"`, id)
	if err != nil {
		return nil, err
	}
	defer qRows.Close()

	questionMap := make(map[uuid.UUID]*domain.QuizQuestion)
	for qRows.Next() {
		q := &domain.QuizQuestion{}
		if err := qRows.Scan(&q.ID, &q.QuizID, &q.Text, &q.Order); err != nil {
			return nil, err
		}
		quiz.Questions = append(quiz.Questions, q)
		questionMap[q.ID] = q
	}
	if err := qRows.Err(); err != nil {
		return nil, err
	}

	if len(quiz.Questions) == 0 {
		return quiz, nil
	}

	oRows, err := r.db.QueryContext(c,
		`SELECT qo.id, qo.question_id, qo.text, qo.is_correct, qo.order
		FROM quiz_options qo
		JOIN quiz_questions qq ON qo.question_id = qq.id
		WHERE qq.quiz_id = $1
		ORDER BY qo.question_id, qo.order`, id)
	if err != nil {
		return nil, err
	}
	defer oRows.Close()

	for oRows.Next() {
		opt := &domain.QuizOption{}
		if err := oRows.Scan(&opt.ID, &opt.QuestionID, &opt.Text, &opt.IsCorrect, &opt.Order); err != nil {
			return nil, err
		}
		if q, ok := questionMap[opt.QuestionID]; ok {
			q.Options = append(q.Options, opt)
		}
	}

	return quiz, oRows.Err()
}

func (r *QuizRepository) SaveResult(ctx interface{}, result *domain.UserQuizResult) error {
	c := toContext(ctx)
	if result.ID == uuid.Nil {
		result.ID = uuid.New()
	}
	result.CompletedAt = time.Now()
	_, err := r.db.ExecContext(c,
		`INSERT INTO user_quiz_results (id, user_id, quiz_id, score, xp_earned, is_perfect, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		result.ID, result.UserID, result.QuizID, result.Score, result.XPEarned, result.IsPerfect, result.CompletedAt,
	)
	return err
}

func (r *QuizRepository) GetUserResults(ctx interface{}, userID uuid.UUID) ([]*domain.UserQuizResult, error) {
	c := toContext(ctx)
	rows, err := r.db.QueryContext(c,
		`SELECT uqr.id, uqr.user_id, uqr.quiz_id, uqr.score, uqr.xp_earned, uqr.is_perfect, uqr.completed_at,
		q.id, q.title, q.description, q.difficulty, q.xp_reward, q.bonus_xp, q.pass_threshold, q.created_at
		FROM user_quiz_results uqr
		JOIN quizzes q ON uqr.quiz_id = q.id
		WHERE uqr.user_id = $1
		ORDER BY uqr.completed_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.UserQuizResult
	for rows.Next() {
		r := &domain.UserQuizResult{Quiz: &domain.Quiz{}}
		if err := rows.Scan(
			&r.ID, &r.UserID, &r.QuizID, &r.Score, &r.XPEarned, &r.IsPerfect, &r.CompletedAt,
			&r.Quiz.ID, &r.Quiz.Title, &r.Quiz.Description, &r.Quiz.Difficulty,
			&r.Quiz.XPReward, &r.Quiz.BonusXP, &r.Quiz.PassThreshold, &r.Quiz.CreatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
