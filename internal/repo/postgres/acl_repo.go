package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

type ACLRepo struct {
	pool *pgxpool.Pool
}

func NewACLRepo(pool *pgxpool.Pool) *ACLRepo {
	return &ACLRepo{pool: pool}
}

// GetByRoleName загружает ACL-записи для роли по её имени.
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
	if entries == nil {
		entries = make([]entity.ACLEntry, 0)
	}
	return entries, nil
}

// GetByRoleID загружает ACL-записи для роли по её идентификатору.
func (r *ACLRepo) GetByRoleID(ctx context.Context, roleID int) ([]entity.ACLEntry, error) {
	query := `
		SELECT a.role_id, a.resource_id, a.can_read, a.can_write, a.can_delete,
		       res.id, res.name, res.description
		FROM acl_entries a
		JOIN resources res ON res.id = a.resource_id
		WHERE a.role_id = $1
		ORDER BY a.resource_id`

	rows, err := r.pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("ACLRepo.GetByRoleID: %w", err)
	}
	defer rows.Close()

	var entries []entity.ACLEntry
	for rows.Next() {
		var entry entity.ACLEntry
		if err := rows.Scan(
			&entry.RoleID, &entry.ResourceID, &entry.CanRead, &entry.CanWrite, &entry.CanDelete,
			&entry.Resource.ID, &entry.Resource.Name, &entry.Resource.Description,
		); err != nil {
			return nil, fmt.Errorf("ACLRepo.GetByRoleID: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ACLRepo.GetByRoleID: %w", err)
	}
	if entries == nil {
		entries = make([]entity.ACLEntry, 0)
	}
	return entries, nil
}

// UpdateEntry обновляет права доступа для заданной пары роль‑ресурс.
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
