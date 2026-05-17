// internal/domain/repository/user_repository.go
package repository

import (
	"context"

	"github.com/artyomstank/go_auth_template/internal/domain/entity"
)

// UserRepository defines methods for user data access.
type UserRepository interface {
	// CreateUser stores a new user (password_hash is already hashed).
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)

	// GetUserByID retrieves a user by ID.
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)

	// GetUserByEmail retrieves a user by email.
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateUser updates user data (excluding password).
	UpdateUser(ctx context.Context, id int64, input *entity.UpdateUserInput) (*entity.User, error)

	// DeleteUser marks user as deleted or removes from database.
	DeleteUser(ctx context.Context, id int64) error

	// ListUsers returns paginated list of users.
	ListUsers(ctx context.Context, limit int, offset int) ([]*entity.User, error)

	// UserExists checks if user exists by email.
	UserExists(ctx context.Context, email string) (bool, error)
}
