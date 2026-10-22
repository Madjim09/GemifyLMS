package handler_test

import (
	"testing"

	"github.com/corporate-training/gamification/internal/delivery/http/handler"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testDeps struct {
	userRepo        *testutil.MockUserRepository
	achievementRepo *testutil.MockAchievementRepository
	skillRepo       *testutil.MockSkillRepository
	eventRepo       *testutil.MockEventRepository
	quizRepo        *testutil.MockQuizRepository
	handler         *handler.Handler
}

func newTestDeps(t *testing.T) *testDeps {
	t.Helper()

	userRepo := testutil.NewMockUserRepository()
	achievementRepo := testutil.NewMockAchievementRepository()
	skillRepo := testutil.NewMockSkillRepository()
	eventRepo := testutil.NewMockEventRepository()
	quizRepo := testutil.NewMockQuizRepository()
	redisClient := cache.NewDisabledRedisClient()

	gamificationUC := usecase.NewGamificationUseCase(userRepo, achievementRepo, skillRepo, eventRepo, redisClient)
	leaderboardUC := usecase.NewLeaderboardUseCase(userRepo, redisClient)
	achievementUC := usecase.NewAchievementUseCase(achievementRepo, userRepo)
	skillUC := usecase.NewSkillUseCase(skillRepo, userRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	authUC := usecase.NewAuthUseCase(userRepo)
	quizUC := usecase.NewQuizUseCase(quizRepo, userRepo, eventRepo, gamificationUC, redisClient)

	h := handler.NewHandler(gamificationUC, leaderboardUC, achievementUC, skillUC, userUC, authUC, quizUC)

	return &testDeps{
		userRepo:        userRepo,
		achievementRepo: achievementRepo,
		skillRepo:       skillRepo,
		eventRepo:       eventRepo,
		quizRepo:        quizRepo,
		handler:         h,
	}
}
