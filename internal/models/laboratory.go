package models

import (
	"time"

	"boreholedata-ms/internal/utils" // For BinaryUUID
)

// LabSample represents a laboratory sample taken from a borehole station.
type LabSample struct {
	ID           utils.BinaryUUID  `gorm:"primaryKey;type:binary(16)" json:"id"` // Changed to utils.BinaryUUID
	StationID    utils.BinaryUUID  `gorm:"not null" json:"stationId"`            // Changed to utils.BinaryUUID
	SampleCode   string            `gorm:"size:50;not null;uniqueIndex" json:"sampleCode"`
	DepthFrom    utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"depthFrom"` // Changed to utils.GormDecimal
	DepthTo      utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"depthTo"`   // Changed to utils.GormDecimal
	SampleType   string            `gorm:"size:50" json:"sampleType"`
	SamplingDate time.Time       `gorm:"type:datetime" json:"samplingDate"` // Changed to pointer and removed not null
	TestedBy     string            `gorm:"size:100" json:"testedBy"`
	LabName      string            `gorm:"size:100" json:"labName"`
	TestDate     time.Time       `gorm:"type:datetime" json:"testDate"`     // Changed to pointer and removed not null
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"createdAt"` // Corrected: Use autoCreateTime
	UpdatedAt    time.Time         `gorm:"autoUpdateTime" json:"updatedAt"` // Corrected: Use autoUpdateTime

	// Relations
	Station *Station `gorm:"foreignKey:StationID" json:"-"` // Added json:"-" to prevent recursion in JSON
}
