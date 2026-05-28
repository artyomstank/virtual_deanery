package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/repository"
)

// Проверки соответствия интерфейсам на этапе компиляции.
var (
	_ repository.UserRepository = (*UserRepo)(nil)
	_ repository.RoleRepository = (*RoleRepo)(nil)
	_ repository.ACLRepository  = (*ACLRepo)(nil)
)

// UserRepo — in‑memory реализация репозитория пользователей.
type UserRepo struct {
	mu      sync.RWMutex
	users   map[int]*entity.User
	counter int
}

// NewUserRepo создаёт новый экземпляр UserRepo с инициализированным хранилищем.
func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[int]*entity.User),
	}
}

// Create добавляет нового пользователя в хранилище.
// Проверяет уникальность email и username.
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверка уникальности email и username
	for _, u := range r.users {
		if u.Email == user.Email {
			return fmt.Errorf("email already exists")
		}
		if u.Username == user.Username {
			return fmt.Errorf("username already exists")
		}
	}

	r.counter++
	user.ID = r.counter

	// Сохраняем копию, чтобы извне нельзя было изменить внутреннее состояние
	clone := *user
	r.users[user.ID] = &clone

	return nil
}

// GetByID возвращает копию пользователя по идентификатору.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	clone := *u
	return &clone, nil
}

// GetByEmail возвращает копию пользователя по адресу электронной почты.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			clone := *u
			return &clone, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetByUsername возвращает копию пользователя по имени пользователя.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Username == username {
			clone := *u
			return &clone, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// Update обновляет флаги IsActive и UpdatedAt существующего пользователя.
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[user.ID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	u.IsActive = user.IsActive
	u.UpdatedAt = time.Now()
	return nil
}

// AssignRole назначает роль пользователю по его идентификатору.
// Создаёт роль с указанным именем, не проверяя её существование в RoleRepo.
func (r *UserRepo) AssignRole(ctx context.Context, userID int, roleName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	u.Role = entity.Role{Name: roleName}
	return nil
}

// RoleRepo — in‑memory реализация репозитория ролей.
type RoleRepo struct {
	mu    sync.RWMutex
	roles map[string]*entity.Role
}

// NewRoleRepo создаёт экземпляр RoleRepo с предустановленными базовыми ролями.
func NewRoleRepo() *RoleRepo {
	return &RoleRepo{
		roles: map[string]*entity.Role{
			"admin":   {ID: 1, Name: "admin"},
			"teacher": {ID: 2, Name: "teacher"},
			"student": {ID: 3, Name: "student"},
			"dean":    {ID: 4, Name: "dean"},
		},
	}
}

// GetByName возвращает копию роли по её имени.
func (r *RoleRepo) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, ok := r.roles[name]
	if !ok {
		return nil, fmt.Errorf("role not found")
	}

	clone := *role
	return &clone, nil
}

// GetAll возвращает срез со всеми ролями (копиями объектов).
func (r *RoleRepo) GetAll(ctx context.Context) ([]entity.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.Role, 0, len(r.roles))
	for _, role := range r.roles {
		result = append(result, *role)
	}
	return result, nil
}

// ACLRepo — in‑memory реализация репозитория записей ACL.
type ACLRepo struct {
	mu      sync.RWMutex
	entries []entity.ACLEntry
	roles   map[int]string // roleID → roleName, заполняется извне
}

// NewACLRepo создаёт экземпляр ACLRepo с пустым слайсом записей.
func NewACLRepo() *ACLRepo {
	return &ACLRepo{
		entries: make([]entity.ACLEntry, 0),
		roles:   make(map[int]string),
	}
}

// SetRoles устанавливает отображение идентификаторов ролей в имена.
func (r *ACLRepo) SetRoles(roles map[int]string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.roles = roles
}

// GetByRoleName возвращает все ACL‑записи, соответствующие указанному имени роли.
// Если роль не найдена в маппинге, возвращается пустой (не nil) слайс.
func (r *ACLRepo) GetByRoleName(ctx context.Context, roleName string) ([]entity.ACLEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Поиск идентификатора роли по имени
	var roleID int
	found := false
	for id, name := range r.roles {
		if name == roleName {
			roleID = id
			found = true
			break
		}
	}

	if !found {
		return []entity.ACLEntry{}, nil
	}

	result := make([]entity.ACLEntry, 0)
	for _, entry := range r.entries {
		if entry.RoleID == roleID {
			result = append(result, entry)
		}
	}
	return result, nil
}

// UpdateEntry обновляет существующую ACL‑запись по RoleID и ResourceID,
// либо добавляет новую, если запись не найдена (upsert‑поведение).
func (r *ACLRepo) UpdateEntry(ctx context.Context, entry entity.ACLEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, e := range r.entries {
		if e.RoleID == entry.RoleID && e.ResourceID == entry.ResourceID {
			r.entries[i] = entry
			return nil
		}
	}

	r.entries = append(r.entries, entry)
	return nil
}
