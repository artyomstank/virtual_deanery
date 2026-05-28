// Package service содержит интерфейсы сервисного слоя доменной области.
// Реализации находятся в internal/service.
package service

import (
	"context"
	"time"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// RegisterInput содержит входные данные для регистрации нового пользователя.
// Password — открытый пароль; хэширование bcrypt (cost >= 12) происходит в реализации сервиса.
// Если RoleName пустой — реализация подставляет "student" по умолчанию.
type RegisterInput struct {
	Username string
	Email    string
	Password string // открытый пароль, не хранится — только хэш
	RoleName string // если пусто — будет присвоена роль "student"
}

// LoginInput содержит входные данные для аутентификации пользователя.
type LoginInput struct {
	Email    string
	Password string // открытый пароль для сравнения с bcrypt-хэшем
}

// AuthResult содержит результат успешной аутентификации.
// Включает JWT access-токен и профиль пользователя с заполненной ролью.
type AuthResult struct {
	Token     string
	ExpiresAt time.Time
	User      *entity.User
}

// UserProfile содержит облегчённый профиль пользователя для HTTP-ответов.
// Не содержит PasswordHash — безопасен для передачи клиенту.
type UserProfile struct {
	ID       int
	Username string
	Email    string
	Role     string
	IsActive bool
}

// UserService определяет контракт бизнес-логики для работы с пользователями.
type UserService interface {
	// Register регистрирует нового пользователя в системе.
	// Проверяет уникальность email и username, хэширует пароль, назначает роль.
	// Ошибки: "email already exists", "username already exists".
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)

	// Login аутентифицирует пользователя и возвращает JWT access-токен.
	// Проверяет существование пользователя, IsActive и bcrypt-хэш пароля.
	// ВАЖНО: при любой ошибке (не найден / неверный пароль / заблокирован)
	// всегда возвращает одну ошибку "invalid credentials" —
	// разная реакция раскрыла бы злоумышленнику существование email в системе.
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)

	// GetProfile возвращает профиль пользователя по ID из JWT-токена.
	// Ошибки: "user not found".
	GetProfile(ctx context.Context, userID int) (*UserProfile, error)

	// CheckPermission проверяет право пользователя на выполнение действия над ресурсом.
	// Загружает ACL роли пользователя из кэша или репозитория.
	// Возвращает nil если доступ разрешён, "permission denied" если запрещён.
	CheckPermission(ctx context.Context, userID int, resource string, action entity.Action) error

	// ChangeUserStatus блокирует или разблокирует учётную запись пользователя.
	// Сначала проверяет что adminID принадлежит пользователю с ролью "admin".
	// Ошибки: "permission denied", "user not found".
	ChangeUserStatus(ctx context.Context, adminID int, targetUserID int, isActive bool) error
}
