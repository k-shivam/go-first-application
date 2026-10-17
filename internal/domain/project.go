package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Title       string `json:"title" validate:"required,min=2,max=150"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateProjectRequest struct {
	Title       string `json:"title" validate:"required,min=2,max=150"`
	Description string `json:"description" validate:"max=1000"`
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type ProjectService interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateProjectRequest) (*Project, error)
	GetByID(ctx context.Context, userID, projectID uuid.UUID) (*Project, error)
	List(ctx context.Context, userID uuid.UUID) ([]Project, error)
	Update(ctx context.Context, userID, projectID uuid.UUID, req UpdateProjectRequest) (*Project, error)
	Delete(ctx context.Context, userID, projectID uuid.UUID) error
}
