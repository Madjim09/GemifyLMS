package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
	"github.com/google/uuid"
)

type QuizUseCase struct {
	quizRepo    domain.QuizRepository
	userRepo    domain.UserRepository
	eventRepo   domain.EventRepository
	gamification *GamificationUseCase
	cache       *cache.RedisClient
}

func NewQuizUseCase(
	quizRepo domain.QuizRepository,
	userRepo domain.UserRepository,
	eventRepo domain.EventRepository,
	gamification *GamificationUseCase,
	redisClient *cache.RedisClient,
) *QuizUseCase {
	return &QuizUseCase{
		quizRepo:    quizRepo,
		userRepo:    userRepo,
		eventRepo:   eventRepo,
		gamification: gamification,
		cache:       redisClient,
	}
}

type SubmitQuizOutput struct {
	Score     int                    `json:"score"`
	Passed    bool                   `json:"passed"`
	IsPerfect bool                   `json:"is_perfect"`
	XPEarned  int64                  `json:"xp_earned"`
	Result    *domain.UserQuizResult `json:"result"`
}

func (uc *QuizUseCase) ListQuizzes(ctx context.Context) ([]*domain.Quiz, error) {
	return uc.quizRepo.GetAll(ctx)
}

// GetQuizForClient returns quiz without correct answer markers.
func (uc *QuizUseCase) GetQuizForClient(ctx context.Context, id uuid.UUID) (*domain.Quiz, error) {
	quiz, err := uc.quizRepo.GetByIDWithQuestions(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, q := range quiz.Questions {
		for _, opt := range q.Options {
			opt.IsCorrect = false
		}
	}
	return quiz, nil
}

// SubmitQuiz validates answers server-side and awards XP.
func (uc *QuizUseCase) SubmitQuiz(ctx context.Context, userID, quizID uuid.UUID, answers map[string]string) (*SubmitQuizOutput, error) {
	quiz, err := uc.quizRepo.GetByIDWithQuestions(ctx, quizID)
	if err != nil {
		return nil, fmt.Errorf("тест не найден")
	}

	total := len(quiz.Questions)
	if total == 0 {
		return nil, fmt.Errorf("тест не содержит вопросов")
	}

	correct := 0
	for _, q := range quiz.Questions {
		submittedID, ok := answers[q.ID.String()]
		if !ok {
			continue
		}
		for _, opt := range q.Options {
			if opt.ID.String() == submittedID && opt.IsCorrect {
				correct++
				break
			}
		}
	}

	score := (correct * 100) / total
	passed := score >= quiz.PassThreshold
	isPerfect := score == 100

	var xpEarned int64
	if passed {
		xpEarned = quiz.XPReward
		if isPerfect {
			xpEarned += quiz.BonusXP
		}
	}

	result := &domain.UserQuizResult{
		UserID:      userID,
		QuizID:      quizID,
		Score:       score,
		XPEarned:    xpEarned,
		IsPerfect:   isPerfect,
		CompletedAt: time.Now(),
	}

	if err := uc.quizRepo.SaveResult(ctx, result); err != nil {
		return nil, fmt.Errorf("ошибка сохранения результата: %w", err)
	}

	if passed {
		updatedUser, err := uc.userRepo.AddXP(ctx, userID, xpEarned)
		if err != nil {
			return nil, fmt.Errorf("ошибка начисления XP: %w", err)
		}

		uc.userRepo.IncrementTestsPassed(ctx, userID)
		uc.cache.InvalidateLeaderboard(ctx)

		passedEvent := &domain.EventRecord{
			ID:        uuid.New(),
			UserID:    userID,
			EventType: domain.EventTestPassed,
			ObjectID:  quizID.String(),
			XPAwarded: xpEarned,
			CreatedAt: time.Now(),
		}
		uc.eventRepo.Save(ctx, passedEvent)

		if isPerfect {
			perfectEvent := &domain.EventRecord{
				ID:        uuid.New(),
				UserID:    userID,
				EventType: domain.EventTestPerfect,
				ObjectID:  quizID.String(),
				XPAwarded: 0,
				CreatedAt: time.Now(),
			}
			uc.eventRepo.Save(ctx, perfectEvent)
		}

		uc.gamification.CheckAchievements(ctx, updatedUser, domain.EventTestPassed)
	}

	return &SubmitQuizOutput{
		Score:     score,
		Passed:    passed,
		IsPerfect: isPerfect,
		XPEarned:  xpEarned,
		Result:    result,
	}, nil
}

func (uc *QuizUseCase) GetUserResults(ctx context.Context, userID uuid.UUID) ([]*domain.UserQuizResult, error) {
	return uc.quizRepo.GetUserResults(ctx, userID)
}
