package repository

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"errors" // For errors.Is
	"fmt"
	"time" // For setting timestamps in Create/Update

	pbdecimal "google.golang.org/genproto/googleapis/type/decimal" // For LithologyLog decimal fields
	"gorm.io/gorm"
)

// LithologyRepositoryImpl is a concrete implementation of the contract.LithologyRepository interface.
type LithologyRepositoryImpl struct {
	db *gorm.DB
}

// NewLithologyRepository creates and returns a new instance of LithologyRepository.
func NewLithologyRepository(db *gorm.DB) contract.LithologyRepository {
	return &LithologyRepositoryImpl{db: db}
}

// CreateLog inserts a new lithology log record into the database.
func (r *LithologyRepositoryImpl) CreateLog(ctx context.Context, log *models.LithologyLog) *exception.AppError {
	// GORM will automatically handle CreatedAt/UpdatedAt if gorm.Model is embedded,
	// otherwise ensure they are set before creation.
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	log.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Create(log)
	if result.Error != nil {
		return exception.NewDatabaseError(fmt.Sprintf("create lithology log for station '%s'", log.StationID.String()), result.Error)
	}
	return nil
}

// GetLogByID retrieves a single lithology log record by its ID.
func (r *LithologyRepositoryImpl) GetLogByID(ctx context.Context, id utils.BinaryUUID) (*models.LithologyLog, *exception.AppError) {
	log := &models.LithologyLog{}
	result := r.db.WithContext(ctx).First(log, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundError("LithologyLog", id)
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("get lithology log by ID '%s'", id.String()), result.Error)
	}
	return log, nil
}

// UpdateLog updates an existing lithology log record in the database.
func (r *LithologyRepositoryImpl) UpdateLog(ctx context.Context, log *models.LithologyLog) *exception.AppError {
	log.UpdatedAt = time.Now() // Update timestamp

	result := r.db.WithContext(ctx).Save(log) // Save will update if primary key exists
	if result.Error != nil {
		return exception.NewDatabaseError(fmt.Sprintf("update lithology log '%s'", log.ID.String()), result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("LithologyLog", log.ID) // No rows affected usually means not found
	}
	return nil
}

// DeleteLog deletes a lithology log record by its ID.
func (r *LithologyRepositoryImpl) DeleteLog(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	result := r.db.WithContext(ctx).Delete(&models.LithologyLog{}, "id = ?", id)
	if result.Error != nil {
		return exception.NewDatabaseError(fmt.Sprintf("delete lithology log '%s'", id.String()), result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("LithologyLog", id) // No rows affected means not found
	}
	return nil
}

// ListLogsByStation retrieves all lithology logs for a given station ID.
func (r *LithologyRepositoryImpl) ListLogsByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LithologyLog, *exception.AppError) {
	var logs []*models.LithologyLog
	result := r.db.WithContext(ctx).Where("station_id = ?", stationID).Find(&logs)
	if result.Error != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("list lithology logs for station '%s'", stationID.String()), result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exception.NewNotFoundError("LithologyLogs for Station", stationID)
	}
	return logs, nil
}

// ListLogsByDepthRange retrieves lithology logs for a given station ID within a specified depth range.
// Note: The interface expects float64 for minDepth/maxDepth, but the model uses pbdecimal.Decimal.
// Conversion is performed here, which might introduce minor precision issues if not handled carefully.
func (r *LithologyRepositoryImpl) ListLogsByDepthRange(ctx context.Context, stationID utils.BinaryUUID, minDepth, maxDepth float64) ([]*models.LithologyLog, *exception.AppError) {
	var logs []*models.LithologyLog

	// Convert float64 to pbdecimal.Decimal for comparison in the database query.
	// This is a critical point for precision. Ensure your database column type (e.g., DECIMAL)
	// and GORM mapping are robust enough to handle the string representation from pbdecimal.Decimal.
	minDepthPb := &pbdecimal.Decimal{Value: fmt.Sprintf("%.2f", minDepth)} // Format to 2 decimal places as per model's (8,2)
	maxDepthPb := &pbdecimal.Decimal{Value: fmt.Sprintf("%.2f", maxDepth)} // Format to 2 decimal places

	result := r.db.WithContext(ctx).
		Where("station_id = ? AND depth_from >= ? AND depth_to <= ?", stationID, minDepthPb.Value, maxDepthPb.Value).
		Find(&logs)

	if result.Error != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("list lithology logs by depth range for station '%s'", stationID.String()), result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exception.NewNotFoundError(
			fmt.Sprintf("LithologyLogs for Station '%s' within depth range %.2f-%.2f", stationID.String(), minDepth, maxDepth),
			stationID,
		)
	}
	return logs, nil
}
