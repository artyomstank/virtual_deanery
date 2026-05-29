// internal/transport/http/handler/acl_handler.go
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/artyomstank/virtual_deanery/apperror"
	postgres_repo "github.com/artyomstank/virtual_deanery/internal/repo/postgres"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// ACLHandler обрабатывает HTTP-запросы, связанные с управлением ACL.
type ACLHandler struct {
	svc       service.ACLService
	adminRepo *postgres_repo.AdminRepo
	logger    *logger.Logger
	validator *validator.Validate
}

// NewACLHandler создаёт новый экземпляр ACLHandler.
func NewACLHandler(svc service.ACLService, adminRepo *postgres_repo.AdminRepo, logger *logger.Logger) *ACLHandler {
	return &ACLHandler{
		svc:       svc,
		adminRepo: adminRepo,
		logger:    logger,
		validator: validator.New(),
	}
}

// GetRoles обрабатывает GET /api/v1/admin/roles.
// Возвращает список всех ролей в системе.
func (h *ACLHandler) GetRoles(c *gin.Context) {
	roles, err := h.svc.GetAllRoles(c.Request.Context())
	if err != nil {
		h.logger.Error("get roles error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	responses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateACLEntry обрабатывает PATCH /api/v1/admin/acl.
// Обновляет запись ACL для пары роль+ресурс.
func (h *ACLHandler) UpdateACLEntry(c *gin.Context) {
	var req dto.ACLEntryRequest

	// Парсим JSON из тела запроса
	if err := c.BindJSON(&req); err != nil {
		h.logger.Warn("invalid request format", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request format",
		})
		return
	}

	// Валидируем поля запроса
	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("validation failed", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "validation failed: " + err.Error(),
		})
		return
	}

	// Получаем user_id из контекста (установленного AuthMiddleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "user_id not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid user_id",
		})
		return
	}

	// Вызываем сервис для обновления ACL
	err := h.svc.UpdateACLEntry(c.Request.Context(), userID, req.RoleID, req.ResourceID, req.CanRead, req.CanWrite, req.CanDelete)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, dto.ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		h.logger.Error("update acl error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ErrorResponse{
		Code:    http.StatusOK,
		Message: "ACL entry updated successfully",
	})
	h.logger.Info("ACL entry updated", map[string]interface{}{"user_id": userID, "role_id": req.RoleID, "resource_id": req.ResourceID})
}

// GetFullACL обрабатывает GET /admin/acl.
// Возвращает полную матрицу прав для всех ролей и ресурсов.
func (h *ACLHandler) GetFullACL(c *gin.Context) {
	perms, err := h.adminRepo.GetFullACL(c.Request.Context())
	if err != nil {
		h.logger.Error("get full acl error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.ACLPermissionItem, len(perms))
	for i, p := range perms {
		items[i] = dto.ACLPermissionItem{
			Resource: p.Resource, Role: p.Role,
			Read: p.CanRead, Write: p.CanWrite, Delete: p.CanDelete,
		}
	}
	c.JSON(http.StatusOK, dto.FullACLResponse{Permissions: items})
}

// UpdateACLByName обрабатывает PATCH /admin/acl.
// Принимает { resource, role, permission, value } — формат фронтенда.
func (h *ACLHandler) UpdateACLByName(c *gin.Context) {
	var req dto.UpdateACLByNameRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
		return
	}

	updated, err := h.adminRepo.UpdateACLByName(c.Request.Context(), req.Resource, req.Role, req.Permission, req.Value)
	if err != nil {
		h.logger.Error("update acl by name error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.ACLPermissionItem{
		Resource: updated.Resource, Role: updated.Role,
		Read: updated.CanRead, Write: updated.CanWrite, Delete: updated.CanDelete,
	})
}

// GetACLByRole обрабатывает GET /api/v1/admin/acl/:role.
// Возвращает все ACL записи для указанной роли.
func (h *ACLHandler) GetACLByRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "role parameter is required",
		})
		return
	}

	entries, err := h.svc.GetACLByRole(c.Request.Context(), role)
	if err != nil {
		h.logger.Error("get acl by role error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	responses := make([]dto.ACLEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = dto.ACLEntryResponse{
			RoleID:     entry.RoleID,
			ResourceID: entry.ResourceID,
			Resource: dto.ResourceResponse{
				ID:          entry.Resource.ID,
				Name:        entry.Resource.Name,
				Description: entry.Resource.Description,
			},
			CanRead:   entry.CanRead,
			CanWrite:  entry.CanWrite,
			CanDelete: entry.CanDelete,
		}
	}

	c.JSON(http.StatusOK, responses)
}
