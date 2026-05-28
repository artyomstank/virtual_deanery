// internal/transport/http/dto/user_dto.go
package dto

import "time"

// RegisterUserRequest — DTO для запроса регистрации пользователя.
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=admin teacher student dean"`
}

// LoginUserRequest — DTO для запроса логина.
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse — DTO для ответа с данными пользователя.
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse — DTO для ответа при успешной аутентификации.
type AuthResponse struct {
	AccessToken string        `json:"access_token"`
	ExpiresAt   time.Time     `json:"expires_at"`
	User        UserResponse  `json:"user"`
}

// ErrorResponse — DTO для ответа с ошибкой.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
