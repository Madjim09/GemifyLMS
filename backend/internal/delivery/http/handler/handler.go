package handler

import (
	"github.com/corporate-training/gamification/internal/usecase"
)

type Handler struct {
	gamification *usecase.GamificationUseCase
	leaderboard  *usecase.LeaderboardUseCase
	achievement  *usecase.AchievementUseCase
	skill        *usecase.SkillUseCase
	user         *usecase.UserUseCase
	auth         *usecase.AuthUseCase
	quiz         *usecase.QuizUseCase
}

func NewHandler(
	gamification *usecase.GamificationUseCase,
	leaderboard *usecase.LeaderboardUseCase,
	achievement *usecase.AchievementUseCase,
	skill *usecase.SkillUseCase,
	user *usecase.UserUseCase,
	auth *usecase.AuthUseCase,
	quiz *usecase.QuizUseCase,
) *Handler {
	return &Handler{
		gamification: gamification,
		leaderboard:  leaderboard,
		achievement:  achievement,
		skill:        skill,
		user:         user,
		auth:         auth,
		quiz:         quiz,
	}
}

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
