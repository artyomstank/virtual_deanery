// internal/transport/http/handler/user_handler.go
package handler

import (
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    service.UserService
	logger logger.Logger
}

// NewUserHandler creates new user HTTP handler.
func NewUserHandler(svc service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger,
	}
}

// RegisterUserDTO is request DTO for user registration.
type RegisterUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserDTO is request DTO for login.
type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterUser handles POST /api/v1/users/register.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	// TODO: Parse RegisterUserDTO from request body

	// TODO: Validate DTO using validator

	// TODO: Call h.svc.RegisterUser()

	// TODO: Map apperror to HTTP status (CONFLICT→409, BAD_REQUEST→400)

	// TODO: Return 201 Created with user data
}

// LoginUser handles POST /api/v1/users/login.
func (h *UserHandler) LoginUser(c *gin.Context) {
	// TODO: Parse LoginUserDTO from request body

	// TODO: Validate DTO using validator

	// TODO: Call h.svc.LoginUser()

	// TODO: Set refreshToken in httpOnly secure cookie
	// name="refresh_token", HttpOnly, Secure, SameSite=Strict

	// TODO: Return 200 OK with access token in body

	// TODO: Handle errors (NOT_FOUND→404, INVALID_CREDENTIALS→401)
}

// GetUser handles GET /api/v1/users/:id.
func (h *UserHandler) GetUser(c *gin.Context) {
	// TODO: Extract user_id from URL param and from context (current user)

	// TODO: Call h.svc.GetUser()

	// TODO: Return 200 OK with user data

	// TODO: Handle NOT_FOUND error → 404
}

// UpdateUser handles PATCH /api/v1/users/:id.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// TODO: Extract user_id from URL param and from context

	// TODO: Check authorization (can only update own profile)

	// TODO: Parse update DTO from request body

	// TODO: Validate DTO

	// TODO: Call h.svc.UpdateUser()

	// TODO: Return 200 OK with updated user data
}

// DeleteUser handles DELETE /api/v1/users/:id.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// TODO: Extract user_id from URL param and from context

	// TODO: Check authorization

	// TODO: Call h.svc.DeleteUser()

	// TODO: Return 204 No Content

	// TODO: Handle NOT_FOUND error → 404
}

// ListUsers handles GET /api/v1/users.
func (h *UserHandler) ListUsers(c *gin.Context) {
	// TODO: Extract limit and offset from query params

	// TODO: Validate pagination params

	// TODO: Call h.svc.ListUsers()

	// TODO: Return 200 OK with paginated list
}

// RefreshAccessToken handles POST /api/v1/users/refresh.
func (h *UserHandler) RefreshAccessToken(c *gin.Context) {
	// TODO: Extract refresh_token from httpOnly cookie

	// TODO: Call h.svc.RefreshAccessToken()

	// TODO: Return 200 OK with new access token

	// TODO: Handle invalid/expired token → 401
}
