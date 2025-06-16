package contract

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils"
	"context"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.ProfileResponse, *exception.AppError)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, *exception.AppError)
	VerifyToken(ctx context.Context, tokenString string) (utils.BinaryUUID, string, *exception.AppError)
	GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*dto.ProfileResponse, *exception.AppError)
	UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, *exception.AppError)

	// **** THIS IS THE CRITICAL CHANGE ****
	// The RefreshToken signature in the contract MUST match the implementation.
	RefreshToken(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, *exception.AppError)

	SendPasswordReset(ctx context.Context, email string) *exception.AppError
	ResetPassword(ctx context.Context, token, newPassword string) *exception.AppError
	VerifyEmail(ctx context.Context, token string) *exception.AppError
	Logout(ctx context.Context, token string) *exception.AppError
}
