package handler

import (
	"encoding/json"
	"net/http"

	"gochat/internal/delivery/dto"
	"gochat/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("Invalid request body"))
		return
	}

	user, err := h.userUsecase.RegisterUser(req.Username)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	respondJSON(w, http.StatusCreated, dto.SuccessResponse(user))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse("user id is required"))
		return
	}

	user, err := h.userUsecase.GetUser(userID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, dto.ErrorResponse(err.Error()))
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse(user))
}
