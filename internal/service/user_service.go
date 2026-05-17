// internal/service/user_service.go
package service

import (
	"context"
	"errors"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

type userService struct {
	repo       repository.UserRepository
	jwtClient  jwt.TokenClient
	logger     logger.Logger
	bcryptCost int
}

// NewUserService creates new user service.
func NewUserService(
	repo repository.UserRepository,
	jwtClient jwt.TokenClient,
	logger logger.Logger,
	bcryptCost int,
) *userService {
	return &userService{
		repo:       repo,
		jwtClient:  jwtClient,
		logger:     logger,
		bcryptCost: bcryptCost,
	}
}

// RegisterUser creates new user with hashed password.
func (s *userService) RegisterUser(ctx context.Context, input *entity.CreateUserInput) (*entity.User, error) {
	// TODO: Validate input

	// TODO: Check if user already exists

	// TODO: Hash password using bcrypt with cost s.bcryptCost

	// TODO: Create user in repository

	// TODO: Return user or handle error
	return nil, errors.New("not implemented")
}

// LoginUser validates credentials and returns JWT tokens.
func (s *userService) LoginUser(ctx context.Context, email, password string) (*entity.UserToken, error) {
	// TODO: Get user from repository by email

	// TODO: Validate password against stored hash using bcrypt

	// TODO: Generate access and refresh tokens using s.jwtClient

	// TODO: Return tokens or handle error
	return nil, errors.New("not implemented")
}

// GetUser retrieves user by ID.
func (s *userService) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	// TODO: Get user from repository

	// TODO: Handle NOT_FOUND error

	return nil, errors.New("not implemented")
}

// UpdateUser updates user profile.
func (s *userService) UpdateUser(ctx context.Context, id int64, input *entity.UpdateUserInput) (*entity.User, error) {
	// TODO: Check user exists

	// TODO: Update user in repository

	// TODO: Return updated user or handle error
	return nil, errors.New("not implemented")
}

// DeleteUser removes user.
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	// TODO: Check user exists

	// TODO: Delete user from repository

	// TODO: Handle errors
	return errors.New("not implemented")
}

// ListUsers returns paginated users list.
func (s *userService) ListUsers(ctx context.Context, limit int, offset int) ([]*entity.User, error) {
	// TODO: Validate pagination params

	// TODO: Get users from repository

	// TODO: Return paginated list or handle error
	return nil, errors.New("not implemented")
}

// RefreshAccessToken generates new access token from refresh token.
func (s *userService) RefreshAccessToken(ctx context.Context, refreshToken string) (*entity.UserToken, error) {
	// TODO: Validate refresh token using s.jwtClient

	// TODO: Extract user ID from refresh token claims

	// TODO: Get user to ensure still exists

	// TODO: Generate new access token

	// TODO: Return new token pair or handle error
	return nil, errors.New("not implemented")
}
