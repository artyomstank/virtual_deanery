package service

import (
	"context"
	"sync"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
	service_iface "github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

type aclService struct {
	roleRepo repository.RoleRepository
	aclRepo  repository.ACLRepository
	logger   *logger.Logger
	cache    map[string][]entity.ACLEntry
	cacheMu  sync.RWMutex
}

func NewACLService(
	aclRepo repository.ACLRepository,
	roleRepo repository.RoleRepository,
	logger *logger.Logger,
) service_iface.ACLService {
	return &aclService{
		aclRepo:  aclRepo,
		roleRepo: roleRepo,
		logger:   logger,
		cache:    make(map[string][]entity.ACLEntry),
	}
}

// GetAllRoles возвращает все роли системы.
func (s *aclService) GetAllRoles(ctx context.Context) ([]entity.Role, error) {
	roles, err := s.roleRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all roles", err, nil)
		return nil, apperror.ErrInternal
	}
	return roles, nil
}

// GetACLByRole возвращает список ACL‑записей для указанной роли.
func (s *aclService) GetACLByRole(ctx context.Context, roleID int) ([]entity.ACLEntry, error) {
	entries, err := s.aclRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		s.logger.Error("failed to get ACL by role", err, nil)
		return nil, apperror.ErrInternal
	}
	return entries, nil
}

// UpdateACLEntry обновляет права для пары роль‑ресурс.
func (s *aclService) UpdateACLEntry(ctx context.Context, roleID, resourceID int, canRead, canWrite, canDelete bool) error {
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

	// Инвалидируем кэш для затронутой роли
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err == nil {
		s.invalidateCache(role.Name)
	} else {
		s.logger.Warn("could not find role for cache invalidation", map[string]interface{}{"role_id": roleID})
	}

	s.logger.Info("ACL entry updated", map[string]interface{}{
		"role_id":     roleID,
		"resource_id": resourceID,
	})
	return nil
}

// invalidateCache сбрасывает кэш ACL для одной роли или полностью (если роль пустая).
func (s *aclService) invalidateCache(roleName string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	if roleName == "" {
		s.cache = make(map[string][]entity.ACLEntry)
	} else {
		delete(s.cache, roleName)
	}
}
