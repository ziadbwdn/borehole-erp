package models

import (
	"time"

	"boreholedata-ms/internal/utils" // For BinaryUUID
)

// Project represents a project in the borehole data management system.
type Project struct {
	ID          utils.BinaryUUID  `gorm:"primaryKey;type:binary(16)" json:"id"` // Changed to utils.BinaryUUID
	Name        string            `gorm:"size:100;not null" json:"name"`
	Description string            `json:"description"`
	Location    string            `gorm:"size:200" json:"location"`
	Latitude      utils.GormDecimal `gorm:"type:decimal(10,6);not null" json:"latitude"`  // Changed to utils.GormDecimal
	Longitude     utils.GormDecimal `gorm:"type:decimal(10,6);not null" json:"longitude"` // Changed to utils.GormDecimal
	Elevation     utils.GormDecimal `gorm:"type:decimal(8,2);not null" json:"elevation"`  // Changed to utils.GormDecimal
	StartDate   *time.Time       `gorm:"type:datetime" json:"startDate"` // Changed to pointer and removed not null
	EndDate     *time.Time       `gorm:"type:datetime" json:"endDate"`   // Changed to pointer and removed not null
	Status      string            `gorm:"size:20;default:'active'" json:"status"`
	CreatedBy   utils.BinaryUUID  `json:"createdBy"`                       // Changed to utils.BinaryUUID
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"createdAt"` // Corrected: Use autoCreateTime
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updatedAt"` // Corrected: Use autoUpdateTime
}
