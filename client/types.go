package main

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
