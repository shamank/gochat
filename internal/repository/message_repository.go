package repository

import (
	"errors"
	"sort"
	"sync"

	"gochat/internal/domain"
)

type InMemoryMessageRepository struct {
	messages     map[string]*domain.Message
	roomMessages map[string][]*domain.Message
	mu           sync.RWMutex
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		messages:     make(map[string]*domain.Message),
		roomMessages: make(map[string][]*domain.Message),
	}
}

func (r *InMemoryMessageRepository) Create(message *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[message.ID] = message
	r.roomMessages[message.RoomID] = append(r.roomMessages[message.RoomID], message)
	return nil
}

func (r *InMemoryMessageRepository) GetByRoomID(roomID string, limit, offset int) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	messages, exists := r.roomMessages[roomID]
	if !exists {
		return []*domain.Message{}, nil
	}

	sortedMessages := make([]*domain.Message, len(messages))
	copy(sortedMessages, messages)
	sort.Slice(sortedMessages, func(i, j int) bool {
		return sortedMessages[i].CreatedAt.Before(sortedMessages[j].CreatedAt)
	})

	start := offset
	if start > len(sortedMessages) {
		return []*domain.Message{}, nil
	}

	end := start + limit
	if end > len(sortedMessages) {
		end = len(sortedMessages)
	}

	return sortedMessages[start:end], nil
}

func (r *InMemoryMessageRepository) GetByID(id string) (*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	message, exists := r.messages[id]
	if !exists {
		return nil, errors.New("message not found")
	}

	return message, nil
}
