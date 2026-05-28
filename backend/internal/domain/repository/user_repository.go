// Package repository содержит интерфейсы репозиториев доменного слоя.
// Реализации находятся в internal/repo/postgres и internal/repo/memory.
package repository

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
)

// UserRepository определяет контракт для работы с пользователями и их ролями.
type UserRepository interface {
	// Create атомарно создаёт запись в таблице users и назначает роль в user_roles
	// в одной транзакции. Поля Username, Email, PasswordHash, Role.Name должны быть заполнены.
	// Если email уже существует — возвращает ошибку "email already exists".
	// Если username уже существует — возвращает ошибку "username already exists".
	Create(ctx context.Context, user *entity.User) error

	// GetByID возвращает пользователя с заполненным полем Role по его идентификатору.
	// Реализация выполняет JOIN на таблицы user_roles и roles.
	// Если пользователь не найден — возвращает ошибку "user not found".
	GetByID(ctx context.Context, id int) (*entity.User, error)

	// GetByEmail возвращает пользователя с заполненным полем Role по email.
	// Используется при логине и проверке уникальности email при регистрации.
	// Если пользователь не найден — возвращает ошибку "user not found".
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername возвращает пользователя с заполненным полем Role по имени пользователя.
	// Используется при регистрации для проверки уникальности username.
	// Если пользователь не найден — возвращает ошибку "user not found".
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// Update обновляет поля IsActive и UpdatedAt существующего пользователя.
	// Используется администратором при блокировке и разблокировке учётных записей.
	Update(ctx context.Context, user *entity.User) error

	// AssignRole изменяет роль пользователя, соблюдая правило "1 пользователь = 1 роль".
	// Реализация: DELETE FROM user_roles WHERE user_id = $1, затем INSERT новой роли.
	// Проверка прав (только admin может вызывать) осуществляется в сервисном слое.
	AssignRole(ctx context.Context, userID int, roleName string) error
}

// RoleRepository определяет контракт для работы со справочником ролей.
type RoleRepository interface {
	// GetByName возвращает роль по её уникальному имени.
	// Используется для валидации роли при регистрации и при смене роли администратором.
	// Если роль не найдена — возвращает ошибку "role not found".
	GetByName(ctx context.Context, name string) (*entity.Role, error)

	// GetAll возвращает полный список всех ролей системы.
	// Используется в административной панели для отображения доступных ролей.
	GetAll(ctx context.Context) ([]entity.Role, error)
}

// ACLRepository определяет контракт для работы со списками контроля доступа.
type ACLRepository interface {
	// GetByRoleName загружает все ACL-записи для указанной роли одним запросом
	// с JOIN на таблицы roles и resources — поле Resource в каждом ACLEntry заполнено.
	// Результат кэшируется в сервисном слое; метод вызывается только при промахе кэша.
	// Если записей нет — возвращает пустой слайс (не nil).
	GetByRoleName(ctx context.Context, roleName string) ([]entity.ACLEntry, error)

	// UpdateEntry обновляет флаги CanRead, CanWrite, CanDelete для пары роль+ресурс.
	// Идентификация записи по RoleID и ResourceID из переданного ACLEntry.
	// Проверка прав (только admin) осуществляется в сервисном слое.
	UpdateEntry(ctx context.Context, entry entity.ACLEntry) error
}
