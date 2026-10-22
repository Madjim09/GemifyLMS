package testutil

import (
	"fmt"
	"sync"
	"time"

	"github.com/corporate-training/gamification/internal/domain"
	"github.com/google/uuid"
)

// MockUserRepository is an in-memory domain.UserRepository for testing.
type MockUserRepository struct {
	mu              sync.RWMutex
	users           map[uuid.UUID]*domain.User
	emailIndex      map[string]*domain.User
	externalIDIndex map[string]*domain.User

	CreateErr          error
	GetByIDErr         error
	GetByEmailErr      error
	GetByExternalIDErr error
	AddXPErr           error
	GetTopByXPErr      error
	ListErr            error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:           make(map[uuid.UUID]*domain.User),
		emailIndex:      make(map[string]*domain.User),
		externalIDIndex: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx interface{}, user *domain.User) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	copied := *user
	m.users[user.ID] = &copied
	m.emailIndex[user.Email] = &copied
	m.externalIDIndex[user.ExternalID] = &copied
	return nil
}

func (m *MockUserRepository) GetByID(ctx interface{}, id uuid.UUID) (*domain.User, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	copied := *u
	return &copied, nil
}

func (m *MockUserRepository) GetByExternalID(ctx interface{}, externalID string) (*domain.User, error) {
	if m.GetByExternalIDErr != nil {
		return nil, m.GetByExternalIDErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.externalIDIndex[externalID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	copied := *u
	return &copied, nil
}

func (m *MockUserRepository) GetByEmail(ctx interface{}, email string) (*domain.User, error) {
	if m.GetByEmailErr != nil {
		return nil, m.GetByEmailErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.emailIndex[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	copied := *u
	return &copied, nil
}

func (m *MockUserRepository) Update(ctx interface{}, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	copied := *user
	m.users[user.ID] = &copied
	return nil
}

func (m *MockUserRepository) AddXP(ctx interface{}, userID uuid.UUID, amount int64) (*domain.User, error) {
	if m.AddXPErr != nil {
		return nil, m.AddXPErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	u.XP += amount
	u.RecalculateLevel()
	copied := *u
	return &copied, nil
}

func (m *MockUserRepository) IncrementTestsPassed(ctx interface{}, userID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[userID]; ok {
		u.TestsPassed++
	}
	return nil
}

func (m *MockUserRepository) GetTopByXP(ctx interface{}, limit int) ([]*domain.User, error) {
	if m.GetTopByXPErr != nil {
		return nil, m.GetTopByXPErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var users []*domain.User
	for _, u := range m.users {
		copied := *u
		users = append(users, &copied)
	}
	if len(users) > limit {
		users = users[:limit]
	}
	return users, nil
}

func (m *MockUserRepository) List(ctx interface{}, offset, limit int) ([]*domain.User, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var users []*domain.User
	for _, u := range m.users {
		copied := *u
		users = append(users, &copied)
	}
	total := int64(len(users))
	if offset >= len(users) {
		return []*domain.User{}, total, nil
	}
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}
	return users[offset:end], total, nil
}

// MockAchievementRepository is an in-memory domain.AchievementRepository for testing.
type MockAchievementRepository struct {
	Achievements     []*domain.Achievement
	UserAchievements map[uuid.UUID][]*domain.UserAchievement

	GetAllErr                error
	AwardToUserErr           error
	HasResult                map[string]bool
	GetUserAchievementsErr   error
}

func NewMockAchievementRepository() *MockAchievementRepository {
	return &MockAchievementRepository{
		UserAchievements: make(map[uuid.UUID][]*domain.UserAchievement),
		HasResult:        make(map[string]bool),
	}
}

func (m *MockAchievementRepository) GetAll(ctx interface{}) ([]*domain.Achievement, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	return m.Achievements, nil
}

func (m *MockAchievementRepository) GetByID(ctx interface{}, id uuid.UUID) (*domain.Achievement, error) {
	for _, a := range m.Achievements {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, fmt.Errorf("achievement not found")
}

func (m *MockAchievementRepository) GetUserAchievements(ctx interface{}, userID uuid.UUID) ([]*domain.UserAchievement, error) {
	if m.GetUserAchievementsErr != nil {
		return nil, m.GetUserAchievementsErr
	}
	return m.UserAchievements[userID], nil
}

func (m *MockAchievementRepository) AwardToUser(ctx interface{}, userID, achievementID uuid.UUID) (*domain.UserAchievement, error) {
	if m.AwardToUserErr != nil {
		return nil, m.AwardToUserErr
	}
	ua := &domain.UserAchievement{
		ID:            uuid.New(),
		UserID:        userID,
		AchievementID: achievementID,
		EarnedAt:      time.Now(),
	}
	m.UserAchievements[userID] = append(m.UserAchievements[userID], ua)
	return ua, nil
}

func (m *MockAchievementRepository) HasUserAchievement(ctx interface{}, userID, achievementID uuid.UUID) (bool, error) {
	key := userID.String() + ":" + achievementID.String()
	return m.HasResult[key], nil
}

// MockSkillRepository is an in-memory domain.SkillRepository for testing.
type MockSkillRepository struct {
	Skills     []*domain.Skill
	UserSkills map[uuid.UUID][]*domain.UserSkill

	GetTreeErr       error
	GetAllErr        error
	GetUserSkillsErr error
}

func NewMockSkillRepository() *MockSkillRepository {
	return &MockSkillRepository{
		UserSkills: make(map[uuid.UUID][]*domain.UserSkill),
	}
}

func (m *MockSkillRepository) GetAll(ctx interface{}) ([]*domain.Skill, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	return m.Skills, nil
}

func (m *MockSkillRepository) GetByID(ctx interface{}, id uuid.UUID) (*domain.Skill, error) {
	for _, s := range m.Skills {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, fmt.Errorf("skill not found")
}

func (m *MockSkillRepository) GetTree(ctx interface{}) ([]*domain.Skill, error) {
	if m.GetTreeErr != nil {
		return nil, m.GetTreeErr
	}
	return m.Skills, nil
}

func (m *MockSkillRepository) GetUserSkills(ctx interface{}, userID uuid.UUID) ([]*domain.UserSkill, error) {
	if m.GetUserSkillsErr != nil {
		return nil, m.GetUserSkillsErr
	}
	return m.UserSkills[userID], nil
}

func (m *MockSkillRepository) UpsertUserSkill(ctx interface{}, userSkill *domain.UserSkill) error {
	m.UserSkills[userSkill.UserID] = append(m.UserSkills[userSkill.UserID], userSkill)
	return nil
}

// MockEventRepository is an in-memory domain.EventRepository for testing.
type MockEventRepository struct {
	Events   []*domain.EventRecord
	CountMap map[string]int

	SaveErr error
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		CountMap: make(map[string]int),
	}
}

func (m *MockEventRepository) Save(ctx interface{}, event *domain.EventRecord) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	m.Events = append(m.Events, event)
	key := event.UserID.String() + ":" + string(event.EventType)
	m.CountMap[key]++
	return nil
}

func (m *MockEventRepository) GetUserEventCount(ctx interface{}, userID uuid.UUID, eventType domain.EventType) (int, error) {
	key := userID.String() + ":" + string(eventType)
	return m.CountMap[key], nil
}

func (m *MockEventRepository) GetUserEventCountSince(ctx interface{}, userID uuid.UUID, eventType domain.EventType, since time.Time) (int, error) {
	key := userID.String() + ":" + string(eventType)
	return m.CountMap[key], nil
}

// MockQuizRepository is an in-memory domain.QuizRepository for testing.
type MockQuizRepository struct {
	Quizzes     map[uuid.UUID]*domain.Quiz
	UserResults map[uuid.UUID][]*domain.UserQuizResult

	GetAllErr               error
	GetByIDWithQuestionsErr error
	SaveResultErr           error
	GetUserResultsErr       error
}

func NewMockQuizRepository() *MockQuizRepository {
	return &MockQuizRepository{
		Quizzes:     make(map[uuid.UUID]*domain.Quiz),
		UserResults: make(map[uuid.UUID][]*domain.UserQuizResult),
	}
}

func (m *MockQuizRepository) GetAll(ctx interface{}) ([]*domain.Quiz, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	var quizzes []*domain.Quiz
	for _, q := range m.Quizzes {
		quizzes = append(quizzes, q)
	}
	return quizzes, nil
}

func (m *MockQuizRepository) GetByIDWithQuestions(ctx interface{}, id uuid.UUID) (*domain.Quiz, error) {
	if m.GetByIDWithQuestionsErr != nil {
		return nil, m.GetByIDWithQuestionsErr
	}
	q, ok := m.Quizzes[id]
	if !ok {
		return nil, fmt.Errorf("quiz not found")
	}
	return q, nil
}

func (m *MockQuizRepository) SaveResult(ctx interface{}, result *domain.UserQuizResult) error {
	if m.SaveResultErr != nil {
		return m.SaveResultErr
	}
	m.UserResults[result.UserID] = append(m.UserResults[result.UserID], result)
	return nil
}

func (m *MockQuizRepository) GetUserResults(ctx interface{}, userID uuid.UUID) ([]*domain.UserQuizResult, error) {
	if m.GetUserResultsErr != nil {
		return nil, m.GetUserResultsErr
	}
	return m.UserResults[userID], nil
}
