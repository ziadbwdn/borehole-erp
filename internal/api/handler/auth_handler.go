// Replace your entire auth_handler.go file content with this version.

package handler

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response" // Using the finalized response package
	"net/http"
	"strings" // Import strings package

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

// Register handles user registration requests.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	profile, appErr := h.authService.Register(c.Request.Context(), req)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	http_response.RespondWithSuccess(c, http.StatusCreated, profile) // CORRECTED
}

// Login handles user login requests.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	tokenResponse, appErr := h.authService.Login(c.Request.Context(), req)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, tokenResponse) // CORRECTED
}

// Logout handles user logout requests.
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	// LOGIC FIX: Extract the raw token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader || tokenString == "" {
		appErr := exception.NewAuthError("Authorization token is missing or malformed")
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	// Pass the RAW token to the service.
	appErr := h.authService.Logout(c.Request.Context(), tokenString)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Logout successful"}) // CORRECTED
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
