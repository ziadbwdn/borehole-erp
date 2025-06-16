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
// CreateProject handles the business logic for creating a new project.
func (s *ProjectServiceImpl) CreateProject(
	ctx context.Context,
	project *models.Project,
	createdBy utils.BinaryUUID,
) (*models.Project, *exception.AppError) {
	project.ID = utils.NewBinaryUUID() // Generate a new UUID for the project
	project.CreatedBy = createdBy
	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	project.UpdatedAt = time.Now()

	appErr := s.projectRepo.Create(ctx, project)
	if appErr != nil {
		return nil, appErr
	}

	return project, nil
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
// It includes a role-based check for updating the 'Status' field.
func (s *ProjectServiceImpl) UpdateProject(
	ctx context.Context,
	project *models.Project,
	userRole models.UserRole, // Added userRole parameter
) *exception.AppError {
	existingProject, appErr := s.projectRepo.GetByID(ctx, project.ID)
	if appErr != nil {
		return appErr
	}

	// --- Granular Authorization Check for Project Status ---
	// If the request includes a status and it's different from the existing status,
	// check if the user has the Engineer role.
	if project.Status != "" && existingProject.Status != project.Status {
		if userRole != models.RoleEngineer {
			return exception.NewPermissionError("Only Engineers are authorized to update project status.")
		}
		existingProject.Status = project.Status // Allow status update for Engineers
	}
	// --- End Granular Authorization Check ---

	// Update other fields if provided. These updates are allowed if the user
	// passed the initial router-level authorization (RoleEngineer for any project update).
	if project.Name != "" {
		existingProject.Name = project.Name
	}
	if project.Description != "" {
		existingProject.Description = project.Description
	}
	if project.Location != "" {
		existingProject.Location = project.Location
	}
	if !project.StartDate.IsZero() {
		existingProject.StartDate = project.StartDate
	}
	if !project.EndDate.IsZero() {
		existingProject.EndDate = project.EndDate
	}
	// The 'Status' field is handled above based on role.

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
