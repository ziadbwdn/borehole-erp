package models

import (
	"time"

	"boreholedata-ms/internal/utils" // For BinaryUUID
)

// DrillingStatus defines the status of a drilling operation for a station.
type DrillingStatus string

const (
	DrillingStatusPlanned    DrillingStatus = "planned"
	DrillingStatusInProgress DrillingStatus = "in_progress"
	DrillingStatusCompleted  DrillingStatus = "completed"
	DrillingStatusAbandoned  DrillingStatus = "abandoned"
	DrillingStatusSuspended  DrillingStatus = "suspended"
	DrillingStatusOnHold     DrillingStatus = "on_hold" // Added a new common status
)

// Station represents a borehole station within a project.
type Station struct {
	ID             utils.BinaryUUID  `gorm:"primaryKey;type:binary(16)" json:"id"` // Changed to utils.BinaryUUID
	ProjectID      utils.BinaryUUID  `gorm:"not null" json:"projectId"`            // Changed to utils.BinaryUUID
	Project        *Project          `gorm:"foreignKey:ProjectID" json:"-"`        // Added json:"-" to prevent recursion in JSON
	StationCode    string            `gorm:"size:50;not null;uniqueIndex" json:"stationCode"`
	StationName    string            `gorm:"size:100" json:"stationName"`
	StationType    string            `gorm:"size:50;not null" json:"stationType"`
	Latitude       utils.GormDecimal `gorm:"type:decimal(10,6);not null" json:"latitude"`        // Changed to utils.GormDecimal
	Longitude      utils.GormDecimal `gorm:"type:decimal(10,6);not null" json:"longitude"`       // Changed to utils.GormDecimal
	Elevation      utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"elevation"`        // Changed to utils.GormDecimal
	GWL            utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"groundwaterLevel"` // Changed to utils.GormDecimal
	TotalDepth     utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"totalDepth"`       // Changed to utils.GormDecimal
	DrillingDate   *time.Time        `gorm:"type:datetime" json:"drillingDate"`
	DrillingStatus DrillingStatus    `gorm:"type:varchar(50);default:'planned'" json:"drillingStatus"` // New field for drilling status
	GeologistName  string            `gorm:"size:100" json:"geologistName"`
	Notes          string            `json:"notes"`
	CreatedAt      time.Time         `gorm:"autoCreateTime" json:"createdAt"` // Corrected: Use autoCreateTime
	UpdatedAt      time.Time         `gorm:"autoUpdateTime" json:"updatedAt"` // Corrected: Use autoUpdateTime
}
