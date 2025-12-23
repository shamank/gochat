package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gochat/internal/domain"
)

type RoomUsecase struct {
	roomRepo domain.RoomRepository
}

func NewRoomUsecase(roomRepo domain.RoomRepository) *RoomUsecase {
	return &RoomUsecase{
		roomRepo: roomRepo,
	}
}

func (uc *RoomUsecase) CreateRoom(name string) (*domain.Room, error) {
	if name == "" {
		return nil, errors.New("room name cannot be empty")
	}

	room := &domain.Room{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := uc.roomRepo.Create(room); err != nil {
		return nil, err
	}

	return room, nil
}

func (uc *RoomUsecase) GetRoom(id string) (*domain.Room, error) {
	return uc.roomRepo.GetByID(id)
}

func (uc *RoomUsecase) GetAllRooms() ([]*domain.Room, error) {
	return uc.roomRepo.GetAll()
}

func (uc *RoomUsecase) RoomExists(id string) bool {
	return uc.roomRepo.Exists(id)
}
