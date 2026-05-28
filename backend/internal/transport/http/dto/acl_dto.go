// internal/transport/http/dto/acl_dto.go
package dto

// ACLEntryRequest — DTO для запроса обновления ACL записи.
type ACLEntryRequest struct {
	RoleID     int  `json:"role_id" validate:"required,gt=0"`
	ResourceID int  `json:"resource_id" validate:"required,gt=0"`
	CanRead    bool `json:"can_read"`
	CanWrite   bool `json:"can_write"`
	CanDelete  bool `json:"can_delete"`
}

// ACLEntryResponse — DTO для ответа с ACL записью.
type ACLEntryResponse struct {
	RoleID     int              `json:"role_id"`
	ResourceID int              `json:"resource_id"`
	Resource   ResourceResponse `json:"resource"`
	CanRead    bool             `json:"can_read"`
	CanWrite   bool             `json:"can_write"`
	CanDelete  bool             `json:"can_delete"`
}

// RoleResponse — DTO для ответа с информацией о роли.
type RoleResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ResourceResponse — DTO для ответа с информацией о ресурсе.
type ResourceResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
