// Package service содержит интерфейсы сервисного слоя доменной области.
// Реализации находятся в internal/service.
package service

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// ACLService определяет контракт бизнес-логики для работы с управлением доступом.
type ACLService interface {
	// GetAllRoles возвращает список всех ролей в системе.
	GetAllRoles(ctx context.Context) ([]entity.Role, error)

	// UpdateACLEntry обновляет ACL запись для пары роль+ресурс.
	// Проверяет, что adminID принадлежит пользователю с ролью "admin".
	// Ошибки: "permission denied", "role not found", "resource not found".
	UpdateACLEntry(ctx context.Context, adminID int, roleID int, resourceID int, canRead bool, canWrite bool, canDelete bool) error

	// GetACLByRole возвращает все ACL записи для указанной роли.
	GetACLByRole(ctx context.Context, roleName string) ([]entity.ACLEntry, error)
}
