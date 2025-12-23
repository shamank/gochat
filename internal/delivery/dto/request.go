package dto

type RegisterUserRequest struct {
	Username string `json:"username"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

type GetMessagesRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
