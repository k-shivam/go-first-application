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

type ProjectHandler struct {
	projectService domain.ProjectService
}

func NewProjectHandler(projectService domain.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	var req domain.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}

	if fieldErrors, ok := validator.ValidateStruct(req); !ok {
		response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed", fieldErrors)
		return
	}

	project, err := h.projectService.Create(r.Context(), userID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create project", nil)
		return
	}

	response.JSON(w, http.StatusCreated, project)

}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	projects, err := h.projectService.List(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list projects", nil)
		return
	}

	response.JSON(w, http.StatusOK, response.Envelope{"projects": projects})

}

func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid project ID", nil)
		return
	}

	project, err := h.projectService.GetByID(r.Context(), userID, projectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied to project", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get project", nil)
		return
	}

	response.JSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid project ID", nil)
		return
	}

	var req domain.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}

	if fieldErrors, ok := validator.ValidateStruct(req); !ok {
		response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed", fieldErrors)
		return
	}

	project, err := h.projectService.Update(r.Context(), userID, projectID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Access denied to project", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update project", nil)
		return
	}

	response.JSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User identity missing", nil)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid project ID", nil)
		return
	}

	if err := h.projectService.Delete(r.Context(), userID, projectID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Project not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete project", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
