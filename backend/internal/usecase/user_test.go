package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/corporate-training/gamification/internal/testutil"
	"github.com/corporate-training/gamification/internal/usecase"
	"github.com/google/uuid"
)

func TestUserUseCase_GetByID_Success(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewUserUseCase(userRepo)

	id := uuid.New()
	user := &domain.User{ID: id, Name: "Eve", Email: "eve@example.com"}
	userRepo.Create(context.Background(), user)

	got, err := uc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != id {
		t.Errorf("GetByID returned wrong user: got %v, want %v", got.ID, id)
	}
}

func TestUserUseCase_GetByID_NotFound(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewUserUseCase(userRepo)

	_, err := uc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for unknown ID, got nil")
	}
}

func TestUserUseCase_GetByID_RepoError(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	userRepo.GetByIDErr = fmt.Errorf("db error")
	uc := usecase.NewUserUseCase(userRepo)

	_, err := uc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserUseCase_List(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewUserUseCase(userRepo)

	for i := 0; i < 5; i++ {
		userRepo.Create(context.Background(), &domain.User{
			ID:    uuid.New(),
			Email: fmt.Sprintf("user%d@example.com", i),
		})
	}

	users, total, err := uc.List(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(users) != 5 {
		t.Errorf("len(users) = %d, want 5", len(users))
	}
}

func TestUserUseCase_List_Pagination(t *testing.T) {
	userRepo := testutil.NewMockUserRepository()
	uc := usecase.NewUserUseCase(userRepo)

	for i := 0; i < 5; i++ {
		userRepo.Create(context.Background(), &domain.User{
			ID:    uuid.New(),
			Email: fmt.Sprintf("pager%d@example.com", i),
		})
	}

	users, total, err := uc.List(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(users) != 0 {
		t.Errorf("expected empty page, got %d users", len(users))
	}
}
