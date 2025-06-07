package handler

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // For BinaryUUID, StringToGormDecimal, GormDecimalToString
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LaboratoryHandler handles HTTP requests related to laboratory samples and UCS results.
type LaboratoryHandler struct {
	labService contract.LaboratoryService // Dependency on the laboratory service
}

// NewLaboratoryHandler creates and returns a new instance of LaboratoryHandler.
func NewLaboratoryHandler(labService contract.LaboratoryService) *LaboratoryHandler {
	return &LaboratoryHandler{
		labService: labService,
	}
}

// CreateLabSample handles the creation of a new laboratory sample.
// @Router /api/lab-samples [post]
func (h *LaboratoryHandler) CreateLabSample(c *gin.Context) {
	var req dto.CreateLabSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Extract createdBy (userID) from context, set by authentication middleware
	createdBy, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}
	_ = createdBy // Use blank identifier to suppress "declared and not used" warning for now

	// Convert string decimal values from DTO to utils.GormDecimal
	depthFromGd, appErr := utils.StringToGormDecimal(req.DepthFrom)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	depthToGd, appErr := utils.StringToGormDecimal(req.DepthTo)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Convert optional dates:
	var samplingDate time.Time
	if req.SamplingDate != nil {
		samplingDate = *req.SamplingDate
	}
	var testDate time.Time
	if req.TestDate != nil {
		testDate = *req.TestDate
	}

	// Map DTO to model
	sample := &models.LabSample{
		StationID:    req.StationID,
		SampleCode:   req.SampleCode,
		DepthFrom:    *depthFromGd,
		DepthTo:      *depthToGd,
		SampleType:   req.SampleType,
		SamplingDate: samplingDate,
		TestedBy:     req.TestedBy,
		LabName:      req.LabName,
		TestDate:     testDate,
	}

	createdSample, appErr := h.labService.CreateSample(c.Request.Context(), sample)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// When mapping the created model back to the response DTO,
	// convert value types back to pointers.
	var samplingDatePtr *time.Time
	if !createdSample.SamplingDate.IsZero() {
		samplingDatePtr = &createdSample.SamplingDate
	}
	var testDatePtr *time.Time
	if !createdSample.TestDate.IsZero() {
		testDatePtr = &createdSample.TestDate
	}

	// Map created model back to response DTO
	resp := &dto.LabSampleResponse{
		ID:           createdSample.ID,
		StationID:    createdSample.StationID,
		SampleCode:   createdSample.SampleCode,
		DepthFrom:    utils.GormDecimalToString(&createdSample.DepthFrom),
		DepthTo:      utils.GormDecimalToString(&createdSample.DepthTo),
		SampleType:   createdSample.SampleType,
		SamplingDate: samplingDatePtr,
		TestedBy:     createdSample.TestedBy,
		LabName:      createdSample.LabName,
		TestDate:     testDatePtr,
		CreatedAt:    createdSample.CreatedAt,
		UpdatedAt:    createdSample.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetLabSampleByID handles retrieving a single lab sample by ID.
// @Router /api/lab-samples/{id} [get]
func (h *LaboratoryHandler) GetLabSampleByID(c *gin.Context) {
	sampleID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lab sample")
	if appErr != nil {
		return // Response already handled by ParseIDFromContext
	}

	sample, appErr := h.labService.GetSample(c.Request.Context(), sampleID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	var samplingDatePtr *time.Time
	if !sample.SamplingDate.IsZero() {
		samplingDatePtr = &sample.SamplingDate
	}
	var testDatePtr *time.Time
	if !sample.TestDate.IsZero() {
		testDatePtr = &sample.TestDate
	}

	// Map model to response DTO
	resp := &dto.LabSampleResponse{
		ID:           sample.ID,
		StationID:    sample.StationID,
		SampleCode:   sample.SampleCode,
		DepthFrom:    utils.GormDecimalToString(&sample.DepthFrom),
		DepthTo:      utils.GormDecimalToString(&sample.DepthTo),
		SampleType:   sample.SampleType,
		SamplingDate: samplingDatePtr,
		TestedBy:     sample.TestedBy,
		LabName:      sample.LabName,
		TestDate:     testDatePtr,
		CreatedAt:    sample.CreatedAt,
		UpdatedAt:    sample.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateLabSample handles updating an existing laboratory sample.
// @Router /api/lab-samples/{id} [put]
func (h *LaboratoryHandler) UpdateLabSample(c *gin.Context) {
	sampleID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lab sample")
	if appErr != nil {
		return // Response already handled
	}

	var req dto.UpdateLabSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Create a models.LabSample to pass to the service, applying only provided fields
	sampleToUpdate := &models.LabSample{
		ID: sampleID,
	}

	if req.StationID != nil {
		sampleToUpdate.StationID = *req.StationID
	}
	if req.SampleCode != nil {
		sampleToUpdate.SampleCode = *req.SampleCode
	}
	if req.DepthFrom != nil {
		dfGd, err := utils.StringToGormDecimal(*req.DepthFrom)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		sampleToUpdate.DepthFrom = *dfGd
	}
	if req.DepthTo != nil {
		dtGd, err := utils.StringToGormDecimal(*req.DepthTo)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		sampleToUpdate.DepthTo = *dtGd
	}
	if req.SampleType != nil {
		sampleToUpdate.SampleType = *req.SampleType
	}

	// Safely handle optional SamplingDate and TestDate.
	// Check if the pointer itself is not nil before dereferencing.
	if req.SamplingDate != nil {

		if !req.SamplingDate.IsZero() {
			sampleToUpdate.SamplingDate = *req.SamplingDate
		} else {

		}
	}

	if req.TestedBy != nil {
		sampleToUpdate.TestedBy = *req.TestedBy
	}
	if req.LabName != nil {
		sampleToUpdate.LabName = *req.LabName
	}

	if req.TestDate != nil {
		// Same logic as SamplingDate
		if !req.TestDate.IsZero() {
			sampleToUpdate.TestDate = *req.TestDate
		}
	}

	appErr = h.labService.UpdateSample(c.Request.Context(), sampleToUpdate)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Retrieve the updated sample to return the full, current state
	updatedSample, appErr := h.labService.GetSample(c.Request.Context(), sampleID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // Should not happen if update was successful
		return
	}

	var samplingDateNewptr *time.Time
	if !updatedSample.SamplingDate.IsZero() {
		samplingDateNewptr = &updatedSample.SamplingDate
	}
	var testDateNewptr *time.Time
	if !updatedSample.TestDate.IsZero() {
		testDateNewptr = &updatedSample.TestDate
	}

	// Map updated model back to response DTO
	resp := &dto.LabSampleResponse{
		ID:           updatedSample.ID,
		StationID:    updatedSample.StationID,
		SampleCode:   updatedSample.SampleCode,
		DepthFrom:    utils.GormDecimalToString(&updatedSample.DepthFrom),
		DepthTo:      utils.GormDecimalToString(&updatedSample.DepthTo),
		SampleType:   updatedSample.SampleType,
		SamplingDate: samplingDateNewptr,
		TestedBy:     updatedSample.TestedBy,
		LabName:      updatedSample.LabName,
		TestDate:     testDateNewptr,
		CreatedAt:    updatedSample.CreatedAt,
		UpdatedAt:    updatedSample.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteLabSample handles deleting a laboratory sample by ID.
// @Router /api/lab-samples/{id} [delete]
func (h *LaboratoryHandler) DeleteLabSample(c *gin.Context) {
	sampleID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lab sample")
	if appErr != nil {
		return // Response already handled
	}

	appErr = h.labService.DeleteSample(c.Request.Context(), sampleID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusNoContent, nil) // 204 No Content for successful deletion
}

// ListLabSamplesByStation handles listing lab samples for a specific station.
// @Router /api/stations/{station_id}/lab-samples [get]
func (h *LaboratoryHandler) ListLabSamplesByStation(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "id", "station")
	if appErr != nil {
		return // Response already handled
	}

	samples, appErr := h.labService.ListSamplesByStation(c.Request.Context(), stationID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map models to response DTOs
	sampleResponses := make([]dto.LabSampleResponse, len(samples))
	for i, s := range samples {

		var samplingDatePtr *time.Time
		if !s.SamplingDate.IsZero() { // valid (non-zero) date
			samplingDatePtr = &s.SamplingDate
		}

		// Convert TestDate from value to pointer:
		var testDatePtr *time.Time
		if !s.TestDate.IsZero() {
			testDatePtr = &s.TestDate
		}

		sampleResponses[i] = dto.LabSampleResponse{
			ID:           s.ID,
			StationID:    s.StationID,
			SampleCode:   s.SampleCode,
			DepthFrom:    utils.GormDecimalToString(&s.DepthFrom),
			DepthTo:      utils.GormDecimalToString(&s.DepthTo),
			SampleType:   s.SampleType,
			SamplingDate: samplingDatePtr,
			TestedBy:     s.TestedBy,
			LabName:      s.LabName,
			TestDate:     testDatePtr,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
		}
	}

	http_response.RespondWithSuccess(c, http.StatusOK, dto.ListLabSamplesResponse{
		Samples: sampleResponses,
		Total:   len(sampleResponses),
	})
}

// CreateUCSResult handles the creation of a new UCS result.
// @Router /api/ucs-results [post]
func (h *LaboratoryHandler) CreateUCSResult(c *gin.Context) {
	var req dto.CreateUCSResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Convert string decimal values from DTO to utils.GormDecimal
	ucsValueGd, appErr := utils.StringToGormDecimal(req.UCSValue)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	specimenDiameterGd, appErr := utils.StringToGormDecimal(req.SpecimenDiameter)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	specimenHeightGd, appErr := utils.StringToGormDecimal(req.SpecimenHeight)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map DTO to model
	ucs := &models.UCSResult{
		SampleID:         req.SampleID,
		UCSValue:         *ucsValueGd,
		Unit:             req.Unit,
		TestMethod:       req.TestMethod,
		SpecimenDiameter: *specimenDiameterGd,
		SpecimenHeight:   *specimenHeightGd,
		FailureMode:      req.FailureMode,
		Notes:            req.Notes,
	}

	createdUCS, appErr := h.labService.CreateUCSResult(c.Request.Context(), ucs)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map created model back to response DTO
	resp := &dto.UCSResultResponse{
		ID:               createdUCS.ID,
		SampleID:         createdUCS.SampleID,
		UCSValue:         utils.GormDecimalToString(&createdUCS.UCSValue),
		Unit:             createdUCS.Unit,
		TestMethod:       createdUCS.TestMethod,
		SpecimenDiameter: utils.GormDecimalToString(&createdUCS.SpecimenDiameter),
		SpecimenHeight:   utils.GormDecimalToString(&createdUCS.SpecimenHeight),
		FailureMode:      createdUCS.FailureMode,
		Notes:            createdUCS.Notes,
		CreatedAt:        createdUCS.CreatedAt,
		UpdatedAt:        createdUCS.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetUCSResultByID handles retrieving a single UCS result by ID.
// @Router /api/ucs-results/{id} [get]
func (h *LaboratoryHandler) GetUCSResultByID(c *gin.Context) {
	ucsID, appErr := gin_helpers.ParseIDFromContext(c, "id", "UCS result")
	if appErr != nil {
		return // Response already handled
	}

	ucs, appErr := h.labService.GetUCSResult(c.Request.Context(), ucsID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map model to response DTO
	resp := &dto.UCSResultResponse{
		ID:               ucs.ID,
		SampleID:         ucs.SampleID,
		UCSValue:         utils.GormDecimalToString(&ucs.UCSValue),
		Unit:             ucs.Unit,
		TestMethod:       ucs.TestMethod,
		SpecimenDiameter: utils.GormDecimalToString(&ucs.SpecimenDiameter),
		SpecimenHeight:   utils.GormDecimalToString(&ucs.SpecimenHeight),
		FailureMode:      ucs.FailureMode,
		Notes:            ucs.Notes,
		CreatedAt:        ucs.CreatedAt,
		UpdatedAt:        ucs.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateUCSResult handles updating an existing UCS result.
// @Router /api/ucs-results/{id} [put]
func (h *LaboratoryHandler) UpdateUCSResult(c *gin.Context) {
	ucsID, appErr := gin_helpers.ParseIDFromContext(c, "id", "UCS result")
	if appErr != nil {
		return // Response already handled
	}

	var req dto.UpdateUCSResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Create a models.UCSResult to pass to the service, applying only provided fields
	ucsToUpdate := &models.UCSResult{
		ID: ucsID,
	}

	if req.SampleID != nil {
		ucsToUpdate.SampleID = *req.SampleID
	}
	if req.UCSValue != nil {
		ucsValueGd, err := utils.StringToGormDecimal(*req.UCSValue)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		ucsToUpdate.UCSValue = *ucsValueGd
	}
	if req.Unit != nil {
		ucsToUpdate.Unit = *req.Unit
	}
	if req.TestMethod != nil {
		ucsToUpdate.TestMethod = *req.TestMethod
	}
	if req.SpecimenDiameter != nil {
		sdGd, err := utils.StringToGormDecimal(*req.SpecimenDiameter)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		ucsToUpdate.SpecimenDiameter = *sdGd
	}
	if req.SpecimenHeight != nil {
		shGd, err := utils.StringToGormDecimal(*req.SpecimenHeight)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		ucsToUpdate.SpecimenHeight = *shGd
	}
	if req.FailureMode != nil {
		ucsToUpdate.FailureMode = *req.FailureMode
	}
	if req.Notes != nil {
		ucsToUpdate.Notes = *req.Notes
	}

	appErr = h.labService.UpdateUCSResult(c.Request.Context(), ucsToUpdate)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Retrieve the updated UCS result to return the full, current state
	updatedUCS, appErr := h.labService.GetUCSResult(c.Request.Context(), ucsID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // Should not happen if update was successful
		return
	}

	// Map updated model back to response DTO
	resp := &dto.UCSResultResponse{
		ID:               updatedUCS.ID,
		SampleID:         updatedUCS.SampleID,
		UCSValue:         utils.GormDecimalToString(&updatedUCS.UCSValue),
		Unit:             updatedUCS.Unit,
		TestMethod:       updatedUCS.TestMethod,
		SpecimenDiameter: utils.GormDecimalToString(&updatedUCS.SpecimenDiameter),
		SpecimenHeight:   utils.GormDecimalToString(&updatedUCS.SpecimenHeight),
		FailureMode:      updatedUCS.FailureMode,
		Notes:            updatedUCS.Notes,
		CreatedAt:        updatedUCS.CreatedAt,
		UpdatedAt:        updatedUCS.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteUCSResult handles deleting a UCS result by ID.
// @Router /api/ucs-results/{id} [delete]
func (h *LaboratoryHandler) DeleteUCSResult(c *gin.Context) {
	ucsID, appErr := gin_helpers.ParseIDFromContext(c, "id", "UCS result")
	if appErr != nil {
		return // Response already handled
	}

	appErr = h.labService.DeleteUCSResult(c.Request.Context(), ucsID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusNoContent, nil) // 204 No Content for successful deletion
}

// ListUCSResultsBySample handles listing UCS results for a specific lab sample.
// @Router /api/lab-samples/{sample_id}/ucs-results [get]
func (h *LaboratoryHandler) ListUCSResultsBySample(c *gin.Context) {
	sampleID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lab sample")
	if appErr != nil {
		return // Response already handled
	}

	ucsResults, appErr := h.labService.ListUCSResultsBySample(c.Request.Context(), sampleID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map models to response DTOs
	ucsResponses := make([]dto.UCSResultResponse, len(ucsResults))
	for i, ucs := range ucsResults {
		ucsResponses[i] = dto.UCSResultResponse{
			ID:               ucs.ID,
			SampleID:         ucs.SampleID,
			UCSValue:         utils.GormDecimalToString(&ucs.UCSValue),
			Unit:             ucs.Unit,
			TestMethod:       ucs.TestMethod,
			SpecimenDiameter: utils.GormDecimalToString(&ucs.SpecimenDiameter),
			SpecimenHeight:   utils.GormDecimalToString(&ucs.SpecimenHeight),
			FailureMode:      ucs.FailureMode,
			Notes:            ucs.Notes,
			CreatedAt:        ucs.CreatedAt,
			UpdatedAt:        ucs.UpdatedAt,
		}
	}

	http_response.RespondWithSuccess(c, http.StatusOK, dto.ListUCSResultsResponse{
		UCSResults: ucsResponses,
		Total:      len(ucsResponses),
	})
}
