package services

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/pkg/auth"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 12
)

// AuthServiceImpl implements the contract.AuthService interface.
type AuthServiceImpl struct {
	userRepo    contract.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewAuthService creates a new instance of AuthServiceImpl.
func NewAuthService(userRepo contract.UserRepository, jwtSecret string) contract.AuthService {
	return &AuthServiceImpl{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: 24 * time.Hour, // Default token expiry
	}
}

// Register handles user registration.
func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.ProfileResponse, *exception.AppError) {
	if err := validatePassword(req.Password); err != nil {
		return nil, exception.NewValidationError("password requirements not met", err.Error())
	}

	existingUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// Check if the error is a NotFoundError from the repository.
		var appErr *exception.AppError
		if errors.As(err, &appErr) && appErr.Code == exception.ErrNotFound {
			// User not found, which is expected for a new registration. Continue.
		} else {
			// It's a different kind of error from the database.
			return nil, exception.NewDatabaseError("user lookup failed during registration", err)
		}
	} else if existingUser != nil {
		// User was found, meaning username already exists.
		return nil, exception.NewValidationError("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, exception.NewInternalError("password hashing failed", err)
	}

	user := &models.User{
		ID:           utils.NewBinaryUUID(), // Generate UUID for new user
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         models.RoleGeologist, // Default role for new registrations
		IsActive:     true,
		CreatedAt:    time.Now(), // Set creation timestamp
		UpdatedAt:    time.Now(), // Set update timestamp
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, exception.NewDatabaseError("user creation failed", err)
	}

	return &dto.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Login handles user authentication and token generation.
func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, *exception.AppError) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// Check for NotFoundError specifically.
		var appErr *exception.AppError
		if errors.As(err, &appErr) && appErr.Code == exception.ErrNotFound {
			return nil, exception.NewAuthError("invalid credentials")
		}
		// Other database errors.
		return nil, exception.NewDatabaseError("user lookup failed during login", err)
	}

	if !user.IsActive {
		return nil, exception.NewAuthError("account disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, exception.NewAuthError("invalid credentials")
	}

	// Generate JWT token using the auth utility
	token, err := auth.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return nil, exception.NewInternalError("token generation failed", err)
	}

	return &dto.TokenResponse{
		AccessToken: token,
		ExpiresAt:   time.Now().Add(s.tokenExpiry),
		TokenType:   "Bearer",
	}, nil
}

// RefreshToken re-issues a new access token if the provided token is valid.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, token string) (string, *exception.AppError) {
	// Validate the provided token (assumed to be a refresh token or an expired access token for re-issuance)
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return "", exception.NewAuthError("Invalid or expired refresh token")
	}

	// Re-generate a new access token
	newToken, err := auth.GenerateAccessToken(claims.UserID, claims.Role, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return "", exception.NewInternalError("Failed to generate new access token", err)
	}
	return newToken, nil
}

// VerifyToken validates a JWT token string and returns the UserID and Role from its claims.
func (s *AuthServiceImpl) VerifyToken(ctx context.Context, tokenString string) (utils.BinaryUUID, string, *exception.AppError) {
	// Remove "Bearer " prefix if present.
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	claims, err := auth.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		// Handle specific JWT errors using error checking methods from jwt/v5
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return utils.BinaryUUID{}, "", exception.NewAuthError("Invalid token format")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return utils.BinaryUUID{}, "", exception.NewAuthError("Token expired")
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return utils.BinaryUUID{}, "", exception.NewAuthError("Token not valid yet")
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			return utils.BinaryUUID{}, "", exception.NewAuthError("Invalid token signature")
		}
		// Catch any other parsing errors
		return utils.BinaryUUID{}, "", exception.NewAuthError(fmt.Sprintf("Invalid token: %v", err))
	}

	// Token is valid, return UserID and Role
	return claims.UserID, claims.Role, nil
}

// GetUserProfile retrieves a user's profile by their ID.
func (s *AuthServiceImpl) GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*dto.ProfileResponse, *exception.AppError) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		// GetUserByID returns standard error, so we check if it's an AppError.
		var appErr *exception.AppError
		if errors.As(err, &appErr) {
			return nil, appErr // Propagate NotFoundError or DatabaseError from repo
		}
		return nil, exception.NewInternalError("failed to retrieve user profile", err)
	}

	return &dto.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUserProfile updates a user's profile.
func (s *AuthServiceImpl) UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, *exception.AppError) {
	// First, retrieve the existing user to ensure they exist and to get current values.
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		var appErr *exception.AppError
		if errors.As(err, &appErr) {
			return nil, appErr // Propagate NotFoundError or DatabaseError from repo
		}
		return nil, exception.NewInternalError("failed to retrieve user for update", err)
	}

	// Apply updates only if the corresponding field is provided in the request (not nil).
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Password != nil && *req.Password != "" {
		// Validate and hash the new password if provided.
		if err := validatePassword(*req.Password); err != nil {
			return nil, exception.NewValidationError("new password requirements not met", err.Error())
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, exception.NewInternalError("new password hashing failed", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Call the repository to update the user.
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		var appErr *exception.AppError
		if errors.As(err, &appErr) {
			return nil, appErr // Propagate NotFoundError or DatabaseError from repo
		}
		return nil, exception.NewInternalError("failed to update user profile", err)
	}

	// Return the updated profile response.
	return &dto.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Logout handles token invalidation
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) *exception.AppError {
	// In a real application, this would typically invalidate the token (e.g., add to a blacklist).
	// For this example, we're simply acknowledging the logout.
	return nil // Always returns nil for now, as no actual invalidation is done
}

// validatePassword checks if the password meets the minimum length and complexity requirements.
func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?/", char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(missing, ", "))
	}

	return nil
}
