package service

import (
	"context"
	"time"

	"go-first-application/internal/domain"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
}

func NewTaskService(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, userID, projectID uuid.UUID, req domain.CreateTaskRequest) (*domain.Task, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrForbidden
	}

	status := req.Status
	if status == "" {
		status = domain.TaskStatusTodo
	}

	priority := req.Priority
	if priority == "" {
		priority = domain.TaskPriorityMedium
	}

	now := time.Now().UTC()
	task := &domain.Task{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ListByProject(ctx context.Context, userID, projectID uuid.UUID) ([]domain.Task, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return s.taskRepo.ListByProjectID(ctx, projectID)
}

func (s *TaskService) GetByID(ctx context.Context, userID, taskID uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return task, nil
}

func (s *TaskService) Update(ctx context.Context, userID, taskID uuid.UUID, req domain.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.GetByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.Priority = req.Priority
	task.DueDate = req.DueDate
	task.UpdatedAt = time.Now().UTC()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, userID, taskID uuid.UUID) error {
	_, err := s.GetByID(ctx, userID, taskID)
	if err != nil {
		return err
	}
	return s.taskRepo.Delete(ctx, taskID)
}
