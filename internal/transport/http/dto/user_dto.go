// internal/transport/http/dto/user_dto.go
package dto

import "time"

// UserResponse is DTO for user response.
type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse is DTO for login response.
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// RefreshResponse is DTO for refresh token response.
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
