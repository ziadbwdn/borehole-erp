package handler

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService contract.AuthService
}

func NewAuthHandler(authService contract.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

    // Get client IP address
    ipAddress := c.ClientIP()

    // Pass ipAddress to the service
	profile, appErr := h.authService.Register(c.Request.Context(), req, ipAddress)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, profile)
}

// Login 
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

    // Get client IP address
    ipAddress := c.ClientIP()

    // Pass ipAddress to the service
	tokenResponse, appErr := h.authService.Login(c.Request.Context(), req, ipAddress)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, tokenResponse)
}

// Logout handles user logout requests.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request: missing refresh_token in body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	ipAddress := c.ClientIP()
	appErr := h.authService.Logout(c.Request.Context(), req.RefreshToken, ipAddress)
	
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Logout successful"})
}

// GetProfile retrieves the authenticated user's profile.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		// GetUserIDFromContext already calls Abort, but we use our handler for a consistent response body
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	profile, appErr := h.authService.GetUserProfile(c.Request.Context(), userID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, profile) // CORRECTED
}

// UpdateProfile handles updating the authenticated user's profile.
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	profile, appErr := h.authService.UpdateUserProfile(c.Request.Context(), userID, req)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, profile) // CORRECTED
}

// refresh handler
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	// 1. Bind the incoming JSON request to the RefreshRequest DTO.
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request: missing refresh_token in body", err.Error())
		http_response.HandleAppError(c, appErr)
		return
	}

	// 3. Call the auth service to perform the token refresh logic.
	tokenResponse, appErr := h.authService.RefreshToken(c.Request.Context(), req)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, tokenResponse)
}