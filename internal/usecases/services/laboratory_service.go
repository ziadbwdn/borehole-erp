package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// NOTE: Please ensure you have helper functions in your gin_helpers package
// to extract Username and IP Address from the context if you need that data.
// Example:
// func GetUsernameFromContext(c context.Context) (string, error) { ... }
// func GetIPAddressFromContext(c context.Context) string { ... }


// LaboratoryServiceImpl implements the contract.LaboratoryService interface.
type LaboratoryServiceImpl struct {
	labRepo         contract.LaboratoryRepository
	stationRepo     contract.StationRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// NewLaboratoryService creates and returns a new instance of LaboratoryServiceImpl.
func NewLaboratoryService(
	labRepo contract.LaboratoryRepository,
	stationRepo contract.StationRepository,
	activityService contract.UserActivityService,
	logger logger.Logger,
) contract.LaboratoryService {
	if labRepo == nil { panic("labRepo must not be nil") }
	if stationRepo == nil { panic("stationRepo must not be nil") }
	if activityService == nil { panic("activityService must not be nil") }
	if logger == nil { panic("logger must not be nil") }

	return &LaboratoryServiceImpl{
		labRepo:         labRepo,
		stationRepo:     stationRepo,
		activityService: activityService,
		logger:          logger,
	}
}

// === LabSample Methods ===

// CREATE
func (s *LaboratoryServiceImpl) CreateSample(ctx context.Context, sample *models.LabSample, userRole models.UserRole, logCtx models.ActivityLogContext) (*models.LabSample, *exception.AppError) {
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician && userRole != models.RoleGeologist {
		return nil, exception.NewPermissionError("User not authorized to create laboratory samples")
	}
	_, appErr := s.stationRepo.GetByID(ctx, sample.StationID); 
	if appErr != nil {
		/* ... error handling ... */ 
		return nil, appErr 
	}
	
	sample.ID = utils.NewBinaryUUID()
	sample.CreatedAt = time.Now()
	sample.UpdatedAt = time.Now()
	
	appErr = s.labRepo.CreateSample(ctx, sample); 
	if appErr != nil { 
		return nil, appErr 
	}

	newValueJSON, _ := json.Marshal(sample); newValueStr := string(newValueJSON)
	sampleIDStr := sample.ID.String(); details := fmt.Sprintf("New Lab Sample '%s' created.", sample.SampleCode); ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeCreateSample, models.ResourceTypeSample, &sampleIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil { 
		s.logger.Error(ctx, "Failed to log CreateSample activity", logErr, logger.Field{Key: "sampleID", Value: sample.ID.String()}) 
	}
	
	return sample, nil
}

// GET
func (s *LaboratoryServiceImpl) GetSample(ctx context.Context, id utils.BinaryUUID) (*models.LabSample, *exception.AppError) { 
	return s.labRepo.GetSampleByID(ctx, id) 
}

// UPDATE
func (s *LaboratoryServiceImpl) UpdateSample(ctx context.Context, sample *models.LabSample, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician && userRole != models.RoleGeologist {
		return exception.NewPermissionError("User not authorized to update laboratory samples")
	}
	existingSample, appErr := s.labRepo.GetSampleByID(ctx, sample.ID); if appErr != nil { return appErr }
	oldValueJSON, err := json.Marshal(existingSample); if err != nil { s.logger.Warn(ctx, "Failed to marshal old lab sample", logger.Field{Key: "error", Value: err.Error()}) }
	oldValueStr := string(oldValueJSON)

	// Update logic from your original file
	if sample.StationID != (utils.BinaryUUID{}) { /* ... */ existingSample.StationID = sample.StationID }
	if sample.SampleCode != "" { existingSample.SampleCode = sample.SampleCode }
	if sample.DepthFrom.Internal.Value != "" { existingSample.DepthFrom = sample.DepthFrom }
	if sample.DepthTo.Internal.Value != "" { existingSample.DepthTo = sample.DepthTo }
	if sample.SampleType != "" { existingSample.SampleType = sample.SampleType }
	if !sample.SamplingDate.IsZero() { existingSample.SamplingDate = sample.SamplingDate }
	if sample.TestedBy != "" { existingSample.TestedBy = sample.TestedBy }
	if sample.LabName != "" { existingSample.LabName = sample.LabName }
	if !sample.TestDate.IsZero() { existingSample.TestDate = sample.TestDate }
	existingSample.UpdatedAt = time.Now()
	
	appErr = s.labRepo.UpdateSample(ctx, existingSample); if appErr != nil { return appErr }
	
	newValueJSON, err := json.Marshal(existingSample); if err != nil { s.logger.Warn(ctx, "Failed to marshal new lab sample", logger.Field{Key: "error", Value: err.Error()}) }
	newValueStr := string(newValueJSON)
	sampleIDStr := existingSample.ID.String(); details := fmt.Sprintf("Lab Sample '%s' updated.", existingSample.SampleCode); ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeUpdateSample, models.ResourceTypeSample, &sampleIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	if logErr != nil { s.logger.Error(ctx, "Failed to log UpdateSample activity", logErr, logger.Field{Key: "sampleID", Value: existingSample.ID.String()}) }

	return nil
}

// DELETE
func (s *LaboratoryServiceImpl) DeleteSample(ctx context.Context, id utils.BinaryUUID, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	if userRole != models.RoleAdmin { return exception.NewPermissionError("User not authorized to delete laboratory samples") }
	sampleToDelete, appErr := s.labRepo.GetSampleByID(ctx, id); if appErr != nil { return appErr }
	oldValueJSON, _ := json.Marshal(sampleToDelete); oldValueStr := string(oldValueJSON)
	
	appErr = s.labRepo.DeleteSample(ctx, id); if appErr != nil { return appErr }

	sampleIDStr := sampleToDelete.ID.String(); details := fmt.Sprintf("Lab Sample '%s' deleted.", sampleToDelete.SampleCode); ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeDeleteSample, models.ResourceTypeSample, &sampleIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil { s.logger.Error(ctx, "Failed to log DeleteSample activity", logErr, logger.Field{Key: "sampleID", Value: id.String()}) }
	
	return nil
}

// SAMPLES BY STATION LIST
func (s *LaboratoryServiceImpl) ListSamplesByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LabSample, *exception.AppError) { 
	return s.labRepo.ListSamplesByStation(ctx, stationID) 
}

// === UCSResult Methods ===

// CREATE
func (s *LaboratoryServiceImpl) CreateUCSResult(ctx context.Context, ucs *models.UCSResult, userRole models.UserRole, logCtx models.ActivityLogContext) (*models.UCSResult, *exception.AppError) {
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician {
		return nil, exception.NewPermissionError("User not authorized to create UCS results")
	}
	_, appErr := s.labRepo.GetSampleByID(ctx, ucs.SampleID); if appErr != nil { /* ... error handling ... */ return nil, appErr }
	
	ucs.ID = utils.NewBinaryUUID(); ucs.CreatedAt = time.Now(); ucs.UpdatedAt = time.Now()
	
	appErr = s.labRepo.CreateUCSResult(ctx, ucs); if appErr != nil { return nil, appErr }
	
	newValueJSON, _ := json.Marshal(ucs); newValueStr := string(newValueJSON)
	ucsIDStr := ucs.ID.String(); details := fmt.Sprintf("New UCS Result created for Sample ID %s.", ucs.SampleID.String()); ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeCreateUCSResult, models.ResourceTypeUCSResult, &ucsIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil { s.logger.Error(ctx, "Failed to log CreateUCSResult activity", logErr, logger.Field{Key: "ucsResultID", Value: ucs.ID.String()}) }
	
	return ucs, nil
}

// GET
func (s *LaboratoryServiceImpl) GetUCSResult(ctx context.Context, id utils.BinaryUUID) (*models.UCSResult, *exception.AppError) { 
	return s.labRepo.GetUCSResultByID(ctx, id) 
}

// UPDATE
func (s *LaboratoryServiceImpl) UpdateUCSResult(ctx context.Context, ucs *models.UCSResult, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	if userRole != models.RoleAdmin && userRole != models.RoleLabTechnician {
		return exception.NewPermissionError("User not authorized to update UCS results")
	}
	existingUCS, appErr := s.labRepo.GetUCSResultByID(ctx, ucs.ID); if appErr != nil { return appErr }
	oldValueJSON, err := json.Marshal(existingUCS); if err != nil { s.logger.Warn(ctx, "Failed to marshal old UCS result", logger.Field{Key: "error", Value: err.Error()}) }
	oldValueStr := string(oldValueJSON)

	// Update logic from your original file
	if ucs.SampleID != (utils.BinaryUUID{}) { /* ... */ existingUCS.SampleID = ucs.SampleID }
	if ucs.UCSValue.Internal.Value != "" { existingUCS.UCSValue = ucs.UCSValue }
	if ucs.Unit != "" { existingUCS.Unit = ucs.Unit }
	if ucs.TestMethod != "" { existingUCS.TestMethod = ucs.TestMethod }
	if ucs.SpecimenDiameter.Internal.Value != "" { existingUCS.SpecimenDiameter = ucs.SpecimenDiameter }
	if ucs.SpecimenHeight.Internal.Value != "" { existingUCS.SpecimenHeight = ucs.SpecimenHeight }
	if ucs.FailureMode != "" { existingUCS.FailureMode = ucs.FailureMode }
	if ucs.Notes != "" { existingUCS.Notes = ucs.Notes }
	existingUCS.UpdatedAt = time.Now()
	
	appErr = s.labRepo.UpdateUCSResult(ctx, existingUCS); if appErr != nil { return appErr }
	
	newValueJSON, err := json.Marshal(existingUCS); if err != nil { s.logger.Warn(ctx, "Failed to marshal new UCS result", logger.Field{Key: "error", Value: err.Error()}) }
	newValueStr := string(newValueJSON)
	ucsIDStr := existingUCS.ID.String(); details := fmt.Sprintf("UCS Result ID %s updated.", existingUCS.ID.String()); ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeUpdateUCSResult, models.ResourceTypeUCSResult, &ucsIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	if logErr != nil { s.logger.Error(ctx, "Failed to log UpdateUCSResult activity", logErr, logger.Field{Key: "ucsResultID", Value: existingUCS.ID.String()}) }
	
	return nil
}

// DELETE
func (s *LaboratoryServiceImpl) DeleteUCSResult(ctx context.Context, id utils.BinaryUUID, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	if userRole != models.RoleAdmin { return exception.NewPermissionError("User not authorized to delete UCS results") }
	ucsToDelete, appErr := s.labRepo.GetUCSResultByID(ctx, id); if appErr != nil { return appErr }
	oldValueJSON, _ := json.Marshal(ucsToDelete); oldValueStr := string(oldValueJSON)
	
	appErr = s.labRepo.DeleteUCSResult(ctx, id); if appErr != nil { return appErr }

	ucsIDStr := ucsToDelete.ID.String(); details := fmt.Sprintf("UCS Result ID %s deleted.", ucsToDelete.ID.String()); ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeDeleteUCSResult, models.ResourceTypeUCSResult, &ucsIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil { s.logger.Error(ctx, "Failed to log DeleteUCSResult activity", logErr, logger.Field{Key: "ucsResultID", Value: id.String()}) }
	
	return nil
}

// UCS RESULT BY SAMPLE LIST
func (s *LaboratoryServiceImpl) ListUCSResultsBySample(ctx context.Context, sampleID utils.BinaryUUID) ([]*models.UCSResult, *exception.AppError) { 
	return s.labRepo.GetUCSResultsBySample(ctx, sampleID) 
}