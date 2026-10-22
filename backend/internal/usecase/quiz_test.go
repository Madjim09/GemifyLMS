package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/google/uuid"
)

func newQuizUC(
	quizRepo *testutil.MockQuizRepository,
	userRepo *testutil.MockUserRepository,
	eventRepo *testutil.MockEventRepository,
) *usecase.QuizUseCase {
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	redisClient := cache.NewDisabledRedisClient()
	gamification := usecase.NewGamificationUseCase(userRepo, achievementRepo, skillRepo, eventRepo, redisClient)
	return usecase.NewQuizUseCase(quizRepo, userRepo, eventRepo, gamification, redisClient)
}

func makeQuiz(passThreshold int, numQ int, allCorrect bool) (*domain.Quiz, map[string]string) {
	quizID := uuid.New()
	answers := make(map[string]string)
	var questions []*domain.QuizQuestion

	for i := 0; i < numQ; i++ {
		qID := uuid.New()
		correctOptID := uuid.New()
		wrongOptID := uuid.New()
		questions = append(questions, &domain.QuizQuestion{
			ID:    qID,
			Text:  fmt.Sprintf("Q%d", i),
			Order: i,
			Options: []*domain.QuizOption{
				{ID: correctOptID, Text: "correct", IsCorrect: true},
				{ID: wrongOptID, Text: "wrong", IsCorrect: false},
			},
		})
		if allCorrect {
			answers[qID.String()] = correctOptID.String()
		} else {
			answers[qID.String()] = wrongOptID.String()
		}
	}

	quiz := &domain.Quiz{
		ID:            quizID,
		Title:         "Test Quiz",
		XPReward:      100,
		BonusXP:       50,
		PassThreshold: passThreshold,
		Questions:     questions,
	}
	return quiz, answers
}

func TestQuizUseCase_ListQuizzes(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	quizRepo.Quizzes[uuid.New()] = &domain.Quiz{ID: uuid.New(), Title: "Q1"}
	quizRepo.Quizzes[uuid.New()] = &domain.Quiz{ID: uuid.New(), Title: "Q2"}

	list, err := uc.ListQuizzes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
}

func TestQuizUseCase_ListQuizzes_Error(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	quizRepo.GetAllErr = fmt.Errorf("db error")
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	_, err := uc.ListQuizzes(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestQuizUseCase_GetQuizForClient(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	optID := uuid.New()
	quizID := uuid.New()
	quizRepo.Quizzes[quizID] = &domain.Quiz{
		ID: quizID,
		Questions: []*domain.QuizQuestion{
			{
				ID: uuid.New(),
				Options: []*domain.QuizOption{
					{ID: optID, IsCorrect: true},
				},
			},
		},
	}

	quiz, err := uc.GetQuizForClient(context.Background(), quizID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// IsCorrect must be masked
	for _, q := range quiz.Questions {
		for _, opt := range q.Options {
			if opt.IsCorrect {
				t.Error("IsCorrect should be masked (false) for client")
			}
		}
	}
}

func TestQuizUseCase_GetQuizForClient_NotFound(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	_, err := uc.GetQuizForClient(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for unknown quiz, got nil")
	}
}

func TestQuizUseCase_SubmitQuiz_Passed(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "u@t.com", XP: 0, Level: 1})

	quiz, answers := makeQuiz(70, 2, true) // 100% correct, pass_threshold=70
	quizRepo.Quizzes[quiz.ID] = quiz

	out, err := uc.SubmitQuiz(context.Background(), userID, quiz.ID, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Passed {
		t.Error("expected Passed = true")
	}
	if !out.IsPerfect {
		t.Error("expected IsPerfect = true")
	}
	if out.XPEarned != quiz.XPReward+quiz.BonusXP {
		t.Errorf("XPEarned = %d, want %d", out.XPEarned, quiz.XPReward+quiz.BonusXP)
	}
}

func TestQuizUseCase_SubmitQuiz_Failed(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "fail@t.com", XP: 0, Level: 1})

	quiz, answers := makeQuiz(70, 2, false) // 0% correct
	quizRepo.Quizzes[quiz.ID] = quiz

	out, err := uc.SubmitQuiz(context.Background(), userID, quiz.ID, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Passed {
		t.Error("expected Passed = false")
	}
	if out.XPEarned != 0 {
		t.Errorf("XPEarned = %d, want 0", out.XPEarned)
	}
}

func TestQuizUseCase_SubmitQuiz_PartialCorrect_Passed(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "partial@t.com", XP: 0, Level: 1})

	// 2 questions, 1 correct → score=50, pass_threshold=50 → passed but not perfect
	q1ID := uuid.New()
	q2ID := uuid.New()
	c1 := uuid.New()
	w1 := uuid.New()
	c2 := uuid.New()
	w2 := uuid.New()

	quiz := &domain.Quiz{
		ID: uuid.New(), Title: "Mixed", XPReward: 100, BonusXP: 50, PassThreshold: 50,
		Questions: []*domain.QuizQuestion{
			{ID: q1ID, Options: []*domain.QuizOption{{ID: c1, IsCorrect: true}, {ID: w1}}},
			{ID: q2ID, Options: []*domain.QuizOption{{ID: c2, IsCorrect: true}, {ID: w2}}},
		},
	}
	quizRepo.Quizzes[quiz.ID] = quiz

	answers := map[string]string{
		q1ID.String(): c1.String(), // correct
		q2ID.String(): w2.String(), // wrong
	}

	out, err := uc.SubmitQuiz(context.Background(), userID, quiz.ID, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Passed {
		t.Error("expected Passed = true (score=50 >= threshold=50)")
	}
	if out.IsPerfect {
		t.Error("expected IsPerfect = false")
	}
	if out.XPEarned != quiz.XPReward {
		t.Errorf("XPEarned = %d, want %d", out.XPEarned, quiz.XPReward)
	}
}

func TestQuizUseCase_SubmitQuiz_NotFound(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	_, err := uc.SubmitQuiz(context.Background(), uuid.New(), uuid.New(), nil)
	if err == nil {
		t.Fatal("expected error for unknown quiz, got nil")
	}
}

func TestQuizUseCase_SubmitQuiz_NoQuestions(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "empty@t.com"})

	quizID := uuid.New()
	quizRepo.Quizzes[quizID] = &domain.Quiz{ID: quizID, Title: "Empty", Questions: nil}

	_, err := uc.SubmitQuiz(context.Background(), userID, quizID, map[string]string{})
	if err == nil {
		t.Fatal("expected error for quiz with no questions, got nil")
	}
}

func TestQuizUseCase_SubmitQuiz_SaveResultError(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	quizRepo.SaveResultErr = fmt.Errorf("db error")
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "saveresult@t.com"})

	quiz, answers := makeQuiz(70, 1, true)
	quizRepo.Quizzes[quiz.ID] = quiz

	_, err := uc.SubmitQuiz(context.Background(), userID, quiz.ID, answers)
	if err == nil {
		t.Fatal("expected error when SaveResult fails, got nil")
	}
}

func TestQuizUseCase_SubmitQuiz_AddXPError(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "addxp@t.com", XP: 0, Level: 1})
	userRepo.AddXPErr = fmt.Errorf("db error")

	quiz, answers := makeQuiz(70, 1, true)
	quizRepo.Quizzes[quiz.ID] = quiz

	_, err := uc.SubmitQuiz(context.Background(), userID, quiz.ID, answers)
	if err == nil {
		t.Fatal("expected error when AddXP fails, got nil")
	}
}

func TestQuizUseCase_GetUserResults(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	quizRepo.UserResults[userID] = []*domain.UserQuizResult{
		{UserID: userID, Score: 80},
		{UserID: userID, Score: 100},
	}

	results, err := uc.GetUserResults(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("len(results) = %d, want 2", len(results))
	}
}

func TestQuizUseCase_SubmitQuiz_UnansweredQuestion(t *testing.T) {
	quizRepo := testutil.NewMockQuizRepository()
	userRepo := testutil.NewMockUserRepository()
	eventRepo := testutil.NewMockEventRepository()
	uc := newQuizUC(quizRepo, userRepo, eventRepo)

	userID := uuid.New()
	userRepo.Create(context.Background(), &domain.User{ID: userID, Email: "unanswered@t.com"})

	quiz, _ := makeQuiz(70, 2, true)
	quizRepo.Quizzes[quiz.ID] = quiz

	// Submit with empty answers (no questions answered)
	out, err := uc.SubmitQuiz(context.Background(), userID, quiz.ID, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Score != 0 {
		t.Errorf("score = %d, want 0 (no answers)", out.Score)
	}
	if out.Passed {
		t.Error("expected Passed = false with no answers")
	}
}
