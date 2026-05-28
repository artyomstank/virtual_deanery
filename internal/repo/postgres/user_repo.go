package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// UserRepo — реализация репозитория пользователей на базе PostgreSQL.
type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create добавляет нового пользователя и назначает ему роль в одной транзакции.
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO users(username, email, password_hash, is_active, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		user.Username, user.Email, user.PasswordHash, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "email") {
				return fmt.Errorf("email already exists")
			}
			if strings.Contains(pgErr.ConstraintName, "username") {
				return fmt.Errorf("username already exists")
			}
		}
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	var roleID int
	err = tx.QueryRow(ctx,
		`SELECT id FROM roles WHERE name = $1`, user.Role.Name,
	).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)`,
		user.ID, roleID,
	)
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	return tx.Commit(ctx)
}

// GetByID возвращает пользователя с ролью по ID.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	user := &entity.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		&user.Role.ID, &user.Role.Name, &user.Role.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("UserRepo.GetByID: %w", err)
	}
	return user, nil
}

// GetByEmail возвращает пользователя по email с ролью.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.email = $1`

	row := r.pool.QueryRow(ctx, query, email)
	user := &entity.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		&user.Role.ID, &user.Role.Name, &user.Role.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("UserRepo.GetByEmail: %w", err)
	}
	return user, nil
}

// GetByUsername возвращает пользователя по username с ролью.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.username = $1`

	row := r.pool.QueryRow(ctx, query, username)
	user := &entity.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		&user.Role.ID, &user.Role.Name, &user.Role.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("UserRepo.GetByUsername: %w", err)
	}
	return user, nil
}

// List возвращает список пользователей с учётом фильтра.
func (r *UserRepo) List(ctx context.Context, isActive *bool, role string) ([]entity.User, error) {
	baseQuery := `
		SELECT u.id, u.username, u.email, u.password_hash, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE 1=1`

	var conditions []string
	var args []interface{}
	argIdx := 1

	if isActive != nil {
		conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", argIdx))
		args = append(args, *isActive)
		argIdx++
	}
	if role != "" {
		conditions = append(conditions, fmt.Sprintf("r.name = $%d", argIdx))
		args = append(args, role)
		argIdx++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}
	baseQuery += " ORDER BY u.id"

	rows, err := r.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.List: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
			&user.Role.ID, &user.Role.Name, &user.Role.Description,
		); err != nil {
			return nil, fmt.Errorf("UserRepo.List: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("UserRepo.List: %w", err)
	}
	return users, nil
}

// Update полностью обновляет данные пользователя, включая роль (в одной транзакции).
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}
	defer tx.Rollback(ctx)

	// Обновляем основные поля
	_, err = tx.Exec(ctx,
		`UPDATE users SET username=$1, email=$2, password_hash=$3, is_active=$4, updated_at=$5 WHERE id=$6`,
		user.Username, user.Email, user.PasswordHash, user.IsActive, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}

	// Обновляем роль: удаляем старую и вставляем новую
	_, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, user.ID)
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}

	var roleID int
	err = tx.QueryRow(ctx,
		`SELECT id FROM roles WHERE name = $1`, user.Role.Name,
	).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("role not found")
		}
		return fmt.Errorf("UserRepo.Update: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)`,
		user.ID, roleID,
	)
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}

	return tx.Commit(ctx)
}

// AssignRole изменяет роль пользователя, соблюдая правило "1 пользователь = 1 роль".
// Удаляет старую роль и назначает новую.
func (r *UserRepo) AssignRole(ctx context.Context, userID int, roleName string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}
	defer tx.Rollback(ctx)

	// Удаляем старую роль
	_, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}

	// Получаем ID новой роли
	var roleID int
	err = tx.QueryRow(ctx,
		`SELECT id FROM roles WHERE name = $1`, roleName,
	).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("role not found")
		}
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}

	// Назначаем новую роль
	_, err = tx.Exec(ctx,
		`INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)`,
		userID, roleID,
	)
	if err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}

	return tx.Commit(ctx)
}
