package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/corporate-training/gamification/internal/delivery/http/handler"
	"github.com/corporate-training/gamification/internal/delivery/http/router"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
	"github.com/corporate-training/gamification/internal/infrastructure/database"
	"github.com/corporate-training/gamification/internal/repository/postgres"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	database.SeedDemoData(db)

	redisClient := cache.NewRedisClient()

	userRepo := postgres.NewUserRepository(db)
	achievementRepo := postgres.NewAchievementRepository(db)
	skillRepo := postgres.NewSkillRepository(db)
	eventRepo := postgres.NewEventRepository(db)
	quizRepo := postgres.NewQuizRepository(db)

	gamificationUC := usecase.NewGamificationUseCase(userRepo, achievementRepo, skillRepo, eventRepo, redisClient)
	leaderboardUC := usecase.NewLeaderboardUseCase(userRepo, redisClient)
	achievementUC := usecase.NewAchievementUseCase(achievementRepo, userRepo)
	skillUC := usecase.NewSkillUseCase(skillRepo, userRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	authUC := usecase.NewAuthUseCase(userRepo)
	quizUC := usecase.NewQuizUseCase(quizRepo, userRepo, eventRepo, gamificationUC, redisClient)

	h := handler.NewHandler(gamificationUC, leaderboardUC, achievementUC, skillUC, userUC, authUC, quizUC)

	r := router.NewRouter(h)

	srv := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", getEnv("PORT", "8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
