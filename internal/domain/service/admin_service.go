// internal/domain/service/admin_service.go
package service

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// UserFilter – параметры фильтрации списка пользователей.
type UserFilter struct {
	IsActive *bool
	Role     string
}

// AdminService определяет контракт для административных операций с пользователями.
type AdminService interface {
	// AdminCreateUser создаёт пользователя с произвольной ролью и статусом.
	AdminCreateUser(ctx context.Context, username, email, password, roleName string, isActive bool) (*entity.User, error)

	// ApproveUser активирует учётную запись (подтверждение регистрации).
	ApproveUser(ctx context.Context, userID int) error

	// ChangeUserStatus блокирует или разблокирует пользователя.
	ChangeUserStatus(ctx context.Context, userID int, isActive bool) error

	// ChangeUserRole изменяет роль пользователя.
	ChangeUserRole(ctx context.Context, userID int, newRole string) error

	// ListUsers возвращает список пользователей с учётом фильтра.
	ListUsers(ctx context.Context, filter UserFilter) ([]entity.User, error)

	// GetUserByID возвращает пользователя по идентификатору.
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}
