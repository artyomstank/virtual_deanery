// internal/domain/service/user_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// Стандартные ошибки сервиса
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrRoleNotFound       = errors.New("role not found")
	ErrResourceNotFound   = errors.New("resource not found")
)

// AuthResult содержит результат успешной аутентификации или регистрации.
type AuthResult struct {
	AccessToken string
	ExpiresAt   time.Time
	User        *entity.User
}

// UserProfile – публичная информация о пользователе (безопасна для клиента).
type UserProfile struct {
	ID       int
	Username string
	Email    string
	Role     string
	IsActive bool
}

// UserService определяет контракт для работы с пользователями (регистрация, аутентификация, профиль, права).
type UserService interface {
	// Register регистрирует нового пользователя (только студент, IsActive = false).
	Register(ctx context.Context, username, email, password, roleName string) (*AuthResult, error)

	// Login аутентифицирует пользователя и возвращает JWT‑токен.
	Login(ctx context.Context, email, password string) (*AuthResult, error)

	// GetProfile возвращает профиль текущего пользователя по его ID.
	GetProfile(ctx context.Context, userID int) (*UserProfile, error)

	// GetUserByID возвращает пользователя по его ID (для получения полной информации).
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)

	// CheckPermission проверяет право пользователя на действие над ресурсом.
	CheckPermission(ctx context.Context, userID int, resource string, action entity.Action) error

	// InvalidateACLCache сбрасывает закэшированные ACL‑записи для роли.
	// Если roleName == "" — очищается весь кэш.
	InvalidateACLCache(roleName string)
}
