// internal/service/user_service.go
package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

type userService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	aclRepo    repository.ACLRepository
	jwtClient  *jwt.Manager
	logger     *logger.Logger
	bcryptCost int
	aclCache   map[string][]entity.ACLEntry
	cacheMu    sync.RWMutex
}

// NewUserService создаёт новый экземпляр пользовательского сервиса.
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	aclRepo repository.ACLRepository,
	jwtClient *jwt.Manager,
	logger *logger.Logger,
	bcryptCost int,
) service.UserService {
	return &userService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		aclRepo:    aclRepo,
		jwtClient:  jwtClient,
		logger:     logger,
		bcryptCost: bcryptCost,
		aclCache:   make(map[string][]entity.ACLEntry),
	}
}

// Register регистрирует нового пользователя в системе.
func (s *userService) Register(ctx context.Context, input service.RegisterInput) (*entity.User, error) {
	// Валидация входных данных
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return nil, apperror.ErrInvalidInput
	}

	if len(input.Password) < 8 {
		return nil, apperror.New(apperror.ErrInvalidInput.Code, "password must be at least 8 characters", nil)
	}

	// Если роль не указана, используем "student" по умолчанию
	if input.RoleName == "" {
		input.RoleName = "student"
	}

	// Проверяем существование роли
	role, err := s.roleRepo.GetByName(ctx, input.RoleName)
	if err != nil {
		s.logger.Warn("role not found", map[string]interface{}{"role": input.RoleName})
		return nil, apperror.ErrInvalidInput
	}

	// Проверяем уникальность email
	_, err = s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, apperror.ErrConflict
	}
	if !strings.Contains(err.Error(), "not found") {
		s.logger.Error("error checking email", nil, map[string]interface{}{"error": err})
		return nil, apperror.ErrInternal
	}

	// Проверяем уникальность username
	_, err = s.userRepo.GetByUsername(ctx, input.Username)
	if err == nil {
		return nil, apperror.ErrConflict
	}
	if !strings.Contains(err.Error(), "not found") {
		s.logger.Error("error checking username", nil, map[string]interface{}{"error": err})
		return nil, apperror.ErrInternal
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		s.logger.Error("failed to hash password", err, nil)
		return nil, apperror.ErrInternal
	}

	// Создаём пользователя
	now := time.Now()
	user := &entity.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		Role:         *role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", err, nil)
		if strings.Contains(err.Error(), "already exists") {
			return nil, apperror.ErrConflict
		}
		return nil, apperror.ErrInternal
	}

	s.logger.Info("user registered successfully", map[string]interface{}{"user_id": user.ID, "role": user.Role.Name})
	return user, nil
}

// Login аутентифицирует пользователя и возвращает JWT access-токен.
func (s *userService) Login(ctx context.Context, input service.LoginInput) (*service.AuthResult, error) {
	if input.Email == "" || input.Password == "" {
		return nil, apperror.ErrUnauthorized
	}

	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Warn("login attempt with non-existent email", map[string]interface{}{"email": input.Email})
		return nil, apperror.ErrUnauthorized
	}

	// Проверяем, активен ли пользователь
	if !user.IsActive {
		s.logger.Warn("login attempt with inactive user", map[string]interface{}{"user_id": user.ID})
		return nil, apperror.ErrUnauthorized
	}

	// Сравниваем пароли
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		s.logger.Warn("invalid password attempt", map[string]interface{}{"user_id": user.ID})
		return nil, apperror.ErrUnauthorized
	}

	// Генерируем JWT access-токен
	token, err := s.jwtClient.Generate(user.ID, user.Role.Name)
	if err != nil {
		s.logger.Error("failed to generate token", err, nil)
		return nil, apperror.ErrInternal
	}

	expiresAt := time.Now().Add(time.Duration(24) * time.Hour) // TODO: Use config TTL

	s.logger.Info("user logged in successfully", map[string]interface{}{"user_id": user.ID})
	return &service.AuthResult{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// GetProfile возвращает профиль пользователя по ID.
func (s *userService) GetProfile(ctx context.Context, userID int) (*service.UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrNotFound
	}

	return &service.UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
		IsActive: user.IsActive,
	}, nil
}

// CheckPermission проверяет право пользователя на выполнение действия над ресурсом.
func (s *userService) CheckPermission(ctx context.Context, userID int, resource string, action entity.Action) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	// Получаем ACL-записи для роли (с кэшированием)
	entries, err := s.getACLEntries(ctx, user.Role.Name)
	if err != nil {
		s.logger.Error("failed to get ACL entries", err, nil)
		return apperror.ErrInternal
	}

	// Проверяем разрешение
	if !user.HasPermission(entries, resource, action) {
		s.logger.Warn("permission denied", map[string]interface{}{
			"user_id":  userID,
			"role":     user.Role.Name,
			"resource": resource,
			"action":   action,
		})
		return apperror.ErrForbidden
	}

	return nil
}

// ChangeUserStatus блокирует или разблокирует учётную запись пользователя.
func (s *userService) ChangeUserStatus(ctx context.Context, adminID int, targetUserID int, isActive bool) error {
	// Проверяем, что adminID принадлежит администратору
	admin, err := s.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return apperror.ErrNotFound
	}

	if !admin.IsAdmin() {
		s.logger.Warn("non-admin user tried to change user status", map[string]interface{}{"admin_id": adminID})
		return apperror.ErrForbidden
	}

	// Получаем целевого пользователя
	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return apperror.ErrNotFound
	}

	// Обновляем статус
	targetUser.IsActive = isActive
	targetUser.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, targetUser); err != nil {
		s.logger.Error("failed to update user status", err, nil)
		return apperror.ErrInternal
	}

	s.logger.Info("user status changed", map[string]interface{}{
		"admin_id":    adminID,
		"target_user": targetUserID,
		"is_active":   isActive,
	})
	return nil
}

// getACLEntries получает ACL-записи для роли с кэшированием.
func (s *userService) getACLEntries(ctx context.Context, roleName string) ([]entity.ACLEntry, error) {
	s.cacheMu.RLock()
	if entries, ok := s.aclCache[roleName]; ok {
		s.cacheMu.RUnlock()
		return entries, nil
	}
	s.cacheMu.RUnlock()

	// Загружаем из репозитория
	entries, err := s.aclRepo.GetByRoleName(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get ACL entries: %w", err)
	}

	// Кэшируем результат
	s.cacheMu.Lock()
	s.aclCache[roleName] = entries
	s.cacheMu.Unlock()

	return entries, nil
}

// InvalidateACLCache инвалидирует кэш ACL для роли или все, если roleName пуст.
func (s *userService) InvalidateACLCache(roleName string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	if roleName == "" {
		// Очищаем весь кэш
		s.aclCache = make(map[string][]entity.ACLEntry)
	} else {
		delete(s.aclCache, roleName)
	}
}
