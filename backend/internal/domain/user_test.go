package domain_test

import (
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
)

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		xp    int64
		level int
	}{
		{0, 1},
		{100, 1},
		{499, 1},
		{500, 2},
		{1499, 2},
		{1500, 3},
		{2999, 3},
		{3000, 4},
		{5499, 4},
		{5500, 5},
		{8999, 5},
		{9000, 6},
		{13999, 6},
		{14000, 7},
		{20999, 7},
		{21000, 8},
		{29999, 8},
		{30000, 9},
		{41999, 9},
		{42000, 10},
		{100000, 10},
	}
	for _, tc := range tests {
		got := domain.CalculateLevel(tc.xp)
		if got != tc.level {
			t.Errorf("CalculateLevel(%d) = %d, want %d", tc.xp, got, tc.level)
		}
	}
}

func TestXPForNextLevel(t *testing.T) {
	tests := []struct {
		level int
		xp    int64
	}{
		{1, 500},
		{2, 1500},
		{3, 3000},
		{4, 5500},
		{5, 9000},
		{6, 14000},
		{7, 21000},
		{8, 30000},
		{9, 42000},
		{10, -1},
		{11, -1},
		{100, -1},
	}
	for _, tc := range tests {
		got := domain.XPForNextLevel(tc.level)
		if got != tc.xp {
			t.Errorf("XPForNextLevel(%d) = %d, want %d", tc.level, got, tc.xp)
		}
	}
}

func TestRecalculateLevel(t *testing.T) {
	tests := []struct {
		xp    int64
		level int
	}{
		{0, 1},
		{500, 2},
		{1500, 3},
		{42000, 10},
	}
	for _, tc := range tests {
		u := &domain.User{XP: tc.xp}
		u.RecalculateLevel()
		if u.Level != tc.level {
			t.Errorf("RecalculateLevel() with XP=%d: level=%d, want %d", tc.xp, u.Level, tc.level)
		}
	}
}
