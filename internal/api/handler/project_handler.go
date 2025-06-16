package handler

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"
	"net/http"
	"time"

	// Import time for handling date/time fields
	"github.com/gin-gonic/gin"
)

// ProjectHandler handles HTTP requests related to project management.
type ProjectHandler struct {
	projectService contract.ProjectService
	// If GetProjectStations needs StationService, it should be here too
	// stationService contract.StationService
}

// NewProjectHandler creates and returns a new instance of ProjectHandler.
func NewProjectHandler(projectService contract.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		// stationService: stationService, // Uncomment if needed
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

	// --- FIX: Correctly map time.Time from DTO to *time.Time in Model ---
	var projectStartDate *time.Time
	if !req.StartDate.IsZero() {
		projectStartDate = &req.StartDate
	}

	var projectEndDate *time.Time
	if !req.EndDate.IsZero() {
		projectEndDate = &req.EndDate
	}
	// --- END FIX ---

	// Map DTO to model
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   projectStartDate, // Assign the correctly mapped *time.Time
		EndDate:     projectEndDate,   // Assign the correctly mapped *time.Time
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

	// Get user role from context
	userRole, appErr := gin_helpers.GetUserRoleFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Create a new Project model to hold only the fields intended for update.
	// The service layer will fetch the full existing project and apply these changes.
	projectToUpdate := &models.Project{
		ID: projectID, // Crucial: Set the ID of the project to be updated
	}

	if req.Name != nil {
		projectToUpdate.Name = *req.Name
	}
	if req.Description != nil {
		projectToUpdate.Description = *req.Description
	}
	if req.Location != nil {
		projectToUpdate.Location = *req.Location
	}
	if req.Status != nil {
		projectToUpdate.Status = *req.Status // Pass status, service will check role
	}

	// --- CORRECTED Date Handling Logic ---
	// If req.StartDate is not nil AND not a zero value, then assign it.
	// Otherwise (req.StartDate is nil OR req.StartDate is a zero value),
	// we do not touch projectToUpdate.StartDate, meaning it retains its default zero value.
	// The service will then interpret this as "no change".
	if req.StartDate != nil {
		if !req.StartDate.IsZero() {
			projectToUpdate.StartDate = req.StartDate
		}
		// The `else` (req.StartDate is not nil, but is zero) is intentionally empty here.
		// This means projectToUpdate.StartDate remains its zero value (time.Time{}).
	}

	if req.EndDate != nil {
		if !req.EndDate.IsZero() {
			projectToUpdate.EndDate = req.EndDate
		}
		// The `else` (req.EndDate is not nil, but is zero) is intentionally empty here.
		// This means projectToUpdate.EndDate remains its zero value (time.Time{}).
	}
	// --- END CORRECTED Date Handling Logic ---

	// FIX: Pass userRole to the service call
	appErr = h.projectService.UpdateProject(c.Request.Context(), projectToUpdate, userRole)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// After a successful update, it's good practice to fetch the updated project
	// to return the most current state to the client.
	updatedProject, appErr := h.projectService.GetProject(c.Request.Context(), projectID)
	if appErr != nil {
		// This should ideally not happen after a successful update, but handle defensively.
		http_response.HandleAppError(c, exception.NewInternalError("Failed to retrieve updated project", appErr))
		return
	}

	// Map updated model back to response DTO
	resp := &dto.ProjectResponse{
		ID:          updatedProject.ID,
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
		Location:    updatedProject.Location,
		StartDate:   updatedProject.StartDate,
		EndDate:     updatedProject.EndDate,
		Status:      updatedProject.Status,
		CreatedBy:   updatedProject.CreatedBy,
		CreatedAt:   updatedProject.CreatedAt,
		UpdatedAt:   updatedProject.UpdatedAt,
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

	// This part requires a StationService instance within the ProjectHandler.
	// You would typically have:
	// h.stationService contract.StationService in the handler struct,
	// and initialize it in NewProjectHandler.
	// Then call:
	// stations, appErr := h.stationService.ListStationsByProject(c.Request.Context(), projectID)
	// if appErr != nil {
	// 	http_response.HandleAppError(c, appErr)
	// 	return
	// }
	// http_response.RespondWithSuccess(c, http.StatusOK, stations)

	// Mock response for now to allow compilation if StationService isn't passed in
	_ = projectID
	http_response.RespondWithSuccess(c, http.StatusOK, []models.Station{}) // Return empty slice
}
