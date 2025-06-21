package services

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract" // Assuming this is the correct path to your contract
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt" // Correct import for bcrypt
)

const (
	minPasswordLength            = 12
	accessTokenExpiry            = 90 * time.Minute   // Short-lived access token
	refreshTokenExpiry           = 7 * 24 * time.Hour // Long-lived refresh token (e.g., 7 days)
	passwordResetTokenExpiry     = 1 * time.Hour      // Password reset token expiry
	emailVerificationTokenExpiry = 24 * time.Hour     // Email verification token expiry
)

// AuthServiceImpl implements the contract.AuthService interface.
type AuthServiceImpl struct {
	userRepo            contract.UserRepository
	userActivityService contract.UserActivityService // <-- NEW: Dependency for logging activities
	jwtSecret           string
}

// NewAuthService creates a new instance of AuthServiceImpl.
func NewAuthService(userRepo contract.UserRepository, userActivityService contract.UserActivityService, jwtSecret string) contract.AuthService {
	if userRepo == nil {
		panic("userRepo must not be nil for AuthServiceImpl")
	}
	if userActivityService == nil {
		panic("userActivityService must not be nil for AuthServiceImpl") // Defensive check
	}
	return &AuthServiceImpl{
		userRepo:            userRepo,
		userActivityService: userActivityService,
		jwtSecret:           jwtSecret,
	}
}

// Register handles user registration.
func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest, ipAddress string) (*dto.ProfileResponse, *exception.AppError) {
	if err := utils.ValidatePasswordWithRegex(req.Password); err != nil { 
		return nil, exception.NewValidationError("password requirements not met", err.Error())
	}

	_, appErr := s.userRepo.GetUserByUsername(ctx, req.Username)
	if appErr == nil {
		return nil, exception.NewValidationError("username already exists")
	}
	if appErr.Code != exception.ErrNotFound {
		return nil, appErr
	}

	_, appErr = s.userRepo.GetUserByEmail(ctx, req.Email)
	if appErr == nil {
		return nil, exception.NewValidationError("email already exists")
	}
	if appErr.Code != exception.ErrNotFound {
		return nil, appErr
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, exception.NewInternalError("password hashing failed", err)
	}

	userRole := models.UserRole(req.Role)
	switch userRole {
	case models.RoleAdmin, models.RoleGeologist, models.RoleEngineer, models.RoleLabTechnician, models.RoleGuest:
	default:
		return nil, exception.NewValidationError("invalid role specified")
	}

	user := &models.User{
		ID:           utils.NewBinaryUUID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         userRole,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if appErr := s.userRepo.CreateUser(ctx, user); appErr != nil {
		return nil, appErr
	}

	// --- Log successful registration activity ---
	resourceID := user.ID.String()
	details := "New user registered successfully"
	s.userActivityService.LogUserActivity(
		ctx,
		user.ID.String(),
		user.Username,
		models.ActionTypeRegisterUser,
		models.ResourceTypeUser,
		&resourceID,     // Take address of string variable
		&ipAddress,     // Take address of string variable
		&details,       // Take address of string variable
		nil,
		nil,
	)

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

// Login handles user authentication and token generation with specific error feedback.
func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest, ipAddress string) (*dto.TokenResponse, *exception.AppError) {
	// Step 1: Attempt to find the user by username.
	user, appErr := s.userRepo.GetUserByUsername(ctx, req.Username)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			details := fmt.Sprintf("Failed login attempt: account not found for username '%s'", req.Username)
			// --- Log failed login activity: User Not Found ---
			s.userActivityService.LogUserActivity(
				ctx,
				"", // User ID is unknown here
				req.Username,
				models.ActionTypeFailedLogin,
				models.ResourceTypeUser,
				nil,             // Resource ID is unknown, so nil
				&ipAddress,      // Take address of string variable
				&details,        // Take address of string variable
				nil,
				nil,
			)
			return nil, exception.NewAuthError(
				"account not found",
				"please check your username or register",
			)
		}
		return nil, appErr
	}

	// Step 2: Check if the found user's account is active.
	if !user.IsActive {
		details := fmt.Sprintf("Failed login attempt: account disabled for user '%s'", user.Username)
		resourceID := user.ID.String()
		// --- Log failed login activity: Account Disabled ---
		s.userActivityService.LogUserActivity(
			ctx,
			user.ID.String(),
			user.Username,
			models.ActionTypeFailedLogin,
			models.ResourceTypeUser,
			&resourceID,    // Take address of string variable
			&ipAddress,     // Take address of string variable
			&details,       // Take address of string variable
			nil,
			nil,
		)
		return nil, exception.NewAuthError("account is disabled")
	}

	// Step 3: Compare the provided password with the stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		details := fmt.Sprintf("Failed login attempt: incorrect password for user '%s'", user.Username)
		resourceID := user.ID.String()
		// --- Log failed login activity: Wrong Password ---
		s.userActivityService.LogUserActivity(
			ctx,
			user.ID.String(),
			user.Username,
			models.ActionTypeFailedLogin,
			models.ResourceTypeUser,
			&resourceID,    // Take address of string variable
			&ipAddress,     // Take address of string variable
			&details,       // Take address of string variable
			nil,
			nil,
		)
		return nil, exception.NewAuthError(
			"invalid credentials",
			"Incorrect password. Please try again",
		)
	}

	// Step 4: Generate tokens.
	accessToken, err := auth.GenerateAccessToken(user.ID, string(user.Role), s.jwtSecret, accessTokenExpiry)
	if err != nil {
		return nil, exception.NewInternalError("access token generation failed", err)
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID, s.jwtSecret, refreshTokenExpiry)
	if err != nil {
		return nil, exception.NewInternalError("refresh token generation failed", err)
	}

	// Step 5: Hash and save the refresh token.
	refreshTokenHash := auth.HashToken(refreshToken)
	refreshExpiresAt := time.Now().Add(refreshTokenExpiry)
	if appErr = s.userRepo.SaveRefreshToken(ctx, user.ID, refreshTokenHash, refreshExpiresAt); appErr != nil {
		return nil, appErr
	}

	// Step 6: Update last login time safely.
	if updateErr := s.userRepo.UpdateLastLogin(ctx, user.ID); updateErr != nil {
		fmt.Printf("Warning: Failed to update last login for user %s: %v\n", user.ID.String(), updateErr)
	}

	// --- Log successful login activity ---
	resourceID := user.ID.String()
	details := "User logged in successfully"
	s.userActivityService.LogUserActivity(
		ctx,
		user.ID.String(),
		user.Username,
		models.ActionTypeLogin,
		models.ResourceTypeUser,
		&resourceID,    // Take address of string variable
		&ipAddress,     // Take address of string variable
		&details,       // Take address of string variable
		nil,
		nil,
	)

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(accessTokenExpiry),
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken re-issues a new access token and a new refresh token if the provided refresh token is valid.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, *exception.AppError) {
	claims, err := auth.ValidateToken(req.RefreshToken, s.jwtSecret)
    if err != nil {
        return nil, exception.NewAuthError(fmt.Sprintf("Invalid refresh token: %v", err))
    }

    if claims.Type != "refresh" {
        return nil, exception.NewAuthError("Provided token is not a refresh token")
    }

    user, appErr := s.userRepo.GetUserByID(ctx, claims.UserID)
    if appErr != nil {
        if appErr.Code == exception.ErrNotFound {
            return nil, exception.NewAuthError("User not found for refresh token")
        }
        return nil, appErr
    }
    if !user.IsActive {
        return nil, exception.NewAuthError("Account disabled")
    }

    providedRefreshTokenHash := auth.HashToken(req.RefreshToken)
    if user.RefreshToken == nil || *user.RefreshToken != providedRefreshTokenHash {
        _ = s.userRepo.ClearRefreshToken(ctx, user.ID)
        return nil, exception.NewAuthError("Invalid or revoked refresh token")
    }

    if user.RefreshTokenExpiresAt == nil || user.RefreshTokenExpiresAt.Before(time.Now()) {
        _ = s.userRepo.ClearRefreshToken(ctx, user.ID)
        return nil, exception.NewAuthError("Refresh token expired in database")
    }

    newAccessToken, err := auth.GenerateAccessToken(user.ID, string(user.Role), s.jwtSecret, accessTokenExpiry)
    if err != nil {
        return nil, exception.NewInternalError("failed to generate new access token", err)
    }

    newRefreshToken, err := auth.GenerateRefreshToken(user.ID, s.jwtSecret, refreshTokenExpiry)
    if err != nil {
        return nil, exception.NewInternalError("failed to generate new refresh token", err)
    }

    newRefreshTokenHash := auth.HashToken(newRefreshToken)
    newRefreshExpiresAt := time.Now().Add(refreshTokenExpiry)
    if appErr = s.userRepo.SaveRefreshToken(ctx, user.ID, newRefreshTokenHash, newRefreshExpiresAt); appErr != nil {
        return nil, appErr
    }

    // --- Log activity: Token Refreshed ---
    resourceID := user.ID.String()
    details := "User tokens refreshed successfully"
    s.userActivityService.LogUserActivity(
        ctx,
        user.ID.String(),
        user.Username,
        models.ActionTypeTokenRefresh,
        models.ResourceTypeUser,
        &resourceID, // Take address of string variable
        nil,         // IP address not directly available here
        &details,    // Take address of string variable
        nil,
        nil,
    )

    return &dto.RefreshResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
        ExpiresAt:    time.Now().Add(accessTokenExpiry),
        TokenType:    "Bearer",
    }, nil
}

// VerifyToken validates a JWT token string and returns the UserID and Role from its claims.
func (s *AuthServiceImpl) VerifyToken(ctx context.Context, tokenString string) (utils.BinaryUUID, string, *exception.AppError) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims, err := auth.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return utils.BinaryUUID{}, "", exception.NewAuthError(fmt.Sprintf("Invalid token: %v", err))
	}
	if claims.Type != "access" {
		return utils.BinaryUUID{}, "", exception.NewAuthError("Provided token is not an access token")
	}

	return claims.UserID, claims.Role, nil
}

// GetUserProfile retrieves a user's profile by their ID.
func (s *AuthServiceImpl) GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*dto.ProfileResponse, *exception.AppError) {
	user, appErr := s.userRepo.GetUserByID(ctx, userID)
	if appErr != nil {
		return nil, appErr
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
	user, appErr := s.userRepo.GetUserByID(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	oldEmail := user.Email
	oldFullName := user.FullName
	oldPasswordHash := user.PasswordHash

	changedFields := []string{}
	oldValuesMap := make(map[string]interface{})
	newValuesMap := make(map[string]interface{})

	if req.Email != nil && *req.Email != user.Email {
		oldValuesMap["email"] = oldEmail // Use stored old value
		user.Email = *req.Email
		newValuesMap["email"] = user.Email
		changedFields = append(changedFields, "email")
	}
	if req.FullName != nil && *req.FullName != user.FullName {
		oldValuesMap["full_name"] = oldFullName // Use stored old value
		user.FullName = *req.FullName
		newValuesMap["full_name"] = user.FullName
		changedFields = append(changedFields, "full_name")
	}
	if req.Password != nil && *req.Password != "" {
		if err := utils.ValidatePasswordWithRegex(*req.Password); err != nil {
			return nil, exception.NewValidationError("new password requirements not met", err.Error())
		}
		oldValuesMap["password_hash"] = oldPasswordHash
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, exception.NewInternalError("new password hashing failed", hashErr)
		}
		user.PasswordHash = string(hashedPassword)
		newValuesMap["password_hash"] = user.PasswordHash
		changedFields = append(changedFields, "password")
	}

	if len(changedFields) == 0 {
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

	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
		return nil, appErr
	}

	var oldValJSON, newValJSON *string
	if len(oldValuesMap) > 0 {
		oldBytes, err := json.Marshal(oldValuesMap)
		if err == nil {
			s := string(oldBytes)
			oldValJSON = &s // Take address of temporary string
		}
	}
	if len(newValuesMap) > 0 {
		newBytes, err := json.Marshal(newValuesMap)
		if err == nil {
			s := string(newBytes)
			newValJSON = &s // Take address of temporary string
		}
	}

	// --- Log successful profile update activity ---
	details := fmt.Sprintf("User profile updated. Changed fields: %s", strings.Join(changedFields, ", "))
	resourceID := user.ID.String()
	s.userActivityService.LogUserActivity(
		ctx,
		user.ID.String(),
		user.Username,
		models.ActionTypeUpdateProfile,
		models.ResourceTypeUser,
		&resourceID, // Take address of string variable
		nil,         // IP address not directly available here
		&details,    // Take address of string variable
		oldValJSON,
		newValJSON,
	)

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

// GET User Details For Logging
func (s *AuthServiceImpl) GetUserDetailsForLogging(ctx context.Context, userID utils.BinaryUUID) (string, *exception.AppError) {
    user, appErr := s.userRepo.GetUserByID(ctx, userID)
    if appErr != nil {
        // It's an internal error if we have a valid UserID from a token but can't find that user.
        return "", exception.NewInternalError("Failed to retrieve user details for logging", appErr)
    }
    return user.Username, nil
}


// SendPasswordReset sends a password reset email to the user.
func (s *AuthServiceImpl) SendPasswordReset(ctx context.Context, email string) *exception.AppError {
	user, appErr := s.userRepo.GetUserByEmail(ctx, email)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			fmt.Printf("Attempted password reset for non-existent email: %s\n", email)
			return nil
		}
		return appErr
	}

	resetToken, err := auth.GeneratePasswordResetToken(user.ID, s.jwtSecret, passwordResetTokenExpiry)
	if err != nil {
		return exception.NewInternalError("failed to generate password reset token", err)
	}

	hashedResetToken := auth.HashToken(resetToken)
	tokenExpiry := time.Now().Add(passwordResetTokenExpiry)

	user.PasswordResetToken = &hashedResetToken
	user.PasswordResetSentAt = &tokenExpiry

	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
		return appErr
	}

	fmt.Printf("Password reset link for %s: YOUR_FRONTEND_URL/reset-password?token=%s\n", email, resetToken)

	// --- Log activity: Password Reset Initiated ---
	resourceID := user.ID.String()
	details := fmt.Sprintf("Password reset initiated for user %s", user.Username)
	s.userActivityService.LogUserActivity(
		ctx,
		user.ID.String(),
		user.Username,
		models.ActionTypePasswordResetRequest,
		models.ResourceTypeUser,
		&resourceID, // Take address of string variable
		nil,         // IP address not directly available here
		&details,    // Take address of string variable
		nil,
		nil,
	)

	return nil
}

// ResetPassword resets the user's password using a valid token.
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, token, newPw string) *exception.AppError {
	if err := utils.ValidatePasswordWithRegex(newPw); err != nil {
		return exception.NewValidationError("new password requirements not met", err.Error())
	}

	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return exception.NewAuthError(fmt.Sprintf("Invalid password reset token: %v", err))
	}

	if claims.Type != "reset" {
		return exception.NewAuthError("Provided token is not a password reset token")
	}

	user, appErr := s.userRepo.GetUserByID(ctx, claims.UserID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return exception.NewAuthError("User not found for password reset")
		}
		return appErr
	}

	providedTokenHash := auth.HashToken(token)
	if user.PasswordResetToken == nil || *user.PasswordResetToken != providedTokenHash {
		return exception.NewAuthError("Invalid or already used password reset token")
	}

	if user.PasswordResetSentAt == nil || user.PasswordResetSentAt.Before(time.Now()) {
		_ = s.userRepo.ClearPasswordResetToken(ctx, user.ID)
		return exception.NewAuthError("Password reset token expired in database")
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(newPw), bcrypt.DefaultCost)
	if hashErr != nil {
		return exception.NewInternalError("failed to hash new password", hashErr)
	}
	user.PasswordHash = string(hashedPassword)

	if appErr := s.userRepo.ClearPasswordResetToken(ctx, user.ID); appErr != nil {
		fmt.Printf("Warning: Failed to clear password reset token for user %s: %v\n", user.ID.String(), appErr)
	}

	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
		return appErr
	}

	// --- Log activity: Password Reset Successful ---
	resourceID := user.ID.String()
	details := "User password reset successfully"
	s.userActivityService.LogUserActivity(
		ctx,
		user.ID.String(),
		user.Username,
		models.ActionTypePasswordReset,
		models.ResourceTypeUser,
		&resourceID, // Take address of string variable
		nil,         // IP address not directly available here
		&details,    // Take address of string variable
		nil,
		nil,
	)

	return nil
}

// VerifyEmail marks a user's email as verified.
func (s *AuthServiceImpl) VerifyEmail(ctx context.Context, token string, ) *exception.AppError {
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return exception.NewAuthError(fmt.Sprintf("Invalid email verification token: %v", err))
	}

	if claims.Type != "email_verify" {
		return exception.NewAuthError("Provided token is not an email verification token")
	}

	user, appErr := s.userRepo.GetUserByID(ctx, claims.UserID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return exception.NewAuthError("User not found for email verification")
		}
		return appErr
	}

	if user.EmailVerified {
		return exception.NewValidationError("Email already verified")
	}

	user.EmailVerified = true
	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
		return appErr
	}

	// --- Log activity: Email Verified ---
	resourceID := user.ID.String()
	details := "User email verified successfully"
	s.userActivityService.LogUserActivity(
		ctx,
		user.ID.String(),
		user.Username,
		models.ActionTypeEmailVerified,
		models.ResourceTypeUser,
		&resourceID, // Take address of string variable
		nil,         // IP address not directly available here
		&details,    // Take address of string variable
		nil,
		nil,
	)

	return nil
}

// Logout handles token invalidation (clearing refresh token from DB).
func (s *AuthServiceImpl) Logout(ctx context.Context, token string, ipAddress string) *exception.AppError {
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		fmt.Printf("Warning: Attempted logout with invalid token: %v\n", err)
		return exception.NewAuthError("Invalid token provided for logout")
	}

	if claims.Type != "refresh" {
		fmt.Printf("Warning: Attempted logout with non-refresh token for user %s\n", claims.UserID.String())
		return exception.NewAuthError("Provided token is not a refresh token")
	}

	if appErr := s.userRepo.ClearRefreshToken(ctx, claims.UserID); appErr != nil {
		return appErr
	}

	user, appErr := s.userRepo.GetUserByID(ctx, claims.UserID)
	username := "unknown"
	if appErr == nil {
		username = user.Username
	} else {
		fmt.Printf("Warning: Could not retrieve username for logout activity log for user ID %s: %v\n", claims.UserID.String(), appErr)
	}

	// --- Log activity: Logout ---
	resourceID := claims.UserID.String()
	details := "User logged out successfully"
	s.userActivityService.LogUserActivity(
		ctx,
		claims.UserID.String(),
		username,
		models.ActionTypeLogout,
		models.ResourceTypeUser,
		&resourceID, // Take address of string variable
		&ipAddress,         // IP address not directly available here
		&details,    // Take address of string variable
		nil,
		nil,
	)

	return nil
}
