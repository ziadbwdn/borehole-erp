package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // For BinaryUUID and decimal conversions
	"context"
	"fmt"
	"time"
)

// LaboratoryServiceImpl implements the contract.LaboratoryService interface.
type LaboratoryServiceImpl struct {
	labRepo     contract.LaboratoryRepository
	stationRepo contract.StationRepository // Dependency to verify StationID exists
}

// NewLaboratoryService creates and returns a new instance of LaboratoryServiceImpl.
func NewLaboratoryService(
	labRepo contract.LaboratoryRepository,
	stationRepo contract.StationRepository,
) contract.LaboratoryService {
	return &LaboratoryServiceImpl{
		labRepo:     labRepo,
		stationRepo: stationRepo,
	}
}

// CreateSample handles the business logic for creating a new laboratory sample.
func (s *LaboratoryServiceImpl) CreateSample(
	ctx context.Context,
	sample *models.LabSample,
	userRole models.UserRole,
) (*models.LabSample, *exception.AppError) {
	// --- Authorization Check for Create ---
	// Only 'admin', 'lab_technician', or 'geologist' roles can create lab samples.
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician && userRole != models.RoleGeologist {
		return nil, exception.NewPermissionError("User not authorized to create laboratory samples")
	}
	// --- End Authorization Check ---

	// Verify StationID exists
	_, appErr := s.stationRepo.GetByID(ctx, sample.StationID)
	if appErr != nil {
		// If station is not found, or any other error from stationRepo.GetByID
		// we return a validation error specifically for the station ID.
		if appErr.Code == exception.ErrNotFound {
			return nil, exception.NewValidationError(fmt.Sprintf("Station with ID '%s' not found for sample creation", sample.StationID.String()))
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to verify station with ID '%s'", sample.StationID.String()), appErr.Err)
	}

	// Set new UUID and timestamps
	sample.ID = utils.NewBinaryUUID()
	if sample.CreatedAt.IsZero() {
		sample.CreatedAt = time.Now()
	}
	sample.UpdatedAt = time.Now()

	appErr = s.labRepo.CreateSample(ctx, sample)
	if appErr != nil {
		return nil, appErr // Propagate error from repository
	}

	return sample, nil
}

// GetSample retrieves a laboratory sample by its ID.
func (s *LaboratoryServiceImpl) GetSample(
	ctx context.Context,
	id utils.BinaryUUID,
) (*models.LabSample, *exception.AppError) {
	// Calling the repository method using its actual name.
	sample, appErr := s.labRepo.GetSampleByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	return sample, nil
}

// UpdateSample handles the business logic for updating an existing laboratory sample.
func (s *LaboratoryServiceImpl) UpdateSample(
	ctx context.Context,
	sample *models.LabSample,
	userRole models.UserRole,
) *exception.AppError {
	// --- Authorization Check for Update ---
	// Only 'admin', 'lab_technician', or 'geologist' roles can update lab samples.
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician && userRole != models.RoleGeologist {
		return exception.NewPermissionError("User not authorized to update laboratory samples")
	}
	// --- End Authorization Check ---

	// First, retrieve the existing sample to ensure it exists and to get current values.
	// Calling the repository method using its actual name.
	existingSample, appErr := s.labRepo.GetSampleByID(ctx, sample.ID)
	if appErr != nil {
		return appErr // Propagate NotFoundError or DatabaseError
	}

	// Apply updates from the provided 'sample' model to the 'existingSample'
	if sample.StationID != (utils.BinaryUUID{}) { // Check if StationID is provided (non-zero UUID)
		// Verify new StationID exists if it's being changed
		_, appErr := s.stationRepo.GetByID(ctx, sample.StationID)
		if appErr != nil {
			if appErr.Code == exception.ErrNotFound {
				return exception.NewValidationError(fmt.Sprintf("New Station with ID '%s' not found for sample update", sample.StationID.String()))
			}
			return exception.NewDatabaseError(fmt.Sprintf("Failed to verify new station with ID '%s'", sample.StationID.String()), appErr.Err)
		}
		existingSample.StationID = sample.StationID
	}
	if sample.SampleCode != "" {
		existingSample.SampleCode = sample.SampleCode
	}
	if sample.DepthFrom.Internal.Value != "" {
		existingSample.DepthFrom = sample.DepthFrom
	}
	if sample.DepthTo.Internal.Value != "" {
		existingSample.DepthTo = sample.DepthTo
	}
	if sample.SampleType != "" {
		existingSample.SampleType = sample.SampleType
	}
	if !sample.SamplingDate.IsZero() { // Check if the time is not its zero value
		existingSample.SamplingDate = sample.SamplingDate
	}
	if sample.TestedBy != "" {
		existingSample.TestedBy = sample.TestedBy
	}
	if sample.LabName != "" {
		existingSample.LabName = sample.LabName
	}
	if !sample.TestDate.IsZero() { // Check if the time is not its zero value
		existingSample.TestDate = sample.TestDate
	}

	existingSample.UpdatedAt = time.Now()

	appErr = s.labRepo.UpdateSample(ctx, existingSample)
	if appErr != nil {
		return appErr
	}

	return nil
}

// DeleteSample handles the business logic for deleting a laboratory sample.
func (s *LaboratoryServiceImpl) DeleteSample(
	ctx context.Context,
	id utils.BinaryUUID,
	userRole models.UserRole,
) *exception.AppError {
	// --- Authorization Check for Delete ---
	// Only 'admin' role can delete lab samples.
	if userRole != models.RoleAdmin {
		return exception.NewPermissionError("User not authorized to delete laboratory samples")
	}
	// --- End Authorization Check ---

	appErr := s.labRepo.DeleteSample(ctx, id)
	if appErr != nil {
		return appErr
	}
	return nil
}

// ListSamplesByStation retrieves a list of laboratory samples for a specific station.
func (s *LaboratoryServiceImpl) ListSamplesByStation(
	ctx context.Context,
	stationID utils.BinaryUUID,
) ([]*models.LabSample, *exception.AppError) {
	samples, appErr := s.labRepo.ListSamplesByStation(ctx, stationID)
	if appErr != nil {
		// If the repository now returns nil for NotFound (as recommended below),
		// this check effectively handles other database errors.
		return nil, appErr
	}
	return samples, nil
}

// CreateUCSResult handles the business logic for creating a new UCS result.
func (s *LaboratoryServiceImpl) CreateUCSResult(
	ctx context.Context,
	ucs *models.UCSResult,
	userRole models.UserRole,
) (*models.UCSResult, *exception.AppError) {
	// --- Authorization Check for Create ---
	// Only 'admin' or 'lab_technician' roles can create UCS results.
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician {
		return nil, exception.NewPermissionError("User not authorized to create UCS results")
	}
	// --- End Authorization Check ---

	// Verify SampleID exists
	// Calling the repository method using its actual name.
	_, appErr := s.labRepo.GetSampleByID(ctx, ucs.SampleID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return nil, exception.NewValidationError(fmt.Sprintf("Lab Sample with ID '%s' not found for UCS result creation", ucs.SampleID.String()))
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to verify lab sample with ID '%s'", ucs.SampleID.String()), appErr.Err)
	}

	// Set new UUID and timestamps
	ucs.ID = utils.NewBinaryUUID()
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
func (s *LaboratoryServiceImpl) GetUCSResult(
	ctx context.Context,
	id utils.BinaryUUID,
) (*models.UCSResult, *exception.AppError) {
	// Calling the repository method using its actual name.
	ucs, appErr := s.labRepo.GetUCSResultByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	return ucs, nil
}

// UpdateUCSResult handles the business logic for updating an existing UCS result.
func (s *LaboratoryServiceImpl) UpdateUCSResult(
	ctx context.Context,
	ucs *models.UCSResult,
	userRole models.UserRole,
) *exception.AppError {
	// --- Authorization Check for Update ---
	// Only 'admin' or 'lab_technician' roles can update UCS results.
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician {
		return exception.NewPermissionError("User not authorized to update UCS results")
	}
	// --- End Authorization Check ---

	// First, retrieve the existing UCS result.
	// Calling the repository method using its actual name.
	existingUCS, appErr := s.labRepo.GetUCSResultByID(ctx, ucs.ID)
	if appErr != nil {
		return appErr
	}

	// Apply updates
	if ucs.SampleID != (utils.BinaryUUID{}) {
		// Verify new SampleID exists if it's being changed
		// Calling the repository method using its actual name.
		_, appErr := s.labRepo.GetSampleByID(ctx, ucs.SampleID)
		if appErr != nil {
			if appErr.Code == exception.ErrNotFound {
				return exception.NewValidationError(fmt.Sprintf("New Lab Sample with ID '%s' not found for UCS result update", ucs.SampleID.String()))
			}
			return exception.NewDatabaseError(fmt.Sprintf("Failed to verify new lab sample with ID '%s'", ucs.SampleID.String()), appErr.Err)
		}
		existingUCS.SampleID = ucs.SampleID
	}
	if ucs.UCSValue.Internal.Value != "" {
		existingUCS.UCSValue = ucs.UCSValue
	}
	if ucs.Unit != "" {
		existingUCS.Unit = ucs.Unit
	}
	if ucs.TestMethod != "" {
		existingUCS.TestMethod = ucs.TestMethod
	}
	if ucs.SpecimenDiameter.Internal.Value != "" {
		existingUCS.SpecimenDiameter = ucs.SpecimenDiameter
	}
	if ucs.SpecimenHeight.Internal.Value != "" {
		existingUCS.SpecimenHeight = ucs.SpecimenHeight
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
func (s *LaboratoryServiceImpl) DeleteUCSResult(
	ctx context.Context,
	id utils.BinaryUUID,
	userRole models.UserRole,
) *exception.AppError {
	// --- Authorization Check for Delete ---
	// Only 'admin' role can delete UCS results.
	if userRole != models.RoleAdmin {
		return exception.NewPermissionError("User not authorized to delete UCS results")
	}
	// --- End Authorization Check ---

	appErr := s.labRepo.DeleteUCSResult(ctx, id)
	if appErr != nil {
		return appErr
	}
	return nil
}

// ListUCSResultsBySample retrieves a list of UCS results for a specific lab sample.
func (s *LaboratoryServiceImpl) ListUCSResultsBySample(
	ctx context.Context,
	sampleID utils.BinaryUUID,
) ([]*models.UCSResult, *exception.AppError) {
	ucsResults, appErr := s.labRepo.GetUCSResultsBySample(ctx, sampleID) // Calling correct repo method
	if appErr != nil {
		// If the repository now returns nil for NotFound (as it already does for this method),
		// this check effectively handles other database errors.
		return nil, appErr
	}
	return ucsResults, nil
}
