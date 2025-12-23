package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gochat/internal/delivery/dto"
	"gochat/internal/delivery/websocket"
	"gochat/internal/usecase"
)

type MessageHandler struct {
	messageUsecase *usecase.MessageUsecase
	wsHub          *websocket.Hub
}

func NewMessageHandler(
	messageUsecase *usecase.MessageUsecase,
	wsHub *websocket.Hub,
) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
		wsHub:          wsHub,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	userID := r.URL.Query().Get("user_id")

	if roomID == "" || userID == "" {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("room_id and user_id are required"))
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("Invalid request body"))
		return
	}

	message, err := h.messageUsecase.SendMessage(roomID, userID, req.Content)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	h.wsHub.BroadcastMessage(roomID, message)

	respondJSON(w, http.StatusCreated, dto.SuccessResponse(message))
}

func (h *MessageHandler) GetMessagesHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("room_id is required"))
		return
	}

	const defaultLimit = 50

	limit := defaultLimit
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	messages, err := h.messageUsecase.GetMessagesHistory(roomID, limit, offset)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse(messages))
}
