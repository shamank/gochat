package usecase

import (
	"errors"
	"testing"

	"gochat/internal/domain"
)

type MockUserRepository struct {
	users           map[string]*domain.User
	usersByUsername map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:           make(map[string]*domain.User),
		usersByUsername: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if _, exists := m.usersByUsername[user.Username]; exists {
		return errors.New("username already exists")
	}
	m.users[user.ID] = user
	m.usersByUsername[user.Username] = user
	return nil
}

func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	user, exists := m.usersByUsername[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) Exists(username string) bool {
	_, exists := m.usersByUsername[username]
	return exists
}

func TestUserUsecase_RegisterUser(t *testing.T) {
	repo := NewMockUserRepository()
	usecase := NewUserUsecase(repo)

	user, err := usecase.RegisterUser("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}

	if user.ID == "" {
		t.Error("Expected user ID to be set")
	}

	_, err = usecase.RegisterUser("testuser")
	if err == nil {
		t.Fatal("Expected error for duplicate username, got nil")
	}

	_, err = usecase.RegisterUser("")
	if err == nil {
		t.Fatal("Expected error for empty username, got nil")
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	repo := NewMockUserRepository()
	usecase := NewUserUsecase(repo)

	user, _ := usecase.RegisterUser("testuser")

	retrieved, err := usecase.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}
}

func TestUserUsecase_GetUserByUsername(t *testing.T) {
	repo := NewMockUserRepository()
	usecase := NewUserUsecase(repo)

	user, _ := usecase.RegisterUser("testuser")

	retrieved, err := usecase.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, retrieved.Username)
	}
}
