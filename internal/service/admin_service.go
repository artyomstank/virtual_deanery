package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
	service_iface "github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// adminService обрабатывает все операции администратора с пользователями.
type adminService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	logger     *logger.Logger
	bcryptCost int
}

func NewAdminService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	logger *logger.Logger,
	bcryptCost int,
) service_iface.AdminService {
	return &adminService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		logger:     logger,
		bcryptCost: bcryptCost,
	}
}

// AdminCreateUser создаёт пользователя с заданной ролью и статусом активности.
func (s *adminService) AdminCreateUser(ctx context.Context, username, email, password, roleName string, isActive bool) (*entity.User, error) {
	if username == "" || email == "" || password == "" || roleName == "" {
		return nil, apperror.ErrInvalidInput
	}

	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		s.logger.Warn("role not found", map[string]interface{}{"role": roleName})
		return nil, apperror.ErrInvalidInput
	}

	// Проверка уникальности
	if _, err := s.userRepo.GetByEmail(ctx, email); err == nil {
		return nil, apperror.ErrConflict
	}
	if _, err := s.userRepo.GetByUsername(ctx, username); err == nil {
		return nil, apperror.ErrConflict
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		s.logger.Error("failed to hash password", err, nil)
		return nil, apperror.ErrInternal
	}

	now := time.Now()
	user := &entity.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPwd),
		IsActive:     isActive,
		CreatedAt:    now,
		UpdatedAt:    now,
		Role:         *role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user by admin", err, nil)
		return nil, apperror.ErrInternal
	}

	s.logger.Info("admin created user", map[string]interface{}{
		"user_id":   user.ID,
		"role":      roleName,
		"is_active": isActive,
	})
	return user, nil
}

// ApproveUser подтверждает регистрацию (делает пользователя активным).
func (s *adminService) ApproveUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	user.IsActive = true
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to approve user", err, nil)
		return apperror.ErrInternal
	}

	s.logger.Info("user approved", map[string]interface{}{"user_id": userID})
	return nil
}

// ChangeUserStatus блокирует или разблокирует пользователя.
func (s *adminService) ChangeUserStatus(ctx context.Context, userID int, isActive bool) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	user.IsActive = isActive
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to change user status", err, nil)
		return apperror.ErrInternal
	}

	s.logger.Info("user status changed", map[string]interface{}{
		"user_id":   userID,
		"is_active": isActive,
	})
	return nil
}

// ChangeUserRole изменяет роль пользователя.
func (s *adminService) ChangeUserRole(ctx context.Context, userID int, newRole string) error {
	role, err := s.roleRepo.GetByName(ctx, newRole)
	if err != nil {
		s.logger.Warn("role not found", map[string]interface{}{"role": newRole})
		return apperror.ErrInvalidInput
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	user.Role = *role
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to change role", err, nil)
		return apperror.ErrInternal
	}

	s.logger.Info("user role changed", map[string]interface{}{
		"user_id":  userID,
		"new_role": newRole,
	})
	return nil
}

// ListUsers возвращает список пользователей с фильтром.
func (s *adminService) ListUsers(ctx context.Context, filter service_iface.UserFilter) ([]entity.User, error) {
	users, err := s.userRepo.List(ctx, filter.IsActive, filter.Role)
	if err != nil {
		s.logger.Error("failed to list users", err, nil)
		return nil, apperror.ErrInternal
	}
	return users, nil
}

// GetUserByID возвращает пользователя по идентификатору.
func (s *adminService) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	return user, nil
}
