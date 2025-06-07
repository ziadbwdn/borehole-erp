package contract

import (
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

// internal/interfaces/repo_interfaces/user_repository.go
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GrantPermission(ctx context.Context, userID, projectID utils.BinaryUUID, permission string) error
	// New methods for profile management
	GetUserByID(ctx context.Context, id utils.BinaryUUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

/**
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) *exception.AppError
	Login(ctx context.Context, req dto.LoginRequest) (string, *exception.AppError)
	RefreshToken(ctx context.Context, token string) (string, *exception.AppError)
	Logout(ctx context.Context, token string) *exception.AppError
	VerifyToken(ctx context.Context, token string) (utils.BinaryUUID, string, *exception.AppError) // Returns userID and role
}
*/
