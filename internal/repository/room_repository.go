package repository

import (
	"errors"
	"sync"

	"gochat/internal/domain"
)

type InMemoryRoomRepository struct {
	rooms map[string]*domain.Room
	mu    sync.RWMutex
}

func NewInMemoryRoomRepository() *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		rooms: make(map[string]*domain.Room),
	}
}

func (r *InMemoryRoomRepository) Create(room *domain.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rooms[room.ID] = room
	return nil
}

func (r *InMemoryRoomRepository) GetByID(id string) (*domain.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, exists := r.rooms[id]
	if !exists {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (r *InMemoryRoomRepository) GetAll() ([]*domain.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rooms := make([]*domain.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *InMemoryRoomRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.rooms[id]
	return exists
}
