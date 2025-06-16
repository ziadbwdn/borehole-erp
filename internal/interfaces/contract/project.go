package contract

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

// ProjectService defines the interface for project-related business logic.
type ProjectService interface {
	CreateProject(ctx context.Context, project *models.Project, createdBy utils.BinaryUUID) (*models.Project, *exception.AppError)
	GetProject(ctx context.Context, id utils.BinaryUUID) (*models.Project, *exception.AppError)
	// Modified: Added userRole parameter for granular authorization
	UpdateProject(ctx context.Context, project *models.Project, userRole models.UserRole) *exception.AppError
	DeleteProject(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListProjectsByUser(ctx context.Context, userID utils.BinaryUUID) ([]*models.Project, *exception.AppError)
}

// ProjectRepository defines the interface for project data access operations.
type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) *exception.AppError
	GetByID(ctx context.Context, id utils.BinaryUUID) (*models.Project, *exception.AppError)
	Update(ctx context.Context, project *models.Project) *exception.AppError
	Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListByUser(ctx context.Context, userID utils.BinaryUUID) ([]*models.Project, *exception.AppError)
}
