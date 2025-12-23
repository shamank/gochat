package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"gochat/internal/domain"
)

type MockMessageRepository struct {
	messages     map[string]*domain.Message
	roomMessages map[string][]*domain.Message
}

func NewMockMessageRepository() *MockMessageRepository {
	return &MockMessageRepository{
		messages:     make(map[string]*domain.Message),
		roomMessages: make(map[string][]*domain.Message),
	}
}

func (m *MockMessageRepository) Create(message *domain.Message) error {
	m.messages[message.ID] = message
	m.roomMessages[message.RoomID] = append(m.roomMessages[message.RoomID], message)
	return nil
}

func (m *MockMessageRepository) GetByRoomID(roomID string, limit, offset int) ([]*domain.Message, error) {
	messages, exists := m.roomMessages[roomID]
	if !exists {
		return []*domain.Message{}, nil
	}

	start := offset
	if start > len(messages) {
		return []*domain.Message{}, nil
	}

	end := start + limit
	if end > len(messages) {
		end = len(messages)
	}

	return messages[start:end], nil
}

func (m *MockMessageRepository) GetByID(id string) (*domain.Message, error) {
	message, exists := m.messages[id]
	if !exists {
		return nil, errors.New("message not found")
	}
	return message, nil
}

type MockRoomRepository struct {
	rooms map[string]*domain.Room
}

func NewMockRoomRepository() *MockRoomRepository {
	return &MockRoomRepository{
		rooms: make(map[string]*domain.Room),
	}
}

func (m *MockRoomRepository) Create(room *domain.Room) error {
	m.rooms[room.ID] = room
	return nil
}

func (m *MockRoomRepository) GetByID(id string) (*domain.Room, error) {
	room, exists := m.rooms[id]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (m *MockRoomRepository) GetAll() ([]*domain.Room, error) {
	rooms := make([]*domain.Room, 0, len(m.rooms))
	for _, room := range m.rooms {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (m *MockRoomRepository) Exists(id string) bool {
	_, exists := m.rooms[id]
	return exists
}

func TestMessageUsecase_SendMessage(t *testing.T) {
	userRepo := NewMockUserRepository()
	roomRepo := NewMockRoomRepository()
	messageRepo := NewMockMessageRepository()

	user := &domain.User{
		ID:        "user1",
		Username:  "testuser",
		CreatedAt: time.Now(),
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	room := &domain.Room{
		ID:        "room1",
		Name:      "Test Room",
		CreatedAt: time.Now(),
	}
	if err := roomRepo.Create(room); err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	usecase := NewMessageUsecase(messageRepo, userRepo, roomRepo)

	message, err := usecase.SendMessage("room1", "user1", "Hello, world!")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if message.Content != "Hello, world!" {
		t.Errorf("Expected content 'Hello, world!', got %s", message.Content)
	}

	if message.RoomID != "room1" {
		t.Errorf("Expected RoomID 'room1', got %s", message.RoomID)
	}

	_, err = usecase.SendMessage("room1", "user1", "")
	if err == nil {
		t.Fatal("Expected error for empty content, got nil")
	}

	_, err = usecase.SendMessage("room1", "nonexistent", "Hello")
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}

	_, err = usecase.SendMessage("nonexistent", "user1", "Hello")
	if err == nil {
		t.Fatal("Expected error for non-existent room, got nil")
	}
}

func TestMessageUsecase_GetMessagesHistory(t *testing.T) {
	userRepo := NewMockUserRepository()
	roomRepo := NewMockRoomRepository()
	messageRepo := NewMockMessageRepository()

	usecase := NewMessageUsecase(messageRepo, userRepo, roomRepo)

	for i := 0; i < 5; i++ {
		message := &domain.Message{
			ID:        fmt.Sprintf("msg%d", i+1),
			RoomID:    "room1",
			UserID:    "user1",
			Username:  "testuser",
			Content:   "Message",
			CreatedAt: time.Now(),
		}
		if err := messageRepo.Create(message); err != nil {
			t.Fatalf("Failed to create message: %v", err)
		}
	}

	messages, err := usecase.GetMessagesHistory("room1", 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(messages) != 5 {
		t.Errorf("Expected 5 messages, got %d", len(messages))
	}

	limited, err := usecase.GetMessagesHistory("room1", 3, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(limited) != 3 {
		t.Errorf("Expected 3 messages with limit, got %d", len(limited))
	}
}
