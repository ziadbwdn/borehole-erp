package services

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/pkg/geo_converter"
	"context"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
)

type StationServiceImpl struct {
	stationRepo     contract.StationRepository
	projectRepo     contract.ProjectRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

func NewStationService(stationRepo contract.StationRepository, projectRepo contract.ProjectRepository, activityService contract.UserActivityService, logger logger.Logger) contract.StationService {
	if stationRepo == nil { panic("stationRepo must not be nil") }
	if projectRepo == nil { panic("projectRepo must not be nil") }
	if activityService == nil { panic("activityService must not be nil") }
	if logger == nil { panic("logger must not be nil") }
	return &StationServiceImpl{
		stationRepo:     stationRepo,
		projectRepo:     projectRepo,
		activityService: activityService,
		logger:          logger,
	}
}

func (s *StationServiceImpl) CreateStation(ctx context.Context, station *models.Station, createdBy utils.BinaryUUID, logCtx models.ActivityLogContext) (*models.Station, *exception.AppError) {
	_, appErr := s.projectRepo.GetByID(ctx, station.ProjectID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return nil, exception.NewValidationError(fmt.Sprintf("Project with ID '%s' not found or inaccessible", station.ProjectID.String()))
		}
		return nil, appErr
	}
	station.ID = utils.NewBinaryUUID()
	station.CreatedAt = time.Now()
	station.UpdatedAt = time.Now()
	
	appErr = s.stationRepo.Create(ctx, station)
	if appErr != nil {
		return nil, appErr
	}
	
	newValueJSON, _ := json.Marshal(station)
	newValueStr := string(newValueJSON)
	stationIDStr := station.ID.String()
	details := fmt.Sprintf("New Station '%s' created.", station.StationCode)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeCreateStation, models.ResourceTypeStation, &stationIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log CreateStation activity", logErr, logger.Field{Key: "stationID", Value: station.ID.String()})
	}

	return station, nil
}

func (s *StationServiceImpl) GetStation(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError) {
	return s.stationRepo.GetByID(ctx, id)
}

func (s *StationServiceImpl) UpdateStation(ctx context.Context, stationID utils.BinaryUUID, req *dto.UpdateStationRequest, logCtx models.ActivityLogContext) (*models.Station, *exception.AppError) {
	existingStation, appErr := s.stationRepo.GetByID(ctx, stationID)
	if appErr != nil {
		return nil, appErr
	}
	oldValueJSON, err := json.Marshal(existingStation)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal old station value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "stationID", Value: stationID.String()})
	}
	oldValueStr := string(oldValueJSON)

	if req.StationCode != nil { existingStation.StationCode = *req.StationCode }
	if req.StationName != nil { existingStation.StationName = *req.StationName }
	if req.StationType != nil { existingStation.StationType = *req.StationType }
	if req.DrillingStatus != nil { existingStation.DrillingStatus = models.DrillingStatus(*req.DrillingStatus) }
	if req.DrillingDate != nil { existingStation.DrillingDate = req.DrillingDate }
	if req.GeologistName != nil { existingStation.GeologistName = *req.GeologistName }
	if req.Notes != nil { existingStation.Notes = *req.Notes }
	if req.Latitude != nil { if gd, err := utils.StringToGormDecimal(*req.Latitude); err == nil { existingStation.Latitude = *gd } }
	if req.Longitude != nil { if gd, err := utils.StringToGormDecimal(*req.Longitude); err == nil { existingStation.Longitude = *gd } }
	if req.Elevation != nil { if gd, err := utils.StringToGormDecimal(*req.Elevation); err == nil { existingStation.Elevation = *gd } }
	if req.GWL != nil { if gd, err := utils.StringToGormDecimal(*req.GWL); err == nil { existingStation.GWL = *gd } }
	if req.TotalDepth != nil { if gd, err := utils.StringToGormDecimal(*req.TotalDepth); err == nil { existingStation.TotalDepth = *gd } }
	existingStation.UpdatedAt = time.Now()

	if appErr := s.stationRepo.Update(ctx, existingStation); appErr != nil {
		return nil, appErr
	}

	newValueJSON, err := json.Marshal(existingStation)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal new station value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "stationID", Value: stationID.String()})
	}
	newValueStr := string(newValueJSON)
	stationIDStr := existingStation.ID.String()
	details := fmt.Sprintf("Station '%s' updated.", existingStation.StationCode)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeUpdateStation, models.ResourceTypeStation, &stationIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log UpdateStation activity", logErr, logger.Field{Key: "stationID", Value: existingStation.ID.String()})
	}
	
	return existingStation, nil
}

func (s *StationServiceImpl) DeleteStation(ctx context.Context, id utils.BinaryUUID, logCtx models.ActivityLogContext) *exception.AppError {
	stationToDelete, appErr := s.stationRepo.GetByID(ctx, id)
	if appErr != nil {
		return appErr
	}
	oldValueJSON, _ := json.Marshal(stationToDelete)
	oldValueStr := string(oldValueJSON)

	appErr = s.stationRepo.Delete(ctx, id)
	if appErr != nil {
		return appErr
	}

	stationIDStr := stationToDelete.ID.String()
	details := fmt.Sprintf("Station '%s' deleted.", stationToDelete.StationCode)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeDeleteStation, models.ResourceTypeStation, &stationIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log DeleteStation activity", logErr, logger.Field{Key: "stationID", Value: stationToDelete.ID.String()})
	}
	
	return nil
}

// Export to common Geographic Long/Lat Service
func (s *StationServiceImpl) ExportStationPointsGeo(ctx context.Context, projectID utils.BinaryUUID, logCtx models.ActivityLogContext) ([]*dto.StationPointGeoResponse, *exception.AppError) {
	stations, appErr := s.stationRepo.ListByProject(ctx, projectID)
	if appErr != nil {
		s.logger.Error(ctx, "Failed to list stations for Geo export", appErr, logger.Field{Key: "projectID", Value: projectID.String()})
		return nil, appErr
	}

	var response []*dto.StationPointGeoResponse
	for _, station := range stations {
		response = append(response, &dto.StationPointGeoResponse{
			StationCode: station.StationCode,
			Latitude:    utils.GormDecimalToString(&station.Latitude),
			Longitude:   utils.GormDecimalToString(&station.Longitude),
			Elevation:   utils.GormDecimalToString(&station.Elevation),
		})
	}
	// s.logger.Info(ctx, "Successfully prepared Geo export for project", logger.Field{Key: "projectID", Value: projectID.String()}, logger.Field{Key: "station_count", Value: len(response)})

    projectIDStr := projectID.String()
    details := fmt.Sprintf("User exported %d station points (UTM) for project %s.", len(response), projectIDStr)
    ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(
        ctx,
        logCtx.UserID,
        logCtx.Username,
        models.ActionTypeExportGeoData, // Use the new constant
        models.ResourceTypeExport,     // The action is on the project resource
        &projectIDStr,
        &ipAddr,
        &details,
        nil,
        nil,
    )
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log ExportStationPointsGeo activity", logErr, logger.Field{Key: "projectID", Value: projectID.String()})
	}

	return response, nil
}

// Export to UTM-Supported Format Service
func (s *StationServiceImpl) ExportStationPointsUTM(ctx context.Context, projectID utils.BinaryUUID, logCtx models.ActivityLogContext) ([]*dto.StationPointUTMResponse, *exception.AppError) {
	stations, appErr := s.stationRepo.ListByProject(ctx, projectID)
	if appErr != nil {
		s.logger.Error(ctx, "Failed to list stations for UTM export", appErr, logger.Field{Key: "projectID", Value: projectID.String()})
		return nil, appErr
	}

	var response []*dto.StationPointUTMResponse
	for _, station := range stations {
		latStr := utils.GormDecimalToString(&station.Latitude)
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			s.logger.Warn(ctx, "Failed to parse latitude for UTM conversion, skipping station",
				logger.Field{Key: "stationCode", Value: station.StationCode},
				logger.Field{Key: "latitudeValue", Value: latStr},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}

		// Step 2: Convert the GormDecimal Longitude to a float64 using the correct helper function.
		lonStr := utils.GormDecimalToString(&station.Longitude)
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			s.logger.Warn(ctx, "Failed to parse longitude for UTM conversion, skipping station",
				logger.Field{Key: "stationCode", Value: station.StationCode},
				logger.Field{Key: "longitudeValue", Value: lonStr},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}

		utmCoords, utmErr := geo_converter.ToUTM(lat, lon)

		if utmErr != nil {
			s.logger.Warn(ctx, "Failed to convert coordinates to UTM for station, skipping in export",
				logger.Field{Key: "stationCode", Value: station.StationCode},
				logger.Field{Key: "error", Value: utmErr.Error()})
			continue
		}

		// The response mapping remains correct.
		response = append(response, &dto.StationPointUTMResponse{
			StationCode: station.StationCode,
			Latitude:    latStr, // Reuse the string we already have
			Longitude:   lonStr, // Reuse the string we already have
			Elevation:   utils.GormDecimalToString(&station.Elevation),
			UTMZone:     fmt.Sprintf("%d%c", utmCoords.Zone, utmCoords.Hemisphere),
			Easting:     strconv.FormatFloat(utmCoords.Easting, 'f', -1, 64),
			Northing:    strconv.FormatFloat(utmCoords.Northing, 'f', -1, 64),
		})
	}
	
    // --- CORRECTED: Added the missing activity logging block ---
    projectIDStr := projectID.String()
    details := fmt.Sprintf("User exported %d station points (UTM) for project %s.", len(response), projectIDStr)
    ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(
        ctx,
        logCtx.UserID,
        logCtx.Username,
        models.ActionTypeExportUTMData, // Use the new constant
        models.ResourceTypeExport,     // The action is on the project resource
        &projectIDStr,
        &ipAddr,
        &details,
        nil,
        nil,
    )
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log ExportStationPointsUTM activity", logErr, logger.Field{Key: "projectID", Value: projectID.String()})
	}
    
	return response, nil
}

// List of stations by project
func (s *StationServiceImpl) ListStationsByProject(ctx context.Context, projectID utils.BinaryUUID) ([]*models.Station, *exception.AppError) {
	return s.stationRepo.ListByProject(ctx, projectID)
}