package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-first-application/internal/domain"
	"go-first-application/pkg/response"
	"go-first-application/pkg/validator"
)

type AuthHandler struct {
	authService domain.UserService
}

func NewAuthHandler(authService domain.UserService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}
	if fieldErrors, ok := validator.ValidateStruct(req); !ok {
		response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed", fieldErrors)
		return
	}
	res, err := h.authService.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			response.Error(w, http.StatusConflict, "CONFLICT", "Email is already registered", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user", nil)
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}

	if fieldErrors, ok := validator.ValidateStruct(req); !ok {
		response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed", fieldErrors)
		return
	}

	res, err := h.authService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to authenticate", nil)
		return
	}

	response.JSON(w, http.StatusOK, res)
}
