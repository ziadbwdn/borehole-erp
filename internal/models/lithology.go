package models

import (
	"time"

	"boreholedata-ms/internal/utils" // For BinaryUUID
)

// LithologyLog represents a single entry in a borehole's lithology log.
type LithologyLog struct {
	ID                 utils.BinaryUUID  `gorm:"primaryKey;type:binary(16)" json:"id"`        // Changed to utils.BinaryUUID
	StationID          utils.BinaryUUID  `gorm:"not null" json:"stationId"`                   // Changed to utils.BinaryUUID
	Station            *Station          `gorm:"foreignKey:StationID" json:"-"`               // Added json:"-" to prevent recursion in JSON
	DepthFrom          utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"depthFrom"` // Changed to utils.GormDecimal
	DepthTo            utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"depthTo"`   // Changed to utils.GormDecimal
	LithologyType      string            `gorm:"size:100;not null" json:"lithologyType"`
	RockColor          string            `gorm:"size:50" json:"rockColor"`
	GrainSize          string            `gorm:"size:50" json:"grainSize"`
	Texture            string            `gorm:"size:100" json:"texture"`
	Structure          string            `gorm:"size:100" json:"structure"`
	Hardness           string            `gorm:"size:50" json:"hardness"`
	Weathering         string            `gorm:"size:50" json:"weathering"`
	Fracturing         string            `gorm:"size:100" json:"fracturing"`
	Description        string            `json:"description"`
	RQDPercentage      utils.GormDecimal `gorm:"type:decimal(5,2)" json:"rqdPercentage"`      // Changed to utils.GormDecimal
	RecoveryPercentage utils.GormDecimal `gorm:"type:decimal(5,2)" json:"recoveryPercentage"` // Changed to utils.GormDecimal
	LoggedBy           string            `gorm:"size:100" json:"loggedBy"`
	LoggedDate         *time.Time       `gorm:"type:datetime" json:"loggedDate"`
	CreatedAt          time.Time         `gorm:"autoCreateTime" json:"createdAt"` // Corrected: Use autoCreateTime
	UpdatedAt          time.Time         `gorm:"autoUpdateTime" json:"updatedAt"` // Corrected: Use autoUpdateTime
}
