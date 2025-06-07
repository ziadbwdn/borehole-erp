package repository // This MUST be "repository"

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// projectRepository is a concrete implementation of the contract.ProjectRepository interface.
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new instance of ProjectRepository.
func NewProjectRepository(db *gorm.DB) contract.ProjectRepository {
	return &projectRepository{db: db}
}

// Create inserts a new project record into the database.
func (r *projectRepository) Create(ctx context.Context, project *models.Project) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return exception.NewDatabaseError("Failed to create project", err)
	}
	return nil
}

// GetByID retrieves a single project record by its ID.
func (r *projectRepository) GetByID(ctx context.Context, id utils.BinaryUUID) (*models.Project, *exception.AppError) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Project", id.String())
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve project with ID %s", id.String()), err)
	}
	return &project, nil
}

// Update updates an existing project record in the database.
func (r *projectRepository) Update(ctx context.Context, project *models.Project) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(project).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to update project with ID %s", project.ID.String()), err)
	}
	return nil
}

// Delete deletes a project record by its ID.
func (r *projectRepository) Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to delete project with ID %s", id.String()), err)
	}
	return nil
}

// ListByUser retrieves all projects associated with a given user ID.
func (r *projectRepository) ListByUser(ctx context.Context, userID utils.BinaryUUID) ([]*models.Project, *exception.AppError) {
	var projects []*models.Project
	if err := r.db.WithContext(ctx).Where("created_by = ?", userID).Find(&projects).Error; err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve projects for user ID %s", userID.String()), err)
	}
	if len(projects) == 0 {
		return nil, exception.NewNotFoundError("Projects", fmt.Sprintf("for user ID %s", userID.String()))
	}
	return projects, nil
}
