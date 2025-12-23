package repository

import (
	"testing"
	"time"

	"gochat/internal/domain"
)

func TestInMemoryMessageRepository_Create(t *testing.T) {
	repo := NewInMemoryMessageRepository()

	message := &domain.Message{
		ID:        "1",
		RoomID:    "room1",
		UserID:    "user1",
		Username:  "testuser",
		Content:   "Hello, world!",
		CreatedAt: time.Now(),
	}

	err := repo.Create(message)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestInMemoryMessageRepository_GetByRoomID(t *testing.T) {
	repo := NewInMemoryMessageRepository()

	roomID := "room1"
	messages := []*domain.Message{
		{
			ID:        "1",
			RoomID:    roomID,
			UserID:    "user1",
			Username:  "user1",
			Content:   "First message",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "2",
			RoomID:    roomID,
			UserID:    "user2",
			Username:  "user2",
			Content:   "Second message",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "3",
			RoomID:    roomID,
			UserID:    "user1",
			Username:  "user1",
			Content:   "Third message",
			CreatedAt: time.Now(),
		},
	}

	for _, msg := range messages {
		if err := repo.Create(msg); err != nil {
			t.Fatalf("Failed to create message: %v", err)
		}
	}

	retrieved, err := repo.GetByRoomID(roomID, 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(retrieved) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(retrieved))
	}

	if retrieved[0].ID != "1" {
		t.Errorf("Expected first message ID to be '1', got %s", retrieved[0].ID)
	}

	limited, err := repo.GetByRoomID(roomID, 2, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(limited) != 2 {
		t.Errorf("Expected 2 messages with limit, got %d", len(limited))
	}

	offset, err := repo.GetByRoomID(roomID, 10, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(offset) != 2 {
		t.Errorf("Expected 2 messages with offset, got %d", len(offset))
	}
}

func TestInMemoryMessageRepository_GetByID(t *testing.T) {
	repo := NewInMemoryMessageRepository()

	message := &domain.Message{
		ID:        "1",
		RoomID:    "room1",
		UserID:    "user1",
		Username:  "testuser",
		Content:   "Hello, world!",
		CreatedAt: time.Now(),
	}

	if err := repo.Create(message); err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	retrieved, err := repo.GetByID("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != message.ID {
		t.Errorf("Expected ID %s, got %s", message.ID, retrieved.ID)
	}

	if retrieved.Content != message.Content {
		t.Errorf("Expected Content %s, got %s", message.Content, retrieved.Content)
	}
}
