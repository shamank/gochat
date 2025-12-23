package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gochat/internal/domain"
)

type MessageUsecase struct {
	messageRepo domain.MessageRepository
	userRepo    domain.UserRepository
	roomRepo    domain.RoomRepository
}

func NewMessageUsecase(
	messageRepo domain.MessageRepository,
	userRepo domain.UserRepository,
	roomRepo domain.RoomRepository,
) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
	}
}

func (uc *MessageUsecase) SendMessage(roomID, userID, content string) (*domain.Message, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !uc.roomRepo.Exists(roomID) {
		return nil, errors.New("room not found")
	}

	message := &domain.Message{
		ID:        uuid.New().String(),
		RoomID:    roomID,
		UserID:    userID,
		Username:  user.Username,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := uc.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (uc *MessageUsecase) GetMessagesHistory(roomID string, limit, offset int) ([]*domain.Message, error) {
	const (
		defaultLimit = 50
		maxLimit     = 100
	)

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return uc.messageRepo.GetByRoomID(roomID, limit, offset)
}
