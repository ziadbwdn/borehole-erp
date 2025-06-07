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

// ProjectHandler handles HTTP requests related to project management.
type ProjectHandler struct {
	projectService contract.ProjectService
}

// NewProjectHandler creates and returns a new instance of ProjectHandler.
func NewProjectHandler(projectService contract.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject handles the creation of a new project.
// @Router /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
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

	// Map DTO to model
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		Latitude:    *latitudeGd,  // Assign the GormDecimal struct
		Longitude:   *longitudeGd, // Assign the GormDecimal struct
		Elevation:   *elevationGd, // Assign the GormDecimal struct
	}

	createdProject, appErr := h.projectService.CreateProject(c.Request.Context(), project, createdBy)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map created model back to response DTO
	resp := &dto.ProjectResponse{
		ID:          createdProject.ID,
		Name:        createdProject.Name,
		Description: createdProject.Description,
		Location:    createdProject.Location,
		Latitude:    utils.GormDecimalToString(&createdProject.Latitude),  // Use GormDecimalToString
		Longitude:   utils.GormDecimalToString(&createdProject.Longitude), // Use GormDecimalToString
		Elevation:   utils.GormDecimalToString(&createdProject.Elevation), // Use GormDecimalToString
		StartDate:   createdProject.StartDate,
		EndDate:     createdProject.EndDate,
		Status:      createdProject.Status,
		CreatedBy:   createdProject.CreatedBy,
		CreatedAt:   createdProject.CreatedAt,
		UpdatedAt:   createdProject.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetProject handles retrieving a single project by ID.
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		return // Response already handled by ParseIDFromContext
	}

	project, appErr := h.projectService.GetProject(c.Request.Context(), projectID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	resp := &dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Location:    project.Location,
		Latitude:    utils.GormDecimalToString(&project.Latitude),
		Longitude:   utils.GormDecimalToString(&project.Longitude),
		Elevation:   utils.GormDecimalToString(&project.Elevation),
		StartDate:   project.StartDate,
		EndDate:     project.EndDate,
		Status:      project.Status,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateProject handles updating an existing project.
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		return // Response already handled
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// Fetch the existing project to apply updates
	existingProject, appErr := h.projectService.GetProject(c.Request.Context(), projectID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Apply updates from DTO to model, converting string pointers to GormDecimal
	if req.Name != nil {
		existingProject.Name = *req.Name
	}
	if req.Description != nil {
		existingProject.Description = *req.Description
	}
	if req.Location != nil {
		existingProject.Location = *req.Location
	}
	if req.Latitude != nil {
		latGd, err := utils.StringToGormDecimal(*req.Latitude)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingProject.Latitude = *latGd // Assign the GormDecimal struct
	}
	if req.Longitude != nil {
		lonGd, err := utils.StringToGormDecimal(*req.Longitude)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingProject.Longitude = *lonGd // Assign the GormDecimal struct
	}
	if req.Elevation != nil {
		elevGd, err := utils.StringToGormDecimal(*req.Elevation)
		if err != nil {
			http_response.HandleAppError(c, err)
			return
		}
		existingProject.Elevation = *elevGd // Assign the GormDecimal struct
	}
	if req.Status != nil {
		existingProject.Status = *req.Status
	}

	// Handle nullable StartDate and EndDate
	if req.StartDate != nil {
		existingProject.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		existingProject.EndDate = req.EndDate
	}
	// Note: StartDate, EndDate are not in UpdateProjectRequest DTO provided,
	// so they are not updated here. Add them to DTO if needed.

	appErr = h.projectService.UpdateProject(c.Request.Context(), existingProject)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map updated model back to response DTO
	resp := &dto.ProjectResponse{
		ID:          existingProject.ID,
		Name:        existingProject.Name,
		Description: existingProject.Description,
		Location:    existingProject.Location,
		Latitude:    utils.GormDecimalToString(&existingProject.Latitude),
		Longitude:   utils.GormDecimalToString(&existingProject.Longitude),
		Elevation:   utils.GormDecimalToString(&existingProject.Elevation),
		StartDate:   existingProject.StartDate,
		EndDate:     existingProject.EndDate,
		Status:      existingProject.Status,
		CreatedBy:   existingProject.CreatedBy,
		CreatedAt:   existingProject.CreatedAt,
		UpdatedAt:   existingProject.UpdatedAt,
	}

	http_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteProject handles deleting a project by ID.
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		return // Response already handled
	}

	appErr = h.projectService.DeleteProject(c.Request.Context(), projectID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusNoContent, nil) // 204 No Content for successful deletion
}

// ListProjects handles listing all projects for the authenticated user.
// @Router /api/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// Extract userID from context, set by authentication middleware
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	projects, appErr := h.projectService.ListProjectsByUser(c.Request.Context(), userID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map models to response DTOs
	projectResponses := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		projectResponses[i] = dto.ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Location:    p.Location,
			Latitude:    utils.GormDecimalToString(&p.Latitude),
			Longitude:   utils.GormDecimalToString(&p.Longitude),
			Elevation:   utils.GormDecimalToString(&p.Elevation),
			StartDate:   p.StartDate,
			EndDate:     p.EndDate,
			Status:      p.Status,
			CreatedBy:   p.CreatedBy,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	http_response.RespondWithSuccess(c, http.StatusOK, dto.ListProjectsResponse{
		Projects: projectResponses,
		Total:    len(projectResponses),
	})
}

// GetProjectStations handles retrieving stations for a specific project.
// @Router /api/projects/{id}/stations [get]
func (h *ProjectHandler) GetProjectStations(c *gin.Context) {
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		return // Response already handled
	}

	// This part requires a StationService or direct repository access.
	// For now, it's a placeholder. You would typically call:
	// stations, appErr := h.stationService.ListStationsByProject(c.Request.Context(), projectID)
	// if appErr != nil {
	//     http_response.HandleAppError(c, appErr)
	//     return
	// }
	// http_response.RespondWithSuccess(c, http.StatusOK, stations)

	// Mock response for now
	_ = projectID                                                          // Use blank identifier to suppress "declared and not used" warning
	http_response.RespondWithSuccess(c, http.StatusOK, []models.Station{}) // Return empty slice
}
