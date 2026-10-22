package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo domain.UserRepository
}

func NewAuthUseCase(userRepo domain.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

type AuthOutput struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

func (uc *AuthUseCase) Register(ctx context.Context, name, email, password string) (*AuthOutput, error) {
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, fmt.Errorf("email уже зарегистрирован")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		ExternalID:   uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		XP:           0,
		Level:        1,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{User: user, Token: token}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*AuthOutput, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{User: user, Token: token}, nil
}

func GenerateToken(userID uuid.UUID) (string, error) {
	secret := getJWTSecret()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string) (uuid.UUID, error) {
	secret := getJWTSecret()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("недействительный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("недействительные данные токена")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("отсутствует subject в токене")
	}

	return uuid.Parse(sub)
}

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "gamification-jwt-secret-change-in-production"
}
