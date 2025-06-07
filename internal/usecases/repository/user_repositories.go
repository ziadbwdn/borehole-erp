package repository

import (
	"context"
	"fmt"
	"time" // Needed for CreatedAt/UpdatedAt fields if not handled by GORM hooks

	"boreholedata-ms/internal/exception"           // Import your custom exception package
	"boreholedata-ms/internal/interfaces/contract" // Import your user contract interface
	"boreholedata-ms/internal/models"              // Import your user model
	"boreholedata-ms/internal/utils"               // Import your utility types like BinaryUUID

	"gorm.io/gorm"        // Import GORM
	"gorm.io/gorm/clause" // For ON CONFLICT DO NOTHING
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
// It takes a context and a User model, returning an error if the operation fails.
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	// GORM will automatically handle setting CreatedAt/UpdatedAt if your models
	// embed gorm.Model or have fields named CreatedAt/UpdatedAt with time.Time type.
	// If not, you might need to set them manually before Create.
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		// Use your custom NewDatabaseError constructor for GORM errors.
		return exception.NewDatabaseError(fmt.Sprintf("create user '%s'", user.Username), result.Error)
	}
	return nil
}

// GetUserByUsername retrieves a single user record from the database by their username using GORM.
// It takes a context and a username string, returning the User model or an error.
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{} // Initialize an empty User struct to hold the results.

	// Use GORM's Where and First methods to find a user by username.
	result := r.db.WithContext(ctx).Where("username = ?", username).First(user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return a specific error if no user is found, which can be handled by the service layer.
			// Use NewNotFoundError for this case.
			return nil, exception.NewNotFoundError("User", username)
		}
		// Wrap other GORM errors using NewDatabaseError.
		return nil, exception.NewDatabaseError(fmt.Sprintf("get user by username '%s'", username), result.Error)
	}
	return user, nil
}

// GrantPermission inserts or updates a user's permission for a specific project using GORM.
// This assumes a separate linking table (e.g., user_project_permissions) for many-to-many relationships.
// For simplicity, we define a temporary struct for the permission record.
// You might have a dedicated GORM model for this table.
type UserProjectPermission struct {
	UserID     utils.BinaryUUID `gorm:"type:binary(16);primaryKey"`
	ProjectID  utils.BinaryUUID `gorm:"type:binary(16);primaryKey"`
	Permission string           `gorm:"type:varchar(255);primaryKey"`
	GrantedAt  time.Time
}

func (UserProjectPermission) TableName() string {
	return "user_project_permissions" // Specify the table name
}

func (r *userRepository) GrantPermission(ctx context.Context, userID, projectID utils.BinaryUUID, permission string) error {
	permissionRecord := UserProjectPermission{
		UserID:     userID,
		ProjectID:  projectID,
		Permission: permission,
		GrantedAt:  time.Now(),
	}

	// Use GORM's Clauses for ON CONFLICT DO NOTHING behavior.
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&permissionRecord)
	if result.Error != nil {
		return exception.NewDatabaseError(
			fmt.Sprintf("grant permission '%s' to user '%s' for project '%s'", permission, userID.String(), projectID.String()),
			result.Error,
		)
	}
	return nil
}

// GetUserByID retrieves a single user record from the database by their ID using GORM.
// It takes a context and a user ID, returning the User model or an error.
func (r *userRepository) GetUserByID(ctx context.Context, id utils.BinaryUUID) (*models.User, error) {
	user := &models.User{} // Initialize an empty User struct to hold the results.

	// Use GORM's First method to find a user by ID.
	result := r.db.WithContext(ctx).First(user, "id = ?", id) // GORM automatically maps ID to primary key

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return a specific error if no user is found.
			return nil, exception.NewNotFoundError("User", id)
		}
		// Wrap other GORM errors using NewDatabaseError.
		return nil, exception.NewDatabaseError(fmt.Sprintf("get user by ID '%s'", id.String()), result.Error)
	}
	return user, nil
}

// UpdateUser updates an existing user record in the database using GORM.
// It takes a context and a User model with updated fields, returning an error if the operation fails.
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	// Ensure the UpdatedAt timestamp is set before updating.
	user.UpdatedAt = time.Now()

	// Use GORM's Save method to update the user.
	// Save will perform an UPDATE if the primary key exists, otherwise an INSERT.
	// We assume the user.ID is already set for an update operation.
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		// Wrap GORM errors using NewDatabaseError.
		return exception.NewDatabaseError(fmt.Sprintf("update user '%s'", user.ID.String()), result.Error)
	}

	// Check if any rows were affected. If not, it might indicate the user was not found.
	if result.RowsAffected == 0 {
		return exception.NewNotFoundError("User", user.ID)
	}

	return nil
}
