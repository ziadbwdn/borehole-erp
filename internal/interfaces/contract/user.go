package contract

import (
	"boreholedata-ms/internal/exception" // Import your custom exception package
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) *exception.AppError
	GetUserByID(ctx context.Context, id utils.BinaryUUID) (*models.User, *exception.AppError)
	GetUserByUsername(ctx context.Context, username string) (*models.User, *exception.AppError)
	GetUserByEmail(ctx context.Context, email string) (*models.User, *exception.AppError)
	UpdateUser(ctx context.Context, user *models.User) *exception.AppError
	DeleteUser(ctx context.Context, id utils.BinaryUUID) *exception.AppError

	// New methods for refresh tokens and reset tokens
	SaveRefreshToken(ctx context.Context, userID utils.BinaryUUID, tokenHash string, expiresAt time.Time) *exception.AppError
	// UpdateRefreshToken(ctx context.Context, userID utils.BinaryUUID, refreshTokenHash string, expiresAt time.Time) *exception.AppError
	ClearRefreshToken(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	GetUserByPasswordResetTokenHash(ctx context.Context, tokenHash string) (*models.User, *exception.AppError)
	ClearPasswordResetToken(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	UpdateLastLogin(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	GrantPermission(ctx context.Context, userID, projectID utils.BinaryUUID, permission string) *exception.AppError
}
