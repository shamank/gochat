package domain

import "time"

type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RoomRepository interface {
	Create(room *Room) error
	GetByID(id string) (*Room, error)
	GetAll() ([]*Room, error)
	Exists(id string) bool
}
