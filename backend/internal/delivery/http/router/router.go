package router

import (
	"net/http"
	"time"

	"github.com/corporate-training/gamification/internal/delivery/http/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	webhook := r.Group("/webhook")
	{
		webhook.POST("/event", h.HandleWebhookEvent)
	}

	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// Public
		api.GET("/leaderboard", h.GetLeaderboard)
		api.GET("/achievements", h.GetAchievements)
		api.GET("/skills", h.GetSkillTree)
		api.GET("/quizzes", h.ListQuizzes)
		api.GET("/quizzes/:id", h.GetQuiz)

		users := api.Group("/users")
		{
			users.GET("", h.ListUsers)
			users.GET("/:id", h.GetUser)
			users.GET("/:id/achievements", h.GetUserAchievements)
			users.GET("/:id/skills", h.GetUserSkills)
			users.GET("/:id/quiz-results", h.GetUserQuizResults)
		}

		// Protected (requires JWT)
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware())
		{
			protected.POST("/quizzes/:id/submit", h.SubmitQuiz)
		}
	}

	return r
}
