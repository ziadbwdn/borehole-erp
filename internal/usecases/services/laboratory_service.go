package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"fmt"
	"time"
)

type laboratoryService struct {
	labRepo     contract.LaboratoryRepository
	stationRepo contract.StationRepository // Dependency to verify StationID exists for samples
}

// NewLaboratoryService creates a new instance of LaboratoryService.
func NewLaboratoryService(
	labRepo contract.LaboratoryRepository,
	stationRepo contract.StationRepository,
) contract.LaboratoryService {
	return &laboratoryService{
		labRepo:     labRepo,
		stationRepo: stationRepo,
	}
}

// CreateSample handles the business logic for creating a new lab sample.
func (s *laboratoryService) CreateSample(ctx context.Context, sample *models.LabSample) (*models.LabSample, *exception.AppError) {
	// Verify StationID exists
	_, appErr := s.stationRepo.GetByID(ctx, sample.StationID)
	if appErr != nil {
		return nil, exception.NewValidationError(fmt.Sprintf("Station with ID '%s' not found or inaccessible", sample.StationID.String()))
	}

	sample.ID = utils.NewBinaryUUID() // Generate a new UUID for the sample
	if sample.CreatedAt.IsZero() {
		sample.CreatedAt = time.Now()
	}
	sample.UpdatedAt = time.Now()

	appErr = s.labRepo.CreateSample(ctx, sample)
	if appErr != nil {
		return nil, appErr
	}
	return sample, nil
}

// GetSample retrieves a lab sample by its ID.
func (s *laboratoryService) GetSample(ctx context.Context, id utils.BinaryUUID) (*models.LabSample, *exception.AppError) {
	sample, appErr := s.labRepo.GetSampleByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	return sample, nil
}

// UpdateSample handles the business logic for updating an existing lab sample.
func (s *laboratoryService) UpdateSample(ctx context.Context, sample *models.LabSample) *exception.AppError {
	existingSample, appErr := s.labRepo.GetSampleByID(ctx, sample.ID)
	if appErr != nil {
		return appErr // Propagate NotFoundError or DatabaseError
	}

	// Apply updates from the provided 'sample' model to 'existingSample'
	if sample.StationID.String() != (utils.BinaryUUID{}).String() { // Check if StationID is provided (not zero UUID)
		// Verify StationID exists if it's being updated
		_, appErr = s.stationRepo.GetByID(ctx, sample.StationID)
		if appErr != nil {
			return exception.NewValidationError(fmt.Sprintf("Station with ID '%s' not found or inaccessible", sample.StationID.String()))
		}
		existingSample.StationID = sample.StationID
	}
	if sample.SampleCode != "" {
		existingSample.SampleCode = sample.SampleCode
	}
	if sample.DepthFrom.Internal.Value != "" {
		existingSample.DepthFrom.Internal.Value = sample.DepthFrom.Internal.Value
	}
	if sample.DepthTo.Internal.Value != "" {
		existingSample.DepthTo.Internal.Value = sample.DepthTo.Internal.Value
	}
	if sample.SampleType != "" {
		existingSample.SampleType = sample.SampleType
	}
	if !sample.SamplingDate.IsZero() {
        existingSample.SamplingDate = sample.SamplingDate
    }
	if sample.TestedBy != "" {
		existingSample.TestedBy = sample.TestedBy
	}
	if sample.LabName != "" {
		existingSample.LabName = sample.LabName
	}
	if !sample.TestDate.IsZero() {
        existingSample.TestDate = sample.TestDate
    }

	existingSample.UpdatedAt = time.Now()

	appErr = s.labRepo.UpdateSample(ctx, existingSample)
	if appErr != nil {
		return appErr
	}
	return nil
}

// DeleteSample handles the business logic for deleting a lab sample.
func (s *laboratoryService) DeleteSample(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	appErr := s.labRepo.DeleteSample(ctx, id)
	if appErr != nil {
		return appErr
	}
	return nil
}

// ListSamplesByStation retrieves all lab samples for a given station.
func (s *laboratoryService) ListSamplesByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LabSample, *exception.AppError) {
	samples, appErr := s.labRepo.ListSamplesByStation(ctx, stationID)
	if appErr != nil {
		// If no samples are found, the repo might return ErrNotFound.
		// We want to return an empty slice in this case, not an error.
		if appErr.Code == exception.ErrNotFound {
			return []*models.LabSample{}, nil
		}
		return nil, appErr
	}
	return samples, nil
}

// CreateUCSResult handles the business logic for creating a new UCS result.
func (s *laboratoryService) CreateUCSResult(ctx context.Context, ucs *models.UCSResult) (*models.UCSResult, *exception.AppError) {
	// Verify SampleID exists
	_, appErr := s.labRepo.GetSampleByID(ctx, ucs.SampleID)
	if appErr != nil {
		return nil, exception.NewValidationError(fmt.Sprintf("Lab sample with ID '%s' not found or inaccessible", ucs.SampleID.String()))
	}

	ucs.ID = utils.NewBinaryUUID() // Generate a new UUID for the UCS result
	if ucs.CreatedAt.IsZero() {
		ucs.CreatedAt = time.Now()
	}
	ucs.UpdatedAt = time.Now()

	appErr = s.labRepo.CreateUCSResult(ctx, ucs)
	if appErr != nil {
		return nil, appErr
	}
	return ucs, nil
}

// GetUCSResult retrieves a UCS result by its ID.
func (s *laboratoryService) GetUCSResult(ctx context.Context, id utils.BinaryUUID) (*models.UCSResult, *exception.AppError) {
	ucs, appErr := s.labRepo.GetUCSResultByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	return ucs, nil
}

// UpdateUCSResult handles the business logic for updating an existing UCS result.
func (s *laboratoryService) UpdateUCSResult(ctx context.Context, ucs *models.UCSResult) *exception.AppError {
	existingUCS, appErr := s.labRepo.GetUCSResultByID(ctx, ucs.ID)
	if appErr != nil {
		return appErr // Propagate NotFoundError or DatabaseError
	}

	// Apply updates from the provided 'ucs' model to 'existingUCS'
	if ucs.SampleID.String() != (utils.BinaryUUID{}).String() { // Check if SampleID is provided
		// Verify SampleID exists if it's being updated
		_, appErr = s.labRepo.GetSampleByID(ctx, ucs.SampleID)
		if appErr != nil {
			return exception.NewValidationError(fmt.Sprintf("Lab sample with ID '%s' not found or inaccessible", ucs.SampleID.String()))
		}
		existingUCS.SampleID = ucs.SampleID
	}
	if ucs.UCSValue.Internal.Value != "" {
		existingUCS.UCSValue.Internal.Value = ucs.UCSValue.Internal.Value
	}
	if ucs.Unit != "" {
		existingUCS.Unit = ucs.Unit
	}
	if ucs.TestMethod != "" {
		existingUCS.TestMethod = ucs.TestMethod
	}
	if ucs.SpecimenDiameter.Internal.Value != "" {
		existingUCS.SpecimenDiameter.Internal.Value = ucs.SpecimenDiameter.Internal.Value
	}
	if ucs.SpecimenHeight.Internal.Value != "" {
		existingUCS.SpecimenHeight.Internal.Value = ucs.SpecimenHeight.Internal.Value
	}
	if ucs.FailureMode != "" {
		existingUCS.FailureMode = ucs.FailureMode
	}
	if ucs.Notes != "" {
		existingUCS.Notes = ucs.Notes
	}

	existingUCS.UpdatedAt = time.Now()

	appErr = s.labRepo.UpdateUCSResult(ctx, existingUCS)
	if appErr != nil {
		return appErr
	}
	return nil
}

// DeleteUCSResult handles the business logic for deleting a UCS result.
func (s *laboratoryService) DeleteUCSResult(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	appErr := s.labRepo.DeleteUCSResult(ctx, id)
	if appErr != nil {
		return appErr
	}
	return nil
}

// ListUCSResultsBySample retrieves all UCS results for a given lab sample.
func (s *laboratoryService) ListUCSResultsBySample(ctx context.Context, sampleID utils.BinaryUUID) ([]*models.UCSResult, *exception.AppError) {
	// Corrected: Call the new GetUCSResultsBySample from the repository
	ucsResults, appErr := s.labRepo.GetUCSResultsBySample(ctx, sampleID)
	if appErr != nil {
		// If no UCS results are found, the repo should return an empty slice and nil error.
		// If it returns an AppError, propagate it.
		return nil, appErr
	}
	return ucsResults, nil
}
