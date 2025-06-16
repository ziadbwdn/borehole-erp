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

	"github.com/gin-gonic/gin"
)

// StationHandler handles HTTP requests related to station management.
type StationHandler struct {
	stationService contract.StationService
}

// NewStationHandler creates and returns a new instance of StationHandler.
func NewStationHandler(stationService contract.StationService) *StationHandler {
	return &StationHandler{
		stationService: stationService,
	}
}

// CreateStation handles the creation of a new station.
// @Router /api/stations [post]
func (h *StationHandler) CreateStation(c *gin.Context) {
	var req dto.CreateStationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	createdBy, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Convert string decimal values from DTO to utils.GormDecimal
	latitudeGd, appErr := utils.StringToGormDecimal(req.Latitude)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	longitudeGd, appErr := utils.StringToGormDecimal(req.Longitude)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	elevationGd, appErr := utils.StringToGormDecimal(req.Elevation)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	GWLGd, appErr := utils.StringToGormDecimal(req.GWL)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	totalDepthGd, appErr := utils.StringToGormDecimal(req.TotalDepth)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map DTO to model
	station := &models.Station{
		ProjectID:     req.ProjectID,
		StationCode:   req.StationCode,
		StationName:   req.StationName,
		StationType:   req.StationType,
		Latitude:      *latitudeGd,
		Longitude:     *longitudeGd,
		GWL:           *GWLGd,
		Elevation:     *elevationGd,
		TotalDepth:    *totalDepthGd,
		DrillingDate:  req.DrillingDate,
		GeologistName: req.GeologistName,
		Notes:         req.Notes,
	}

	// Map DrillingStatus from DTO to model
	if req.DrillingStatus != "" {
		station.DrillingStatus = models.DrillingStatus(req.DrillingStatus)
	}
	// If not provided in DTO, GORM's default value ('planned') will be used when inserting.

	createdStation, appErr := h.stationService.CreateStation(c.Request.Context(), station, createdBy)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map created model back to response DTO
	resp := &dto.StationResponse{
		ID:             createdStation.ID,
		ProjectID:      createdStation.ProjectID,
		StationCode:    createdStation.StationCode,
		StationName:    createdStation.StationName,
		StationType:    createdStation.StationType,
		Latitude:       utils.GormDecimalToString(&createdStation.Latitude),
		Longitude:      utils.GormDecimalToString(&createdStation.Longitude),
		Elevation:      utils.GormDecimalToString(&createdStation.Elevation),
		GWL:            utils.GormDecimalToString(&createdStation.GWL),
		TotalDepth:     utils.GormDecimalToString(&createdStation.TotalDepth),
		DrillingDate:   createdStation.DrillingDate,
		DrillingStatus: string(createdStation.DrillingStatus), // --- FIX: Direct conversion or use as string ---
		GeologistName:  createdStation.GeologistName,
		Notes:          createdStation.Notes,
		CreatedAt:      createdStation.CreatedAt,
		UpdatedAt:      createdStation.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetStation handles retrieving a single station by ID.
// @Router /api/stations/{id} [get]
func (h *StationHandler) GetStation(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "id", "station")
	if appErr != nil {
		return // Response already handled by ParseIDFromContext
	}

	station, appErr := h.stationService.GetStation(c.Request.Context(), stationID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	resp := &dto.StationResponse{
		ID:             station.ID,
		ProjectID:      station.ProjectID,
		StationCode:    station.StationCode,
		StationName:    station.StationName,
		StationType:    station.StationType,
		Latitude:       utils.GormDecimalToString(&station.Latitude),
		Longitude:      utils.GormDecimalToString(&station.Longitude),
		Elevation:      utils.GormDecimalToString(&station.Elevation),
		TotalDepth:     utils.GormDecimalToString(&station.TotalDepth),
		GWL:            utils.GormDecimalToString(&station.GWL),
		DrillingDate:   station.DrillingDate,
		DrillingStatus: string(station.DrillingStatus), // --- FIX: Direct conversion or use as string ---
		GeologistName:  station.GeologistName,
		Notes:          station.Notes,
		CreatedAt:      station.CreatedAt,
		UpdatedAt:      station.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateStation handles updating an existing station.
// @Router /api/stations/{id} [put]
func (h *StationHandler) UpdateStation(c *gin.Context) {
	// Step 1: Get Station ID and request body.
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "id", "station")
	if appErr != nil {
		return
	}

	var req dto.UpdateStationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Step 2: Call the service with the ID and the DTO. No mapping, no model creation.
	updatedStation, appErr := h.stationService.UpdateStation(c.Request.Context(), stationID, &req)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Step 3: Map the final model to a clean response DTO.
	// This also fixes the ugly JSON output you were seeing for decimal fields.
	resp := dto.MapStationToResponse(updatedStation)
	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteStation handles deleting a station by ID.
// @Router /api/stations/{id} [delete]
func (h *StationHandler) DeleteStation(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "id", "station")
	if appErr != nil {
		return // Response already handled
	}

	appErr = h.stationService.DeleteStation(c.Request.Context(), stationID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusNoContent, nil) // 204 No Content for successful deletion
}

// ListStationsByProject handles listing stations for a specific project.
// @Router /api/projects/{project_id}/stations [get]
func (h *StationHandler) ListStationsByProject(c *gin.Context) {
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		return // Response already handled
	}

	stations, appErr := h.stationService.ListStationsByProject(c.Request.Context(), projectID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map models to response DTOs
	stationResponses := make([]dto.StationResponse, len(stations))
	for i, s := range stations {
		resp := dto.StationResponse{
			ID:             s.ID,
			ProjectID:      s.ProjectID,
			StationCode:    s.StationCode,
			StationName:    s.StationName,
			StationType:    s.StationType,
			Latitude:       utils.GormDecimalToString(&s.Latitude),
			Longitude:      utils.GormDecimalToString(&s.Longitude),
			Elevation:      utils.GormDecimalToString(&s.Elevation),
			GWL:            utils.GormDecimalToString(&s.GWL),
			TotalDepth:     utils.GormDecimalToString(&s.TotalDepth),
			DrillingDate:   s.DrillingDate,
			DrillingStatus: string(s.DrillingStatus), // --- FIX: Direct conversion or use as string ---
			GeologistName:  s.GeologistName,
			Notes:          s.Notes,
			CreatedAt:      s.CreatedAt,
			UpdatedAt:      s.UpdatedAt,
		}
		// Assuming you might want to include ProjectName if the Project relationship is loaded
		// if s.Project != nil {
		//  resp.ProjectName = s.Project.Name
		// }
		stationResponses[i] = resp
	}

	http_response.RespondWithSuccess(c, http.StatusOK, dto.ListStationsResponse{
		Stations: stationResponses,
		Total:    len(stationResponses),
	})
}
