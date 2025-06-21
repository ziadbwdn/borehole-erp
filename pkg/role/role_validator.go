package role

import (
	"boreholedata-ms/internal/models"
)

// IsValid checks if a role is one of the predefined valid roles.
func IsValid(role models.UserRole) bool {
	switch role {
	case models.RoleAdmin, models.RoleGeologist, models.RoleEngineer,
		models.RoleLabTechnician, models.RoleGuest:
		return true
	default:
		return false
	}
}

// CanUpdateDrillingStatus checks if a role has permission to update drilling status.
func CanUpdateDrillingStatus(role models.UserRole) bool {
	return role == models.RoleEngineer || role == models.RoleGeologist
}

// CanManageProjects checks if a role has permission to manage projects.
func CanManageProjects(role models.UserRole) bool {
	return role == models.RoleAdmin || role == models.RoleEngineer
}

// CanManageUsers checks if a role has permission to manage users.
func CanManageUsers(role models.UserRole) bool {
	return role == models.RoleAdmin
}

// CanAccessLabData checks if a role has permission to access lab data.
func CanAccessLabData(role models.UserRole) bool {
	return role == models.RoleAdmin || role == models.RoleLabTechnician
}
