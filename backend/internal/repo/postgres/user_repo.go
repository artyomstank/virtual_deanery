package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// UserRepo — реализация репозитория пользователей на базе PostgreSQL.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo создаёт новый экземпляр UserRepo.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create добавляет нового пользователя и назначает ему роль в одной транзакции.
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	defer tx.Rollback(ctx) // безопасно даже после коммита

	// Шаг 1: вставка пользователя
	err = tx.QueryRow(ctx,
		`INSERT INTO users(username, email, password_hash, is_active, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		user.Username, user.Email, user.PasswordHash, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if contains(pgErr.ConstraintName, "email") {
				return fmt.Errorf("email already exists")
			}
			if contains(pgErr.ConstraintName, "username") {
				return fmt.Errorf("username already exists")
			}
		}
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	// Шаг 2: получение ID роли по имени
	var roleID int
	err = tx.QueryRow(ctx,
		`SELECT id FROM roles WHERE name = $1`,
		user.Role.Name,
	).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	// Шаг 3: назначение роли пользователю
	_, err = tx.Exec(ctx,
		`INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)`,
		user.ID, roleID,
	)
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	return nil
}

// GetByID возвращает пользователя с указанным ID вместе с его ролью.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.status, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	user := &entity.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Status, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
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

// GetByEmail возвращает пользователя по email вместе с его ролью.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.status, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.email = $1`

	row := r.pool.QueryRow(ctx, query, email)
	user := &entity.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Status, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
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

// GetByUsername возвращает пользователя по username вместе с его ролью.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.status, u.is_active,
		       u.created_at, u.updated_at, r.id, r.name, r.description
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.username = $1`

	row := r.pool.QueryRow(ctx, query, username)
	user := &entity.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Status, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
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

// Update изменяет статус активности и время обновления пользователя.
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET is_active = $1, updated_at = $2 WHERE id = $3`,
		user.IsActive, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}
	return nil
}

// AssignRole заменяет текущую роль пользователя на новую.
func (r *UserRepo) AssignRole(ctx context.Context, userID int, roleName string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}
	defer tx.Rollback(ctx)

	// Удаляем все существующие роли пользователя
	_, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}

	// Находим ID роли по имени
	var roleID int
	err = tx.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1`, roleName).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("role not found")
		}
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}

	// Назначаем новую роль
	_, err = tx.Exec(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)`, userID, roleID)
	if err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("UserRepo.AssignRole: %w", err)
	}
	return nil
}

// contains проверяет, содержится ли подстрока sub в строке s.
// Используется для анализа имени ограничения.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
