package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-first-application/internal/domain"
	"go-first-application/internal/middleware"
	"go-first-application/pkg/response"
	"go-first-application/pkg/validator"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService domain.TaskService
}

func NewTaskHandler(taskService domain.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid project ID", nil)
		return
	}

	var req domain.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}

	if fieldErrors, ok := validator.ValidateStruct(req); !ok {
		response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed", fieldErrors)
		return
	}

	task, err := h.taskService.Create(r.Context(), userID, projectID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create task", nil)
		return
	}

	response.JSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid project ID", nil)
		return
	}

	tasks, err := h.taskService.ListByProject(r.Context(), userID, projectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list tasks", nil)
		return
	}

	response.JSON(w, http.StatusOK, response.Envelope{"tasks": tasks})
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid task ID", nil)
		return
	}

	task, err := h.taskService.GetByID(r.Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Task not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get task", nil)
		return
	}

	response.JSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid task ID", nil)
		return
	}

	var req domain.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}

	if fieldErrors, ok := validator.ValidateStruct(req); !ok {
		response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed", fieldErrors)
		return
	}

	task, err := h.taskService.Update(r.Context(), userID, taskID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Task not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update task", nil)
		return
	}

	response.JSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid task ID", nil)
		return
	}

	if err := h.taskService.Delete(r.Context(), userID, taskID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Task not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete task", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
