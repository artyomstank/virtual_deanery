package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// ACLRepo — реализация репозитория ACL на базе PostgreSQL.
type ACLRepo struct {
	pool *pgxpool.Pool
}

// NewACLRepo создаёт новый экземпляр ACLRepo.
func NewACLRepo(pool *pgxpool.Pool) *ACLRepo {
	return &ACLRepo{pool: pool}
}

// GetByRoleName загружает все ACL-записи для указанной роли с JOIN на таблицы roles и resources.
func (r *ACLRepo) GetByRoleName(ctx context.Context, roleName string) ([]entity.ACLEntry, error) {
	query := `
		SELECT a.role_id, a.resource_id, a.can_read, a.can_write, a.can_delete,
		       res.id, res.name, res.description
		FROM acl_entries a
		JOIN roles ro ON ro.id = a.role_id
		JOIN resources res ON res.id = a.resource_id
		WHERE ro.name = $1
		ORDER BY a.resource_id`

	rows, err := r.pool.Query(ctx, query, roleName)
	if err != nil {
		return nil, fmt.Errorf("ACLRepo.GetByRoleName: %w", err)
	}
	defer rows.Close()

	var entries []entity.ACLEntry
	for rows.Next() {
		var entry entity.ACLEntry
		if err := rows.Scan(
			&entry.RoleID, &entry.ResourceID, &entry.CanRead, &entry.CanWrite, &entry.CanDelete,
			&entry.Resource.ID, &entry.Resource.Name, &entry.Resource.Description,
		); err != nil {
			return nil, fmt.Errorf("ACLRepo.GetByRoleName: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ACLRepo.GetByRoleName: %w", err)
	}

	// Если записей нет — возвращаем пустой слайс (не nil)
	if entries == nil {
		entries = make([]entity.ACLEntry, 0)
	}

	return entries, nil
}

// UpdateEntry обновляет флаги CanRead, CanWrite, CanDelete для пары роль+ресурс.
func (r *ACLRepo) UpdateEntry(ctx context.Context, entry entity.ACLEntry) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE acl_entries 
		 SET can_read = $1, can_write = $2, can_delete = $3
		 WHERE role_id = $4 AND resource_id = $5`,
		entry.CanRead, entry.CanWrite, entry.CanDelete,
		entry.RoleID, entry.ResourceID,
	)
	if err != nil {
		return fmt.Errorf("ACLRepo.UpdateEntry: %w", err)
	}
	return nil
}
