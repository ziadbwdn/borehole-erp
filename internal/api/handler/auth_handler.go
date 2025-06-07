package handler

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract" // Note the alias for the service contract
	"boreholedata-ms/pkg/gin_helpers"              // Import the new gin_helpers package
	"boreholedata-ms/pkg/http_response"            // Import the new http_response package
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests related to authentication.
// It depends on an implementation of the AuthService interface to perform business logic.
type AuthHandler struct {
	authService contract.AuthService // Dependency on the authentication service
}

// NewAuthHandler creates and returns a new instance of AuthHandler.
// It takes an AuthService implementation as a dependency.
func NewAuthHandler(authService contract.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration requests.
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	// Bind the JSON request body to the RegisterRequest DTO.
	if err := c.ShouldBindJSON(&req); err != nil {
		// If binding fails, it's a bad request.
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	// Call the Register method of the AuthService.
	profile, appErr := h.authService.Register(c.Request.Context(), req)
	if appErr != nil {
		// If the service returns an AppError, map it to an HTTP status and response.
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	// On successful registration, return a 201 Created status with the user profile.
	http_response.RespondWithSuccess(c, http.StatusCreated, profile) // Use the utility function from pkg/http_response
}

// Login handles user login requests.
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	// Bind the JSON request body to the LoginRequest DTO.
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	// Call the Login method of the AuthService.
	tokenResponse, appErr := h.authService.Login(c.Request.Context(), req)
	if appErr != nil {
		// If the service returns an AppError, map it to an HTTP status and response.
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	// On successful login, return a 200 OK status with the token.
	http_response.RespondWithSuccess(c, http.StatusOK, tokenResponse) // Use the utility function from pkg/http_response
}

// Logout handles user logout requests.
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		appErr := exception.NewAuthError("Authorization token missing")
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	// Pass the token to the service for invalidation (e.g., blacklisting)
	appErr := h.authService.Logout(c.Request.Context(), authHeader)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Logout successful"}) // Use the utility function from pkg/http_response
}

// GetProfile retrieves the authenticated user's profile.
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Use the new helper function from pkg/gin_helpers
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Call the AuthService method to get the user profile.
	profile, appErr := h.authService.GetUserProfile(c.Request.Context(), userID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, profile) // Use the utility function from pkg/http_response
}

// UpdateProfile handles updating the authenticated user's profile.
// @Router /api/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Use the new helper function from pkg/gin_helpers
	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	// Call the AuthService method to update the user profile.
	profile, appErr := h.authService.UpdateUserProfile(c.Request.Context(), userID, req)
	if appErr != nil {
		http_response.HandleAppError(c, appErr) // Use the utility function from pkg/http_response
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, profile) // Use the utility function from pkg/http_response
}
