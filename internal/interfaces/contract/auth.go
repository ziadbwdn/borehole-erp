package contract

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.ProfileResponse, *exception.AppError)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, *exception.AppError) // Corrected return type
	RefreshToken(ctx context.Context, token string) (string, *exception.AppError)
	Logout(ctx context.Context, token string) *exception.AppError
	VerifyToken(ctx context.Context, token string) (utils.BinaryUUID, string, *exception.AppError) // Returns userID and role

	// New methods for profile management
	GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*dto.ProfileResponse, *exception.AppError)
	UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, *exception.AppError)
}
