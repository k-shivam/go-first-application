package service

import (
	"context"
	"time"

	"go-first-application/internal/domain"

	"github.com/google/uuid"
)

type ProjectService struct {
	projectRepo domain.ProjectRepository
}

func NewProjectService(projectRepo domain.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateProjectRequest) (*domain.Project, error) {
	now := time.Now().UTC()
	project := &domain.Project{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) GetByID(ctx context.Context, userID, projectID uuid.UUID) (*domain.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return project, nil
}

func (s *ProjectService) List(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	return s.projectRepo.ListByUserID(ctx, userID)
}

func (s *ProjectService) Update(ctx context.Context, userID, projectID uuid.UUID, req domain.UpdateProjectRequest) (*domain.Project, error) {
	project, err := s.GetByID(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	project.Title = req.Title
	project.Description = req.Description
	project.UpdatedAt = time.Now().UTC()

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, userID, projectID uuid.UUID) error {
	return s.projectRepo.Delete(ctx, projectID, userID)
}
