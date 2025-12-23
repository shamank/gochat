package handler

import (
	"encoding/json"
	"net/http"

	"gochat/internal/delivery/dto"
	"gochat/internal/usecase"
)

type RoomHandler struct {
	roomUsecase *usecase.RoomUsecase
}

func NewRoomHandler(roomUsecase *usecase.RoomUsecase) *RoomHandler {
	return &RoomHandler{
		roomUsecase: roomUsecase,
	}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("Invalid request body"))
		return
	}

	room, err := h.roomUsecase.CreateRoom(req.Name)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	respondJSON(w, http.StatusCreated, dto.SuccessResponse(room))
}

func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("room id is required"))
		return
	}

	room, err := h.roomUsecase.GetRoom(roomID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, dto.ErrorResponse(err.Error()))
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse(room))
}

func (h *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms, err := h.roomUsecase.GetAllRooms()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse(rooms))
}
