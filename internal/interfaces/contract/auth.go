package contract

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils"
	"context"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	// Add ipAddress parameter
	Register(ctx context.Context, req dto.RegisterRequest, ipAddress string) (*dto.ProfileResponse, *exception.AppError)
	// Add ipAddress parameter
	Login(ctx context.Context, req dto.LoginRequest, ipAddress string) (*dto.TokenResponse, *exception.AppError)

	VerifyToken(ctx context.Context, tokenString string) (utils.BinaryUUID, string, *exception.AppError)
	GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*dto.ProfileResponse, *exception.AppError)
	UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, *exception.AppError)
	GetUserDetailsForLogging(ctx context.Context, userID utils.BinaryUUID) (username string, appErr *exception.AppError)
	RefreshToken(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, *exception.AppError)

	SendPasswordReset(ctx context.Context, email string) *exception.AppError
	ResetPassword(ctx context.Context, token, newPassword string) *exception.AppError
	VerifyEmail(ctx context.Context, token string) *exception.AppError
	Logout(ctx context.Context, token string, ipAddress string) *exception.AppError
}
