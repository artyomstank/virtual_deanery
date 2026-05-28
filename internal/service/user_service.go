package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
	service_iface "github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// Экспортируем типы из domain/service для удобства импорта
type AuthResult = service_iface.AuthResult
type UserProfile = service_iface.UserProfile

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

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	aclRepo repository.ACLRepository,
	jwtClient *jwt.Manager,
	logger *logger.Logger,
	bcryptCost int,
) service_iface.UserService {
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

// Register регистрирует пользователя (только студент, неактивен по умолчанию).
func (s *userService) Register(ctx context.Context, username, email, password, roleName string) (*AuthResult, error) {
	// Валидация
	if username == "" || email == "" || password == "" {
		return nil, apperror.ErrInvalidInput
	}
	if len(password) < 8 {
		return nil, apperror.New(apperror.ErrInvalidInput.Code, "password must be at least 8 characters", nil)
	}

	// Для публичной регистрации доступна только роль "student", IsActive = false
	role, err := s.roleRepo.GetByName(ctx, "student")
	if err != nil {
		s.logger.Error("role 'student' not found in DB", err, nil)
		return nil, apperror.ErrInternal
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
		IsActive:     false, // студенты требуют подтверждения
		CreatedAt:    now,
		UpdatedAt:    now,
		Role:         *role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", err, nil)
		return nil, apperror.ErrInternal
	}

	// Генерация токена
	token, err := s.jwtClient.Generate(user.ID, user.Role.Name)
	if err != nil {
		s.logger.Error("failed to generate token after register", err, nil)
		return nil, apperror.ErrInternal
	}

	expiresAt := time.Now().Add(time.Duration(s.jwtClient.ExpireHours()) * time.Hour) // предполагаем метод ExpireHours()
	s.logger.Info("new student registered (inactive)", map[string]interface{}{"user_id": user.ID})
	return &AuthResult{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User:        user,
	}, nil
}

// Login аутентифицирует пользователя.
func (s *userService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	if email == "" || password == "" {
		return nil, apperror.ErrUnauthorized
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login attempt with unknown email", map[string]interface{}{"email": email})
		return nil, apperror.ErrUnauthorized
	}

	if !user.IsActive {
		s.logger.Warn("login attempt for inactive user", map[string]interface{}{"user_id": user.ID})
		return nil, apperror.New(apperror.ErrUnauthorized.Code, "account is not active", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("invalid password", map[string]interface{}{"user_id": user.ID})
		return nil, apperror.ErrUnauthorized
	}

	token, err := s.jwtClient.Generate(user.ID, user.Role.Name)
	if err != nil {
		s.logger.Error("failed to generate token", err, nil)
		return nil, apperror.ErrInternal
	}

	expiresAt := time.Now().Add(time.Duration(s.jwtClient.ExpireHours()) * time.Hour)
	return &AuthResult{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User:        user,
	}, nil
}

// GetProfile возвращает профиль текущего пользователя.
func (s *userService) GetProfile(ctx context.Context, userID int) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	return &UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
		IsActive: user.IsActive,
	}, nil
}

// GetUserByID возвращает пользователя по его ID для получения полной информации.
func (s *userService) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	return user, nil
}

// CheckPermission проверяет права пользователя на операцию с ресурсом.
func (s *userService) CheckPermission(ctx context.Context, userID int, resource string, action entity.Action) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	entries, err := s.getACLEntries(ctx, user.Role.Name)
	if err != nil {
		s.logger.Error("failed to load ACL entries for permission check", err, nil)
		return apperror.ErrInternal
	}

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

// InvalidateACLCache сбрасывает кэш ACL для одной роли или полностью (если роль пустая).
func (s *userService) InvalidateACLCache(roleName string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	if roleName == "" {
		s.aclCache = make(map[string][]entity.ACLEntry)
	} else {
		delete(s.aclCache, roleName)
	}
}

// getACLEntries возвращает ACL‑записи для роли (с кэшированием).
func (s *userService) getACLEntries(ctx context.Context, roleName string) ([]entity.ACLEntry, error) {
	s.cacheMu.RLock()
	entries, ok := s.aclCache[roleName]
	s.cacheMu.RUnlock()
	if ok {
		return entries, nil
	}

	entries, err := s.aclRepo.GetByRoleName(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("acl repo: %w", err)
	}

	s.cacheMu.Lock()
	s.aclCache[roleName] = entries
	s.cacheMu.Unlock()
	return entries, nil
}
