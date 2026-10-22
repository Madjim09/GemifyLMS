package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/infrastructure/cache"
)

type LeaderboardEntry struct {
	Rank  int          `json:"rank"`
	User  *domain.User `json:"user"`
	XP    int64        `json:"xp"`
	Level int          `json:"level"`
}

type LeaderboardUseCase struct {
	userRepo domain.UserRepository
	cache    *cache.RedisClient
}

func NewLeaderboardUseCase(userRepo domain.UserRepository, cache *cache.RedisClient) *LeaderboardUseCase {
	return &LeaderboardUseCase{userRepo: userRepo, cache: cache}
}

func (uc *LeaderboardUseCase) GetTopUsers(ctx context.Context, limit int) ([]*LeaderboardEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	cacheKey := "leaderboard:top"
	if cached, err := uc.cache.Get(ctx, cacheKey); err == nil {
		var entries []*LeaderboardEntry
		if err := json.Unmarshal([]byte(cached), &entries); err == nil {
			return entries, nil
		}
	}

	users, err := uc.userRepo.GetTopByXP(ctx, limit)
	if err != nil {
		return nil, err
	}

	entries := make([]*LeaderboardEntry, len(users))
	for i, u := range users {
		entries[i] = &LeaderboardEntry{
			Rank:  i + 1,
			User:  u,
			XP:    u.XP,
			Level: u.Level,
		}
	}

	if data, err := json.Marshal(entries); err == nil {
		if err := uc.cache.Set(ctx, cacheKey, string(data), 5*time.Minute); err != nil {
			log.Printf("Failed to cache leaderboard: %v", err)
		}
	}

	return entries, nil
}
