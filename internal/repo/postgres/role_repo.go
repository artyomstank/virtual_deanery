package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

type RoleRepo struct {
	pool *pgxpool.Pool
}

func NewRoleRepo(pool *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{pool: pool}
}

func (r *RoleRepo) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = $1`
	row := r.pool.QueryRow(ctx, query, name)
	role := &entity.Role{}
	err := row.Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("RoleRepo.GetByName: %w", err)
	}
	return role, nil
}

func (r *RoleRepo) GetByID(ctx context.Context, id int) (*entity.Role, error) {
	query := `SELECT id, name, description FROM roles WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	role := &entity.Role{}
	err := row.Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("RoleRepo.GetByID: %w", err)
	}
	return role, nil
}

func (r *RoleRepo) GetAll(ctx context.Context) ([]entity.Role, error) {
	query := `SELECT id, name, description FROM roles ORDER BY id`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("RoleRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, fmt.Errorf("RoleRepo.GetAll: %w", err)
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RoleRepo.GetAll: %w", err)
	}
	return roles, nil
}
