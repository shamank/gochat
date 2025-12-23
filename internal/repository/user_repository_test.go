package repository

import (
	"testing"
	"time"

	"gochat/internal/domain"
)

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &domain.User{
		ID:        "1",
		Username:  "testuser",
		CreatedAt: time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.Create(user)
	if err == nil {
		t.Fatal("Expected error for duplicate username, got nil")
	}
}

func TestInMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &domain.User{
		ID:        "1",
		Username:  "testuser",
		CreatedAt: time.Now(),
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := repo.GetByID("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}

	if retrieved.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, retrieved.Username)
	}

	_, err = repo.GetByID("999")
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}
}

func TestInMemoryUserRepository_GetByUsername(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &domain.User{
		ID:        "1",
		Username:  "testuser",
		CreatedAt: time.Now(),
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, retrieved.Username)
	}
}

func TestInMemoryUserRepository_Exists(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &domain.User{
		ID:        "1",
		Username:  "testuser",
		CreatedAt: time.Now(),
	}

	if repo.Exists("testuser") {
		t.Error("Expected user to not exist")
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if !repo.Exists("testuser") {
		t.Error("Expected user to exist")
	}
}
