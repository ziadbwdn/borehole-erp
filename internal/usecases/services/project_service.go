package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger" // <-- Added
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"encoding/json" // <-- Added
	"fmt"           // <-- Added
	"time"
)

// ProjectServiceImpl implements the contract.ProjectService interface.
// REFACTORED: Added activityService and logger fields.
type ProjectServiceImpl struct {
	projectRepo     contract.ProjectRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// NewProjectService creates and returns a new instance of ProjectServiceImpl.
// REFACTORED: Now accepts UserActivityService and Logger.
func NewProjectService(
	projectRepo contract.ProjectRepository,
	activityService contract.UserActivityService,
	logger logger.Logger,
) contract.ProjectService {
	// Panic checks as per the guide's best practices
	if projectRepo == nil {
		panic("projectRepo must not be nil for ProjectService")
	}
	if activityService == nil {
		panic("activityService must not be nil for ProjectService")
	}
	if logger == nil {
		panic("logger must not be nil for ProjectService")
	}
	return &ProjectServiceImpl{
		projectRepo:     projectRepo,
		activityService: activityService,
		logger:          logger,
	}
}

// CreateProject handles the business logic for creating a new project.
// REFACTORED: Added user activity logging on success.
func (s *ProjectServiceImpl) CreateProject(ctx context.Context, project *models.Project, createdBy utils.BinaryUUID, logCtx models.ActivityLogContext) (*models.Project, *exception.AppError) {
	project.ID = utils.NewBinaryUUID()
	project.CreatedBy = createdBy
	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	project.UpdatedAt = time.Now()

	appErr := s.projectRepo.Create(ctx, project)
	if appErr != nil {
		return nil, appErr
	}

	// --- LOG USER ACTIVITY ---
	newValueJSON, _ := json.Marshal(project)
	newValueStr := string(newValueJSON)
	projectIDStr := project.ID.String()
	details := fmt.Sprintf("New Project '%s' created.", project.Name)
	ipAddr := logCtx.IPAddress // Use local var to take its address

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeCreateProject, models.ResourceTypeProject, &projectIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log CreateProject activity", logErr, logger.Field{Key: "projectID", Value: project.ID.String()})
	}

	return project, nil
}

// GetProject retrieves a project by its ID.
// REFACTORED: Added optional (commented-out) logging for reads.
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
func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, project *models.Project, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	// Fetching current record from the db
	existingProject, appErr := s.projectRepo.GetByID(ctx, project.ID)
	if appErr != nil {
		return appErr
	}

	// Marshal the original state for logging before any changes are made.
	oldValueJSON, err := json.Marshal(existingProject)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal old project value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "projectID", Value: project.ID.String()})
	}
	oldValueStr := string(oldValueJSON)

	// Apply the requested changes from the input 'project' model to 'existingProject'.
	
	// Special authorization for Status
	if project.Status != "" && existingProject.Status != project.Status {
		if userRole != models.RoleEngineer {
			return exception.NewPermissionError("Only Engineers are authorized to update project status.")
		}
		existingProject.Status = project.Status
	}

	// CORRECTED: Apply all other fields if they were provided in the request.
	if project.Name != "" {
		existingProject.Name = project.Name
	}
	if project.Description != "" {
		existingProject.Description = project.Description
	}
	if project.Location != "" {
		existingProject.Location = project.Location
	}
	// The handler logic ensures that StartDate/EndDate are only non-zero if they were in the request.
	if project.StartDate != nil && !project.StartDate.IsZero() {
		existingProject.StartDate = project.StartDate
	}
	if project.EndDate != nil && !project.EndDate.IsZero() {
		existingProject.EndDate = project.EndDate
	}

	existingProject.UpdatedAt = time.Now()
	appErr = s.projectRepo.Update(ctx, existingProject)
	if appErr != nil {
		return appErr
	}

	// Log the activity with the old and new values.
	newValueJSON, err := json.Marshal(existingProject)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal new project value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "projectID", Value: project.ID.String()})
	}
	newValueStr := string(newValueJSON)
	projectIDStr := existingProject.ID.String()
	details := fmt.Sprintf("Project '%s' updated.", existingProject.Name)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeUpdateProject, models.ResourceTypeProject, &projectIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log UpdateProject activity", logErr, logger.Field{Key: "projectID", Value: existingProject.ID.String()})
	}

	return nil
}

// DeleteProject handles the business logic for deleting a project by its ID.
func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, id utils.BinaryUUID, logCtx models.ActivityLogContext) *exception.AppError {
	projectToDelete, appErr := s.projectRepo.GetByID(ctx, id)
	if appErr != nil {
		return appErr
	}
	oldValueJSON, _ := json.Marshal(projectToDelete)
	oldValueStr := string(oldValueJSON)

	appErr = s.projectRepo.Delete(ctx, id)
	if appErr != nil {
		return appErr
	}

	// --- LOG USER ACTIVITY ---
	projectIDStr := projectToDelete.ID.String()
	details := fmt.Sprintf("Project '%s' deleted.", projectToDelete.Name)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeDeleteProject, models.ResourceTypeProject, &projectIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log DeleteProject activity", logErr, logger.Field{Key: "projectID", Value: projectToDelete.ID.String()})
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