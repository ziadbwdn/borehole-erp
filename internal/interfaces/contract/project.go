package contract

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

// ProjectRepository defines the contract for project data operations.
type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) *exception.AppError
	GetByID(ctx context.Context, id utils.BinaryUUID) (*models.Project, *exception.AppError)
	Update(ctx context.Context, project *models.Project) *exception.AppError
	Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListByUser(ctx context.Context, userID utils.BinaryUUID) ([]*models.Project, *exception.AppError)
}

// ProjectService defines the contract for project business logic.
type ProjectService interface {
	// Corrected CreateProject signature to match implementation:
	// Now returns (*models.Project, *exception.AppError)
	CreateProject(ctx context.Context, project *models.Project, createdBy utils.BinaryUUID) (*models.Project, *exception.AppError)
	GetProject(ctx context.Context, id utils.BinaryUUID) (*models.Project, *exception.AppError)
	UpdateProject(ctx context.Context, project *models.Project) *exception.AppError
	DeleteProject(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListProjectsByUser(ctx context.Context, userID utils.BinaryUUID) ([]*models.Project, *exception.AppError)
}
