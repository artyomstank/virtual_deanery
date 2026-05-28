// internal/domain/service/acl_service.go
package service

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// ACLService определяет контракт для управления правами доступа.
type ACLService interface {
	// GetAllRoles возвращает список всех ролей в системе.
	GetAllRoles(ctx context.Context) ([]entity.Role, error)

	// GetACLByRole возвращает все ACL‑записи для указанной роли (по её ID).
	GetACLByRole(ctx context.Context, roleID int) ([]entity.ACLEntry, error)

	// UpdateACLEntry обновляет права для пары роль‑ресурс.
	// После успешного обновления сбрасывает кэш ACL для затронутой роли.
	UpdateACLEntry(ctx context.Context, roleID, resourceID int, canRead, canWrite, canDelete bool) error
}
