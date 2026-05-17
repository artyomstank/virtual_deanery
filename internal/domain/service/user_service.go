// internal/domain/service/user_service.go
package service

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// UserService defines user business logic interface.
type UserService interface {
	// RegisterUser creates new user with hashed password.
	RegisterUser(ctx context.Context, input *entity.CreateUserInput) (*entity.User, error)

	// LoginUser validates credentials and returns JWT tokens.
	LoginUser(ctx context.Context, email, password string) (*entity.UserToken, error)

	// GetUser retrieves user by ID.
	GetUser(ctx context.Context, id int64) (*entity.User, error)

	// UpdateUser updates user profile.
	UpdateUser(ctx context.Context, id int64, input *entity.UpdateUserInput) (*entity.User, error)

	// DeleteUser removes user.
	DeleteUser(ctx context.Context, id int64) error

	// ListUsers returns paginated users list.
	ListUsers(ctx context.Context, limit int, offset int) ([]*entity.User, error)

	// RefreshAccessToken generates new access token from refresh token.
	RefreshAccessToken(ctx context.Context, refreshToken string) (*entity.UserToken, error)
}
