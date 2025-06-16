package services

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract" // Assuming this is the correct path to your contract
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/pkg/auth"
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt" // Correct import for bcrypt
)

const (
	minPasswordLength            = 12
	accessTokenExpiry            = 15 * time.Minute   // Short-lived access token
	refreshTokenExpiry           = 7 * 24 * time.Hour // Long-lived refresh token (e.g., 7 days)
	passwordResetTokenExpiry     = 1 * time.Hour      // Password reset token expiry
	emailVerificationTokenExpiry = 24 * time.Hour     // Email verification token expiry
)

// AuthServiceImpl implements the contract.AuthService interface.
type AuthServiceImpl struct {
	userRepo  contract.UserRepository
	jwtSecret string
}

// NewAuthService creates a new instance of AuthServiceImpl.
func NewAuthService(userRepo contract.UserRepository, jwtSecret string) contract.AuthService {
	return &AuthServiceImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration.
func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.ProfileResponse, *exception.AppError) {
	if err := validatePassword(req.Password); err != nil {
		return nil, exception.NewValidationError("password requirements not met", err.Error())
	}

	// Check if username already exists
	_, appErr := s.userRepo.GetUserByUsername(ctx, req.Username)
	if appErr == nil { // User found, meaning username already exists
		return nil, exception.NewValidationError("username already exists")
	}
	if appErr.Code != exception.ErrNotFound { // If it's not a NotFoundError, it's a database error
		return nil, appErr // Propagate the database error directly
	}

	// Check if email already exists
	_, appErr = s.userRepo.GetUserByEmail(ctx, req.Email)
	if appErr == nil { // User found, meaning email already exists
		return nil, exception.NewValidationError("email already exists")
	}
	if appErr.Code != exception.ErrNotFound { // If it's not a NotFoundError, it's a database error
		return nil, appErr // Propagate the database error directly
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// **** FIX for bcrypt error ****
		return nil, exception.NewInternalError("password hashing failed", err)
	}

	userRole := models.UserRole(req.Role)
	switch userRole {
	case models.RoleAdmin, models.RoleGeologist, models.RoleEngineer, models.RoleLabTechnician, models.RoleGuest:
		// Valid roles
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
func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, *exception.AppError) {
	// Step 1: Attempt to find the user by username.
	user, appErr := s.userRepo.GetUserByUsername(ctx, req.Username)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			// CASE 1: USER NOT FOUND
			// Return the specific error message as requested.
			return nil, exception.NewAuthError(
				"account not found",
				"please check your username or register",
			)
		}
		// Propagate other errors (e.g., database connection issue).
		return nil, appErr
	}

	// Step 2: Check if the found user's account is active.
	if !user.IsActive {
		// This error is also specific and correct.
		return nil, exception.NewAuthError("account is disabled")
	}

	// Step 3: Compare the provided password with the stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// CASE 2: WRONG PASSWORD
		// The user was found, but the password was incorrect. Return the specific error.
		return nil, exception.NewAuthError(
			"invalid credentials",
			"Incorrect password. Please try again",
		)
	}

	// --- If we reach here, the login is successful ---
	// The rest of the logic remains the same as it's already correct.

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

	// Step 7: Return the successful token response.
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(accessTokenExpiry),
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken re-issues a new access token and a new refresh token if the provided refresh token is valid.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, *exception.AppError) {
	// 1. Validate the provided Refresh Token (JWT validation)
	claims, err := auth.ValidateToken(req.RefreshToken, s.jwtSecret)
	if err != nil {
		return nil, exception.NewAuthError(fmt.Sprintf("Invalid refresh token: %v", err))
	}

	// Ensure it's a refresh token type
	if claims.Type != "refresh" {
		return nil, exception.NewAuthError("Provided token is not a refresh token")
	}

	// 2. Retrieve the user from the database
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

	// 3. Compare the provided refresh token's SHA256 hash with the hash stored in the database
	providedRefreshTokenHash := auth.HashToken(req.RefreshToken)
	if user.RefreshToken == nil || *user.RefreshToken != providedRefreshTokenHash {
		_ = s.userRepo.ClearRefreshToken(ctx, user.ID)
		return nil, exception.NewAuthError("Invalid or revoked refresh token")
	}

	// 4. Check if the refresh token in the database has expired
	if user.RefreshTokenExpiresAt == nil || user.RefreshTokenExpiresAt.Before(time.Now()) {
		_ = s.userRepo.ClearRefreshToken(ctx, user.ID)
		return nil, exception.NewAuthError("Refresh token expired in database")
	}

	// 5. Generate new Access Token
	newAccessToken, err := auth.GenerateAccessToken(user.ID, string(user.Role), s.jwtSecret, accessTokenExpiry)
	if err != nil {
		return nil, exception.NewInternalError("failed to generate new access token", err)
	}

	// 6. Generate new Refresh Token (Refresh Token Rotation)
	newRefreshToken, err := auth.GenerateRefreshToken(user.ID, s.jwtSecret, refreshTokenExpiry)
	if err != nil {
		return nil, exception.NewInternalError("failed to generate new refresh token", err)
	}

	// 7. Hash and save the new refresh token to the database
	newRefreshTokenHash := auth.HashToken(newRefreshToken)
	newRefreshExpiresAt := time.Now().Add(refreshTokenExpiry)
	if appErr = s.userRepo.SaveRefreshToken(ctx, user.ID, newRefreshTokenHash, newRefreshExpiresAt); appErr != nil {
		return nil, appErr
	}

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
		// **** FIX for auth.GenerateToken typo -> using auth.ValidateToken ****
		return utils.BinaryUUID{}, "", exception.NewAuthError(fmt.Sprintf("Invalid token: %v", err))
	}
	// Ensure it's an access token
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

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Password != nil && *req.Password != "" {
		if err := validatePassword(*req.Password); err != nil {
			return nil, exception.NewValidationError("new password requirements not met", err.Error())
		}
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, exception.NewInternalError("new password hashing failed", hashErr)
		}
		user.PasswordHash = string(hashedPassword)
	}

	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
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

	// Generate a JWT for password reset
	resetToken, err := auth.GeneratePasswordResetToken(user.ID, s.jwtSecret, passwordResetTokenExpiry)
	if err != nil {
		return exception.NewInternalError("failed to generate password reset token", err)
	}

	// Hash the JWT reset token using SHA256 for storage
	hashedResetToken := auth.HashToken(resetToken)
	tokenExpiry := time.Now().Add(passwordResetTokenExpiry)

	user.PasswordResetToken = &hashedResetToken
	user.PasswordResetSentAt = &tokenExpiry

	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
		return appErr
	}

	// TODO: Integrate with an email sending service here
	fmt.Printf("Password reset link for %s: YOUR_FRONTEND_URL/reset-password?token=%s\n", email, resetToken)

	return nil
}

// ResetPassword resets the user's password using a valid token.
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, token, newPw string) *exception.AppError {
	if err := validatePassword(newPw); err != nil {
		return exception.NewValidationError("new password requirements not met", err.Error())
	}

	// 1. Validate the provided reset token (JWT validation)
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return exception.NewAuthError(fmt.Sprintf("Invalid password reset token: %v", err))
	}

	// Ensure it's a reset token type
	if claims.Type != "reset" {
		return exception.NewAuthError("Provided token is not a password reset token")
	}

	// 2. Retrieve the user from the database
	user, appErr := s.userRepo.GetUserByID(ctx, claims.UserID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return exception.NewAuthError("User not found for password reset")
		}
		return appErr
	}

	// 3. Compare the provided token's SHA256 hash with the hash stored in the database
	providedTokenHash := auth.HashToken(token)
	if user.PasswordResetToken == nil || *user.PasswordResetToken != providedTokenHash {
		return exception.NewAuthError("Invalid or already used password reset token")
	}

	// 4. Check if the reset token in the database has expired (redundant if JWT validation caught it but good for robustness)
	if user.PasswordResetSentAt == nil || user.PasswordResetSentAt.Before(time.Now()) {
		_ = s.userRepo.ClearPasswordResetToken(ctx, user.ID)
		return exception.NewAuthError("Password reset token expired in database")
	}

	// Hash new password
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(newPw), bcrypt.DefaultCost)
	if hashErr != nil {
		return exception.NewInternalError("failed to hash new password", hashErr)
	}
	user.PasswordHash = string(hashedPassword)

	// Invalidate the reset token after successful use
	if appErr := s.userRepo.ClearPasswordResetToken(ctx, user.ID); appErr != nil {
		fmt.Printf("Warning: Failed to clear password reset token for user %s: %v\n", user.ID.String(), appErr)
	}

	// Update user's password (this implicitly updates UpdatedAt)
	if appErr := s.userRepo.UpdateUser(ctx, user); appErr != nil {
		return appErr
	}

	return nil
}

// VerifyEmail marks a user's email as verified.
func (s *AuthServiceImpl) VerifyEmail(ctx context.Context, token string) *exception.AppError {
	// 1. Validate the provided email verification token (JWT validation)
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return exception.NewAuthError(fmt.Sprintf("Invalid email verification token: %v", err))
	}

	// Ensure it's an email verification token type
	if claims.Type != "email_verify" {
		return exception.NewAuthError("Provided token is not an email verification token")
	}

	// 2. Retrieve the user from the database
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
	// TODO: Consider clearing the email verification token hash from DB after successful use if you store it.

	return nil
}

// Logout handles token invalidation (clearing refresh token from DB).
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) *exception.AppError {
	// We assume the 'token' provided here is the refresh token the client wants to invalidate.
	// 1. Validate the provided Refresh Token (JWT validation) to get UserID
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		fmt.Printf("Warning: Attempted logout with invalid token: %v\n", err)
		return exception.NewAuthError("Invalid token provided for logout")
	}

	// Ensure it's a refresh token type
	if claims.Type != "refresh" {
		fmt.Printf("Warning: Attempted logout with non-refresh token for user %s\n", claims.UserID.String())
		return exception.NewAuthError("Provided token is not a refresh token")
	}

	// 2. Clear the specific refresh token hash from the database for this user
	if appErr := s.userRepo.ClearRefreshToken(ctx, claims.UserID); appErr != nil {
		return appErr
	}

	return nil
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
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?/", char): // Common special characters
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
