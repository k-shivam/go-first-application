package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uuid.UUID    `json:"id"`
	ProjectID   uuid.UUID    `json:"project_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title" validate:"required,min=2,max=200"`
	Description string       `json:"description" validate:"max=2000"`
	Status      TaskStatus   `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Priority    TaskPriority `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time   `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string       `json:"title" validate:"required,min=2,max=200"`
	Description string       `json:"description" validate:"max=2000"`
	Status      TaskStatus   `json:"status" validate:"required,oneof=todo in_progress done"`
	Priority    TaskPriority `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time   `json:"due_date"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TaskService interface {
	Create(ctx context.Context, userID, projectID uuid.UUID, req CreateTaskRequest) (*Task, error)
	ListByProject(ctx context.Context, userID, projectID uuid.UUID) ([]Task, error)
	GetByID(ctx context.Context, userID, taskID uuid.UUID) (*Task, error)
	Update(ctx context.Context, userID, taskID uuid.UUID, req UpdateTaskRequest) (*Task, error)
	Delete(ctx context.Context, userID, taskID uuid.UUID) error
}
