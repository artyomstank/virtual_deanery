// internal/repo/memory/user_repo.go
package memory

import (
	"context"
	"sync"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// userRepository is in-memory implementation for testing.
type userRepository struct {
	mu     sync.RWMutex
	users  map[int64]*entity.User
	nextID int64
}

// NewUserRepository creates new in-memory user repository.
func NewUserRepository() *userRepository {
	return &userRepository{
		users:  make(map[int64]*entity.User),
		nextID: 1,
	}
}

// CreateUser stores a new user (password_hash already hashed).
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// TODO: Check if user with same email exists

	// TODO: Assign ID and store

	// TODO: Return created user or error
	return nil, apperror.ErrInternalServer
}

// GetUserByID retrieves a user by ID.
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// TODO: Find user by ID

	// TODO: Return user or NOT_FOUND error
	return nil, apperror.ErrNotFound
}

// GetUserByEmail retrieves a user by email.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	// TODO: Find user by email

	// TODO: Return user or NOT_FOUND error
	return nil, apperror.ErrNotFound
}

// UpdateUser updates user data (excluding password).
func (r *userRepository) UpdateUser(ctx context.Context, id int64, input *entity.UpdateUserInput) (*entity.User, error) {
	// TODO: Find user by ID

	// TODO: Update fields

	// TODO: Return updated user or NOT_FOUND error
	return nil, apperror.ErrNotFound
}

// DeleteUser marks user as deleted or removes from database.
func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	// TODO: Remove user from map

	// TODO: Handle NOT_FOUND error
	return apperror.ErrNotFound
}

// ListUsers returns paginated list of users.
func (r *userRepository) ListUsers(ctx context.Context, limit int, offset int) ([]*entity.User, error) {
	// TODO: Return paginated slice from map

	// TODO: Handle limit and offset
	return nil, nil
}

// UserExists checks if user exists by email.
func (r *userRepository) UserExists(ctx context.Context, email string) (bool, error) {
	// TODO: Check if user with email exists

	return false, nil
}
