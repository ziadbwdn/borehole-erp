package models

import (
	"time"

	"boreholedata-ms/internal/utils" // For BinaryUUID

)

// UCSResult represents the Unconfined Compressive Strength test results for a lab sample.
type UCSResult struct {
	ID               utils.BinaryUUID  `gorm:"primaryKey;type:binary(16)" json:"id"`       // Changed to utils.BinaryUUID
	SampleID         utils.BinaryUUID  `gorm:"not null" json:"sampleId"`                   // Changed to utils.BinaryUUID
	Sample           *LabSample        `gorm:"foreignKey:SampleID" json:"-"`               // Added json:"-" to prevent recursion in JSON
	UCSValue         utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"ucsValue"` // Changed to utils.GormDecimal
	Unit             string            `gorm:"size:10;default:'MPa'" json:"unit"`
	TestMethod       string            `gorm:"size:50" json:"testMethod"`
	SpecimenDiameter utils.GormDecimal `gorm:"type:decimal(6,3)" json:"specimenDiameter"` // Changed to utils.GormDecimal
	SpecimenHeight   utils.GormDecimal `gorm:"type:decimal(6,3)" json:"specimenHeight"`   // Changed to utils.GormDecimal
	FailureMode      string            `gorm:"size:100" json:"failureMode"`
	Notes            string            `json:"notes"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"createdAt"` // Corrected: Use autoCreateTime
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updatedAt"` // Corrected: Use autoUpdateTime
}
