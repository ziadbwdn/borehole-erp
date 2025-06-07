package repository

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type laboratoryRepository struct {
	db *gorm.DB
}

// NewLaboratoryRepository creates a new instance of LaboratoryRepository.
func NewLaboratoryRepository(db *gorm.DB) contract.LaboratoryRepository {
	return &laboratoryRepository{db: db}
}

// CreateSample creates a new lab sample in the database.
func (r *laboratoryRepository) CreateSample(ctx context.Context, sample *models.LabSample) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(sample).Error; err != nil {
		return exception.NewDatabaseError("Failed to create lab sample", err)
	}
	return nil
}

// GetSampleByID retrieves a lab sample by its ID.
func (r *laboratoryRepository) GetSampleByID(ctx context.Context, id utils.BinaryUUID) (*models.LabSample, *exception.AppError) {
	var sample models.LabSample
	if err := r.db.WithContext(ctx).First(&sample, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Lab sample", id.String())
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve lab sample with ID %s", id.String()), err)
	}
	return &sample, nil
}

// UpdateSample updates an existing lab sample in the database.
func (r *laboratoryRepository) UpdateSample(ctx context.Context, sample *models.LabSample) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(sample).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to update lab sample with ID %s", sample.ID.String()), err)
	}
	return nil
}

// DeleteSample deletes a lab sample by its ID.
func (r *laboratoryRepository) DeleteSample(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&models.LabSample{}, "id = ?", id).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to delete lab sample with ID %s", id.String()), err)
	}
	return nil
}

// ListSamplesByStation retrieves all lab samples associated with a given station ID.
func (r *laboratoryRepository) ListSamplesByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LabSample, *exception.AppError) {
	var samples []*models.LabSample
	if err := r.db.WithContext(ctx).Where("station_id = ?", stationID).Find(&samples).Error; err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve lab samples for station ID %s", stationID.String()), err)
	}
	if len(samples) == 0 {
		return nil, exception.NewNotFoundError("Lab samples", fmt.Sprintf("for station ID %s", stationID.String()))
	}
	return samples, nil
}

// CreateUCSResult creates a new UCS result in the database.
func (r *laboratoryRepository) CreateUCSResult(ctx context.Context, ucs *models.UCSResult) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(ucs).Error; err != nil {
		return exception.NewDatabaseError("Failed to create UCS result", err)
	}
	return nil
}

// GetUCSResultByID retrieves a UCS result by its ID.
func (r *laboratoryRepository) GetUCSResultByID(ctx context.Context, id utils.BinaryUUID) (*models.UCSResult, *exception.AppError) {
	var ucs models.UCSResult
	if err := r.db.WithContext(ctx).First(&ucs, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("UCS result", id.String())
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve UCS result with ID %s", id.String()), err)
	}
	return &ucs, nil
}

// UpdateUCSResult updates an existing UCS result in the database.
func (r *laboratoryRepository) UpdateUCSResult(ctx context.Context, ucs *models.UCSResult) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(ucs).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to update UCS result with ID %s", ucs.ID.String()), err)
	}
	return nil
}

// DeleteUCSResult deletes a UCS result by its ID.
func (r *laboratoryRepository) DeleteUCSResult(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&models.UCSResult{}, "id = ?", id).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to delete UCS result with ID %s", id.String()), err)
	}
	return nil
}

// GetUCSResultsBySample retrieves all UCS results associated with a given sample ID.
func (r *laboratoryRepository) GetUCSResultsBySample(ctx context.Context, sampleID utils.BinaryUUID) ([]*models.UCSResult, *exception.AppError) {
	var ucsResults []*models.UCSResult
	// Find all UCS results where sample_id matches
	if err := r.db.WithContext(ctx).Where("sample_id = ?", sampleID).Find(&ucsResults).Error; err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve UCS results for sample ID %s", sampleID.String()), err)
	}
	// If no records are found, return an empty slice and no error, as it's not an error condition
	// to have no UCS results for a sample, but rather an empty set.
	if len(ucsResults) == 0 {
		return []*models.UCSResult{}, nil // Return empty slice, not NotFoundError, as per reporting needs
	}
	return ucsResults, nil
}
