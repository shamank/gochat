package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gochat/internal/domain"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uc *UserUsecase) RegisterUser(username string) (*domain.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if uc.userRepo.Exists(username) {
		return nil, errors.New("username already exists")
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Username:  username,
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) GetUser(id string) (*domain.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *UserUsecase) GetUserByUsername(username string) (*domain.User, error) {
	return uc.userRepo.GetByUsername(username)
}
