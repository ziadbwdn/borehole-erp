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
	authService    contract.AuthService
}

// NewProjectHandler creates and returns a new instance of ProjectHandler.
func NewProjectHandler(projectService contract.ProjectService, authService contract.AuthService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		authService:    authService,
	}
}

// CreateProject handles the creation of a new project.
// @Router /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_response.HandleAppError(c, exception.NewValidationError("Invalid request body", err.Error()))
		return
	}

	// 1. Get UserID from context
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// 2. Get Username by calling AuthService
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// 3. Prepare the ActivityLogContext
	logCtx := models.ActivityLogContext{
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}
	
	// Map DTO to model (existing logic)
	var projectStartDate, projectEndDate *time.Time
	if !req.StartDate.IsZero() { projectStartDate = &req.StartDate }
	if !req.EndDate.IsZero() { projectEndDate = &req.EndDate }
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   projectStartDate,
		EndDate:     projectEndDate,
	}

	// 4. Call the service with the new logCtx parameter
	createdProject, appErr := h.projectService.CreateProject(c.Request.Context(), project, userID, logCtx)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Map created model back to response DTO (existing logic)
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
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project"); if appErr != nil { return }
	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_response.HandleAppError(c, exception.NewValidationError("Invalid request body", err.Error())); return
	}

	userID, appErr := gin_helpers.GetUserIDFromContext(c); 
	if appErr != nil {
		http_response.HandleAppError(c, appErr); 
		return 
	}

	userRole, appErr := gin_helpers.GetUserRoleFromContext(c); 
	if appErr != nil { 
		http_response.HandleAppError(c, appErr); 
		return 
	}

	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID); 
	if appErr != nil { 
		http_response.HandleAppError(c, appErr); 
		return 
	}
	
	logCtx := models.ActivityLogContext{ 
		UserID: userID.String(), 
		Username: username, 
		IPAddress: c.ClientIP()}

	projectToUpdate := &models.Project{ ID: projectID }
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
		projectToUpdate.Status = *req.Status 
	}

	if req.StartDate != nil && !req.StartDate.IsZero() { 
		projectToUpdate.StartDate = req.StartDate 
	}

	if req.EndDate != nil && !req.EndDate.IsZero() { 
		projectToUpdate.EndDate = req.EndDate 
	}

	appErr = h.projectService.UpdateProject(c.Request.Context(), projectToUpdate, userRole, logCtx)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	updatedProject, appErr := h.projectService.GetProject(c.Request.Context(), projectID)
	if appErr != nil {
		http_response.HandleAppError(c, exception.NewInternalError("Failed to retrieve updated project", appErr)); return
	}

	// map DTO Response
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
		return
	}

	// 1. Get UserID and Username for logging
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// 2. Prepare the ActivityLogContext
	logCtx := models.ActivityLogContext{
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	// 3. Call the service with logCtx
	appErr = h.projectService.DeleteProject(c.Request.Context(), projectID, logCtx)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusNoContent, nil)
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
