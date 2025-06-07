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
	"strconv" // For parsing depth range floats

	"github.com/gin-gonic/gin"
)

// LithologyHandler handles HTTP requests related to lithology log management.
type LithologyHandler struct {
	lithologyService contract.LithologyService
}

// NewLithologyHandler creates and returns a new instance of LithologyHandler.
func NewLithologyHandler(lithologyService contract.LithologyService) *LithologyHandler {
	return &LithologyHandler{
		lithologyService: lithologyService,
	}
}

// CreateLithologyLog handles the creation of a new lithology log entry.
// @Router /api/lithology-logs [post]
func (h *LithologyHandler) CreateLithologyLog(c *gin.Context) {
	var req dto.CreateLithologyLogRequest
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

	rqdPercentageGd, appErr := utils.StringToGormDecimal(req.RQDPercentage)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	recoveryPercentageGd, appErr := utils.StringToGormDecimal(req.RecoveryPercentage)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map DTO to model
	log := &models.LithologyLog{
		StationID:          req.StationID,
		DepthFrom:          *depthFromGd,
		DepthTo:            *depthToGd,
		LithologyType:      req.LithologyType,
		RockColor:          req.RockColor,
		GrainSize:          req.GrainSize,
		Texture:            req.Texture,
		Structure:          req.Structure,
		Hardness:           req.Hardness,
		Weathering:         req.Weathering,
		Fracturing:         req.Fracturing,
		Description:        req.Description,
		RQDPercentage:      *rqdPercentageGd,
		RecoveryPercentage: *recoveryPercentageGd,
		LoggedBy:           req.LoggedBy,
		LoggedDate:         req.LoggedDate,
	}

	createdLog, appErr := h.lithologyService.CreateLog(c.Request.Context(), log, createdBy)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map created model back to response DTO
	resp := &dto.LithologyLogResponse{
		ID:                 createdLog.ID,
		StationID:          createdLog.StationID,
		DepthFrom:          utils.GormDecimalToString(&createdLog.DepthFrom),
		DepthTo:            utils.GormDecimalToString(&createdLog.DepthTo),
		LithologyType:      createdLog.LithologyType,
		RockColor:          createdLog.RockColor,
		GrainSize:          createdLog.GrainSize,
		Texture:            createdLog.Texture,
		Structure:          createdLog.Structure,
		Hardness:           createdLog.Hardness,
		Weathering:         createdLog.Weathering,
		Fracturing:         createdLog.Fracturing,
		Description:        createdLog.Description,
		RQDPercentage:      utils.GormDecimalToString(&createdLog.RQDPercentage),
		RecoveryPercentage: utils.GormDecimalToString(&createdLog.RecoveryPercentage),
		LoggedBy:           createdLog.LoggedBy,
		LoggedDate:         createdLog.LoggedDate,
		CreatedAt:          createdLog.CreatedAt,
		UpdatedAt:          createdLog.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetLithologyLogByID handles retrieving a single lithology log by ID.
// @Router /api/lithology-logs/{id} [get]
func (h *LithologyHandler) GetLithologyLogByID(c *gin.Context) {
	logID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lithology log")
	if appErr != nil {
		return // Response already handled by ParseIDFromContext
	}

	log, appErr := h.lithologyService.GetLogByID(c.Request.Context(), logID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map model to response DTO
	resp := &dto.LithologyLogResponse{
		ID:                 log.ID,
		StationID:          log.StationID,
		DepthFrom:          utils.GormDecimalToString(&log.DepthFrom),
		DepthTo:            utils.GormDecimalToString(&log.DepthTo),
		LithologyType:      log.LithologyType,
		RockColor:          log.RockColor,
		GrainSize:          log.GrainSize,
		Texture:            log.Texture,
		Structure:          log.Structure,
		Hardness:           log.Hardness,
		Weathering:         log.Weathering,
		Fracturing:         log.Fracturing,
		Description:        log.Description,
		RQDPercentage:      utils.GormDecimalToString(&log.RQDPercentage),
		RecoveryPercentage: utils.GormDecimalToString(&log.RecoveryPercentage),
		LoggedBy:           log.LoggedBy,
		LoggedDate:         log.LoggedDate,
		CreatedAt:          log.CreatedAt,
		UpdatedAt:          log.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateLithologyLog handles updating an existing lithology log entry.
// @Router /api/lithology-logs/{id} [put]
func (h *LithologyHandler) UpdateLithologyLog(c *gin.Context) {
	logID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lithology log")
	if appErr != nil {
		return // Response already handled
	}

	var req dto.UpdateLithologyLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Fetch the existing log to apply updates
	existingLog, appErr := h.lithologyService.GetLogByID(c.Request.Context(), logID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Apply updates from DTO to model, converting string pointers to GormDecimal
	if req.StationID != nil {
		existingLog.StationID = *req.StationID
	}
	if req.DepthFrom != nil {
		dfGd, err := utils.StringToGormDecimal(*req.DepthFrom)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingLog.DepthFrom = *dfGd
	}
	if req.DepthTo != nil {
		dtGd, err := utils.StringToGormDecimal(*req.DepthTo)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingLog.DepthTo = *dtGd
	}
	if req.LithologyType != nil {
		existingLog.LithologyType = *req.LithologyType
	}
	if req.RockColor != nil {
		existingLog.RockColor = *req.RockColor
	}
	if req.GrainSize != nil {
		existingLog.GrainSize = *req.GrainSize
	}
	if req.Texture != nil {
		existingLog.Texture = *req.Texture
	}
	if req.Structure != nil {
		existingLog.Structure = *req.Structure
	}
	if req.Hardness != nil {
		existingLog.Hardness = *req.Hardness
	}
	if req.Weathering != nil {
		existingLog.Weathering = *req.Weathering
	}
	if req.Fracturing != nil {
		existingLog.Fracturing = *req.Fracturing
	}
	if req.Description != nil {
		existingLog.Description = *req.Description
	}
	if req.RQDPercentage != nil {
		rqdGd, err := utils.StringToGormDecimal(*req.RQDPercentage)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingLog.RQDPercentage = *rqdGd
	}
	if req.RecoveryPercentage != nil {
		recGd, err := utils.StringToGormDecimal(*req.RecoveryPercentage)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingLog.RecoveryPercentage = *recGd
	}
	if req.LoggedBy != nil {
		existingLog.LoggedBy = *req.LoggedBy
	}
	if req.LoggedDate != nil {
		existingLog.LoggedDate = req.LoggedDate
	}

	appErr = h.lithologyService.UpdateLog(c.Request.Context(), existingLog)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map updated model back to response DTO
	resp := &dto.LithologyLogResponse{
		ID:                 existingLog.ID,
		StationID:          existingLog.StationID,
		DepthFrom:          utils.GormDecimalToString(&existingLog.DepthFrom),
		DepthTo:            utils.GormDecimalToString(&existingLog.DepthTo),
		LithologyType:      existingLog.LithologyType,
		RockColor:          existingLog.RockColor,
		GrainSize:          existingLog.GrainSize,
		Texture:            existingLog.Texture,
		Structure:          existingLog.Structure,
		Hardness:           existingLog.Hardness,
		Weathering:         existingLog.Weathering,
		Fracturing:         existingLog.Fracturing,
		Description:        existingLog.Description,
		RQDPercentage:      utils.GormDecimalToString(&existingLog.RQDPercentage),
		RecoveryPercentage: utils.GormDecimalToString(&existingLog.RecoveryPercentage),
		LoggedBy:           existingLog.LoggedBy,
		LoggedDate:         existingLog.LoggedDate,
		CreatedAt:          existingLog.CreatedAt,
		UpdatedAt:          existingLog.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteLithologyLog handles deleting a lithology log by ID.
// @Router /api/lithology-logs/{id} [delete]
func (h *LithologyHandler) DeleteLithologyLog(c *gin.Context) {
	logID, appErr := gin_helpers.ParseIDFromContext(c, "id", "lithology log")
	if appErr != nil {
		return // Response already handled
	}

	appErr = h.lithologyService.DeleteLog(c.Request.Context(), logID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusNoContent, nil) // 204 No Content for successful deletion
}

// ListLithologyLogsByStation handles listing lithology logs for a specific station.
// @Router /api/stations/{station_id}/lithology-logs [get]
func (h *LithologyHandler) ListLithologyLogsByStation(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "id", "station")
	if appErr != nil {
		return // Response already handled
	}

	logs, appErr := h.lithologyService.ListLogsByStation(c.Request.Context(), stationID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map models to response DTOs
	logResponses := make([]dto.LithologyLogResponse, len(logs))
	for i, l := range logs {
		logResponses[i] = dto.LithologyLogResponse{
			ID:                 l.ID,
			StationID:          l.StationID,
			DepthFrom:          utils.GormDecimalToString(&l.DepthFrom),
			DepthTo:            utils.GormDecimalToString(&l.DepthTo),
			LithologyType:      l.LithologyType,
			RockColor:          l.RockColor,
			GrainSize:          l.GrainSize,
			Texture:            l.Texture,
			Structure:          l.Structure,
			Hardness:           l.Hardness,
			Weathering:         l.Weathering,
			Fracturing:         l.Fracturing,
			Description:        l.Description,
			RQDPercentage:      utils.GormDecimalToString(&l.RQDPercentage),
			RecoveryPercentage: utils.GormDecimalToString(&l.RecoveryPercentage),
			LoggedBy:           l.LoggedBy,
			LoggedDate:         l.LoggedDate,
			CreatedAt:          l.CreatedAt,
			UpdatedAt:          l.UpdatedAt,
		}
	}

	http_response.RespondWithSuccess(c, http.StatusOK, dto.ListLithologyLogsResponse{
		Logs:  logResponses,
		Total: len(logResponses),
	})
}

// ListLithologyLogsByDepthRange handles listing lithology logs for a specific station within a depth range.
// @Router /api/stations/{station_id}/lithology-logs/by-depth [get]
func (h *LithologyHandler) ListLithologyLogsByDepthRange(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "id", "station")
	if appErr != nil {
		return // Response already handled
	}

	minDepthStr := c.Query("min_depth")
	maxDepthStr := c.Query("max_depth")

	minDepth, err := strconv.ParseFloat(minDepthStr, 64)
	if err != nil {
		appErr := exception.NewValidationError("Invalid min_depth format", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	maxDepth, err := strconv.ParseFloat(maxDepthStr, 64)
	if err != nil {
		appErr := exception.NewValidationError("Invalid max_depth format", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	logs, appErr := h.lithologyService.ListLogsByDepthRange(c.Request.Context(), stationID, minDepth, maxDepth)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map models to response DTOs
	logResponses := make([]dto.LithologyLogResponse, len(logs))
	for i, l := range logs {
		logResponses[i] = dto.LithologyLogResponse{
			ID:                 l.ID,
			StationID:          l.StationID,
			DepthFrom:          utils.GormDecimalToString(&l.DepthFrom),
			DepthTo:            utils.GormDecimalToString(&l.DepthTo),
			LithologyType:      l.LithologyType,
			RockColor:          l.RockColor,
			GrainSize:          l.GrainSize,
			Texture:            l.Texture,
			Structure:          l.Structure,
			Hardness:           l.Hardness,
			Weathering:         l.Weathering,
			Fracturing:         l.Fracturing,
			Description:        l.Description,
			RQDPercentage:      utils.GormDecimalToString(&l.RQDPercentage),
			RecoveryPercentage: utils.GormDecimalToString(&l.RecoveryPercentage),
			LoggedBy:           l.LoggedBy,
			LoggedDate:         l.LoggedDate,
			CreatedAt:          l.CreatedAt,
			UpdatedAt:          l.UpdatedAt,
		}
	}

	http_response.RespondWithSuccess(c, http.StatusOK, dto.ListLithologyLogsResponse{
		Logs:  logResponses,
		Total: len(logResponses),
	})
}
