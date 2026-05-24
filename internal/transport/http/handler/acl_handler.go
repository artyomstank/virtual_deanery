// internal/transport/http/handler/acl_handler.go
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// ACLHandler обрабатывает HTTP-запросы, связанные с управлением ACL.
type ACLHandler struct {
	svc       service.ACLService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewACLHandler создаёт новый экземпляр ACLHandler.
func NewACLHandler(svc service.ACLService, logger *logger.Logger) *ACLHandler {
	return &ACLHandler{
		svc:       svc,
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
