// Package entity содержит чистые доменные сущности системы "Виртуальный деканат".
// Пакет не зависит ни от каких внешних библиотек — только стандартная библиотека Go.
package entity

import "time"

// Role описывает роль пользователя в системе.
// Роли статические: admin, teacher, student, dean.
// Один пользователь имеет ровно одну роль.
type Role struct {
	ID          int
	Name        string // "admin" | "teacher" | "student" | "dean"
	Description string
}

// Resource описывает защищаемый ресурс системы.
// Примеры: "students", "grades", "schedule", "teachers".
type Resource struct {
	ID          int
	Name        string
	Description string
}

// ACLEntry представляет одну запись в таблице контроля доступа.
// Связывает роль с ресурсом и набором разрешённых действий.
type ACLEntry struct {
	RoleID     int
	ResourceID int
	Resource   Resource // вложенный объект для удобной проверки имени без дополнительных запросов
	CanRead    bool
	CanWrite   bool
	CanDelete  bool
}

// Action определяет тип действия над ресурсом.
type Action string

// Предопределённые действия системы контроля доступа.
const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
)

// User представляет пользователя системы.
// Содержит учётные данные, статус активности и роль.
// PasswordHash никогда не передаётся за пределы сервисного слоя.
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string // bcrypt-хэш, cost >= 12; открытый пароль в домене не хранится
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Role         Role // вложенная структура, не строка — для работы с ACL без доп. запросов
}

// NewUser создаёт нового пользователя с указанными учётными данными.
// Устанавливает IsActive = true и временные метки в текущее время.
// Поле Role заполняется отдельно после валидации роли в сервисном слое.
func NewUser(username, email, passwordHash string) *User {
	now := time.Now()
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// HasPermission проверяет, разрешено ли пользователю выполнять указанное действие над ресурсом.
// entries — слайс ACL-записей для роли пользователя, загружается сервисным слоем.
// Если пользователь равен nil — возвращает false.
// Если запись для ресурса не найдена или флаг равен false — возвращает false.
func (u *User) HasPermission(entries []ACLEntry, resource string, action Action) bool {
	if u == nil {
		return false
	}

	for _, entry := range entries {
		if entry.Resource.Name != resource {
			continue
		}

		switch action {
		case ActionRead:
			return entry.CanRead
		case ActionWrite:
			return entry.CanWrite
		case ActionDelete:
			return entry.CanDelete
		default:
			return false
		}
	}

	return false
}

// IsAdmin возвращает true, если пользователь обладает ролью "admin".
// Если пользователь равен nil — возвращает false.
func (u *User) IsAdmin() bool {
	if u == nil {
		return false
	}

	return u.Role.Name == "admin"
}
