// internal/repo/postgres/user_repo.go
package postgres

import (
	"context"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

// NewUserRepository creates new postgres user repository.
func NewUserRepository(pool *pgxpool.Pool, logger logger.Logger) *userRepository {
	return &userRepository{
		pool:   pool,
		logger: logger,
	}
}

// CreateUser stores a new user (password_hash already hashed).
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// TODO: Execute INSERT query
	// TODO: Handle UNIQUE constraint violation (email) → CONFLICT error
	// TODO: Return created user with ID or error
	return nil, apperror.ErrInternalServer
}

// GetUserByID retrieves a user by ID.
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// TODO: Execute SELECT query by ID
	// TODO: Handle no rows → NOT_FOUND error
	// TODO: Return user or error
	return nil, apperror.ErrInternalServer
}

// GetUserByEmail retrieves a user by email.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	// TODO: Execute SELECT query by email
	// TODO: Handle no rows → NOT_FOUND error
	// TODO: Return user or error
	return nil, apperror.ErrInternalServer
}

// UpdateUser updates user data (excluding password).
func (r *userRepository) UpdateUser(ctx context.Context, id int64, input *entity.UpdateUserInput) (*entity.User, error) {
	// TODO: Execute UPDATE query
	// TODO: Handle no rows → NOT_FOUND error
	// TODO: Return updated user or error
	return nil, apperror.ErrInternalServer
}

// DeleteUser marks user as deleted or removes from database.
func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	// TODO: Execute DELETE query
	// TODO: Handle no rows → NOT_FOUND error
	return apperror.ErrInternalServer
}

// ListUsers returns paginated list of users.
func (r *userRepository) ListUsers(ctx context.Context, limit int, offset int) ([]*entity.User, error) {
	// TODO: Execute SELECT query with LIMIT and OFFSET
	// TODO: Return paginated list or error
	return nil, apperror.ErrInternalServer
}

// UserExists checks if user exists by email.
func (r *userRepository) UserExists(ctx context.Context, email string) (bool, error) {
	// TODO: Execute COUNT query
	// TODO: Return exists bool or error
	return false, apperror.ErrInternalServer
}
