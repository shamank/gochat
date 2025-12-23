package repository

import (
	"errors"
	"sync"

	"gochat/internal/domain"
)

type InMemoryUserRepository struct {
	users           map[string]*domain.User
	usersByUsername map[string]*domain.User
	mu              sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:           make(map[string]*domain.User),
		usersByUsername: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByUsername[user.Username]; exists {
		return errors.New("username already exists")
	}

	r.users[user.ID] = user
	r.usersByUsername[user.Username] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByUsername[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *InMemoryUserRepository) Exists(username string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.usersByUsername[username]
	return exists
}
