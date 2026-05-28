package service

import (
	"context"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

type aclService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	aclRepo  repository.ACLRepository
	logger   *logger.Logger
	userSvc  service.UserService // для проверки админа через интерфейс
}

// NewACLService создаёт новый экземпляр ACL сервиса.
func NewACLService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	aclRepo repository.ACLRepository,
	logger *logger.Logger,
	userSvc service.UserService,
) service.ACLService {
	return &aclService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		aclRepo:  aclRepo,
		logger:   logger,
		userSvc:  userSvc,
	}
}

// GetAllRoles возвращает список всех ролей в системе.
func (s *aclService) GetAllRoles(ctx context.Context) ([]entity.Role, error) {
	roles, err := s.roleRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all roles", err, nil)
		return nil, apperror.ErrInternal
	}

	return roles, nil
}

// UpdateACLEntry обновляет ACL запись для пары роль+ресурс.
func (s *aclService) UpdateACLEntry(ctx context.Context, adminID int, roleID int, resourceID int, canRead bool, canWrite bool, canDelete bool) error {
	// Проверяем, что adminID принадлежит администратору
	admin, err := s.userRepo.GetByID(ctx, adminID)
	if err != nil {
		s.logger.Warn("admin not found", map[string]interface{}{"admin_id": adminID})
		return apperror.ErrNotFound
	}

	if !admin.IsAdmin() {
		s.logger.Warn("non-admin user tried to update ACL", map[string]interface{}{"admin_id": adminID})
		return apperror.ErrForbidden
	}

	// Обновляем ACL запись
	entry := entity.ACLEntry{
		RoleID:     roleID,
		ResourceID: resourceID,
		CanRead:    canRead,
		CanWrite:   canWrite,
		CanDelete:  canDelete,
	}

	if err := s.aclRepo.UpdateEntry(ctx, entry); err != nil {
		s.logger.Error("failed to update ACL entry", err, nil)
		return apperror.ErrInternal
	}

	// Кэш инвалидируется при следующем обращении через getACLEntries в userService
	// или можно добавить явный метод инвалидации в интерфейс UserService

	s.logger.Info("ACL entry updated", map[string]interface{}{
		"admin_id":    adminID,
		"role_id":     roleID,
		"resource_id": resourceID,
	})

	return nil
}

// GetACLByRole возвращает все ACL записи для указанной роли.
func (s *aclService) GetACLByRole(ctx context.Context, roleName string) ([]entity.ACLEntry, error) {
	entries, err := s.aclRepo.GetByRoleName(ctx, roleName)
	if err != nil {
		s.logger.Error("failed to get ACL entries by role", err, nil)
		return nil, apperror.ErrInternal
	}

	return entries, nil
}
