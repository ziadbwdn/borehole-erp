package repository

import (
	"context"
	"fmt"
	"time"

	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// userRepository is a concrete implementation of the contract.UserRepository interface.
// It holds the GORM database connection.
type userRepository struct {
	db *gorm.DB // Your GORM database connection
}

// NewUserRepository creates and returns a new instance of UserRepository.
// It takes a GORM database connection as a dependency.
func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &userRepository{db: db}
}

// CreateUser inserts a new user record into the database using GORM.
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) *exception.AppError {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return exception.NewDatabaseError(fmt.Sprintf("create user '%s'", user.Username), result.Error)
	}
	return nil
}

// GetUserByUsername retrieves a single user record from the database by their username using GORM.
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, *exception.AppError) {
	user := &models.User{}

	result := r.db.WithContext(ctx).Where("username = ?", username).First(user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("User", username)
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("get user by username '%s'", username), result.Error)
	}
	return user, nil
}

// GetUserByEmail retrieves a single user record from the database by their email using GORM.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, *exception.AppError) {
	user := &models.User{}
	result := r.db.WithContext(ctx).Where("email = ?", email).First(user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("User", email)
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("get user by email '%s'", email), result.Error)
	}
	return user, nil
}

// GetUserByID retrieves a single user record from the database by their ID using GORM.
func (r *userRepository) GetUserByID(ctx context.Context, id utils.BinaryUUID) (*models.User, *exception.AppError) {
	user := &models.User{}

	result := r.db.WithContext(ctx).First(user, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("User", id)
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("get user by ID '%s'", id.String()), result.Error)
	}
	return user, nil
}

// UpdateUser updates an existing user record in the database using GORM.
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) *exception.AppError {
	user.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return exception.NewDatabaseError(fmt.Sprintf("update user '%s'", user.ID.String()), result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("User", user.ID)
	}

	return nil
}

// DeleteUser deletes a user record from the database.
func (r *userRepository) DeleteUser(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return exception.NewDatabaseError(fmt.Sprintf("delete user '%s'", id.String()), result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("User", id)
	}
	return nil
}

// GrantPermission inserts or updates a user's permission for a specific project using GORM.
type UserProjectPermission struct {
	UserID     utils.BinaryUUID `gorm:"type:binary(16);primaryKey"`
	ProjectID  utils.BinaryUUID `gorm:"type:binary(16);primaryKey"`
	Permission string           `gorm:"type:varchar(255);primaryKey"`
	GrantedAt  time.Time
}

func (UserProjectPermission) TableName() string {
	return "user_project_permissions" // Specify the table name
}

func (r *userRepository) GrantPermission(ctx context.Context, userID, projectID utils.BinaryUUID, permission string) *exception.AppError {
	permissionRecord := UserProjectPermission{
		UserID:     userID,
		ProjectID:  projectID,
		Permission: permission,
		GrantedAt:  time.Now(),
	}

	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&permissionRecord)
	if result.Error != nil {
		return exception.NewDatabaseError(
			fmt.Sprintf("grant permission '%s' to user '%s' for project '%s'", permission, userID.String(), projectID.String()),
			result.Error,
		)
	}
	return nil
}

// SaveRefreshToken updates a user's refresh token hash and expiry in the database.
func (r *userRepository) SaveRefreshToken(ctx context.Context, userID utils.BinaryUUID, tokenHash string, expiresAt time.Time) *exception.AppError {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"refresh_token":            tokenHash,
		"refresh_token_expires_at": expiresAt,
		"updated_at":               time.Now(), // Update updated_at
	})
	if result.Error != nil {
		return exception.NewDatabaseError("Failed to save refresh token", result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("User", userID.String())
	}
	return nil
}

// ClearRefreshToken clears a user's refresh token and expiry from the database.
func (r *userRepository) ClearRefreshToken(ctx context.Context, userID utils.BinaryUUID) *exception.AppError {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"refresh_token":            gorm.Expr("NULL"), // Set to NULL
		"refresh_token_expires_at": gorm.Expr("NULL"), // Set to NULL
		"updated_at":               time.Now(),        // Update updated_at
	})
	if result.Error != nil {
		return exception.NewDatabaseError("Failed to clear refresh token", result.Error)
	}
	// If user not found, it means there's nothing to clear, which is fine for logout.
	// But since you want AppError consistency, we'll return NotFound if 0 rows affected.
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("User", userID.String())
	}
	return nil
}

// GetUserByPasswordResetTokenHash retrieves a user by their stored password reset token hash.
func (r *userRepository) GetUserByPasswordResetTokenHash(ctx context.Context, tokenHash string) (*models.User, *exception.AppError) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "password_reset_token = ?", tokenHash).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("User by password reset token", tokenHash)
		}
		return nil, exception.NewDatabaseError("Failed to retrieve user by password reset token hash", err)
	}
	return &user, nil
}

// ClearPasswordResetToken clears a user's password reset token and sent at timestamp.
func (r *userRepository) ClearPasswordResetToken(ctx context.Context, userID utils.BinaryUUID) *exception.AppError {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_reset_token":   gorm.Expr("NULL"),
		"password_reset_sent_at": gorm.Expr("NULL"),
		"updated_at":             time.Now(),
	})
	if result.Error != nil {
		return exception.NewDatabaseError("Failed to clear password reset token", result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("User", userID.String())
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID utils.BinaryUUID) *exception.AppError {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now)

	if result.Error != nil {
		return exception.NewDatabaseError("update last login", result.Error)
	}
	// We don't strictly need to check RowsAffected here, as failing to update a timestamp is not critical.
	// You can add it if you want to be strict.
	return nil
}
