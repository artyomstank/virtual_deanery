// internal/domain/entity/user.go
package entity

import "time"

// User represents a user in the system.
type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	FullName     string    `db:"full_name"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// CreateUserInput is DTO for user creation.
type CreateUserInput struct {
	Email    string
	FullName string
	Password string
}

// UpdateUserInput is DTO for user updates.
type UpdateUserInput struct {
	FullName string
}

// UserToken represents JWT token pair.
type UserToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}
