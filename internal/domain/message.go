package domain

import "time"

type Message struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRepository interface {
	Create(message *Message) error
	GetByRoomID(roomID string, limit, offset int) ([]*Message, error)
	GetByID(id string) (*Message, error)
}
