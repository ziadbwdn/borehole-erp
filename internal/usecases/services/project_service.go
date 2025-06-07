package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"time"
)

// ProjectServiceImpl implements the contract.ProjectService interface.
type ProjectServiceImpl struct {
	projectRepo contract.ProjectRepository
}

// NewProjectService creates and returns a new instance of ProjectServiceImpl.
func NewProjectService(projectRepo contract.ProjectRepository) contract.ProjectService {
	return &ProjectServiceImpl{
		projectRepo: projectRepo,
	}
}

// CreateProject handles the business logic for creating a new project.
// This method's signature MUST match the contract.ProjectService interface.
func (s *ProjectServiceImpl) CreateProject(
	ctx context.Context,
	project *models.Project,
	createdBy utils.BinaryUUID,
) (*models.Project, *exception.AppError) { // <-- ENSURE THIS RETURN SIGNATURE IS EXACTLY AS SHOWN
	// Assign the creator and set timestamps if not already set by GORM hooks.
	project.ID = utils.NewBinaryUUID() // Generate a new UUID for the project
	project.CreatedBy = createdBy
	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	project.UpdatedAt = time.Now()

	appErr := s.projectRepo.Create(ctx, project)
	if appErr != nil {
		return nil, appErr // Return nil for project on error
	}

	return project, nil // Return the created project on success
}

// GetProject retrieves a project by its ID.
func (s *ProjectServiceImpl) GetProject(
	ctx context.Context,
	id utils.BinaryUUID,
) (*models.Project, *exception.AppError) {
	project, appErr := s.projectRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	return project, nil
}

// UpdateProject handles the business logic for updating an existing project.
func (s *ProjectServiceImpl) UpdateProject(
	ctx context.Context,
	project *models.Project,
) *exception.AppError {
	existingProject, appErr := s.projectRepo.GetByID(ctx, project.ID)
	if appErr != nil {
		return appErr
	}

	if project.Name != "" {
		existingProject.Name = project.Name
	}
	if project.Description != "" {
		existingProject.Description = project.Description
	}
	if project.Location != "" {
		existingProject.Location = project.Location
	}
	// Use .Value != "" for google.golang.org/genproto/googleapis/type/decimal.Decimal
	if project.Latitude.Internal.Value != "" {
		existingProject.Latitude.Internal.Value = project.Latitude.Internal.Value
	}
	if project.Longitude.Internal.Value != "" {
		existingProject.Longitude.Internal.Value = project.Longitude.Internal.Value
	}
	if project.Elevation.Internal.Value != "" {
		existingProject.Elevation.Internal.Value = project.Elevation.Internal.Value
	}
	if !project.StartDate.IsZero() {
		existingProject.StartDate = project.StartDate
	}
	if !project.EndDate.IsZero() {
		existingProject.EndDate = project.EndDate
	}
	if project.Status != "" {
		existingProject.Status = project.Status
	}

	existingProject.UpdatedAt = time.Now()

	appErr = s.projectRepo.Update(ctx, existingProject)
	if appErr != nil {
		return appErr
	}

	return nil
}

// DeleteProject handles the business logic for deleting a project by its ID.
func (s *ProjectServiceImpl) DeleteProject(
	ctx context.Context,
	id utils.BinaryUUID,
) *exception.AppError {
	appErr := s.projectRepo.Delete(ctx, id)
	if appErr != nil {
		return appErr
	}
	return nil
}

// ListProjectsByUser retrieves a list of projects associated with a specific user.
func (s *ProjectServiceImpl) ListProjectsByUser(
	ctx context.Context,
	userID utils.BinaryUUID,
) ([]*models.Project, *exception.AppError) {
	projects, appErr := s.projectRepo.ListByUser(ctx, userID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return []*models.Project{}, nil
		}
		return nil, appErr
	}
	return projects, nil
}
