// internal/transport/http/handler/acl_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// ACLHandler обрабатывает запросы управления правами доступа.
type ACLHandler struct {
	aclService service.ACLService
	logger     *logger.Logger
}

func NewACLHandler(aclService service.ACLService, log *logger.Logger) *ACLHandler {
	return &ACLHandler{
		aclService: aclService,
		logger:     log,
	}
}

// GetRoles возвращает список всех ролей.
func (h *ACLHandler) GetRoles(c *gin.Context) {
	roles, err := h.aclService.GetAllRoles(c.Request.Context())
	if err != nil {
		h.logger.Errorf("get roles: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: http.StatusInternalServerError, Message: "ошибка получения ролей"})
		return
	}

	resp := make([]dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, dto.RoleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// GetACLByRole возвращает список ACL-записей для указанной роли.
func (h *ACLHandler) GetACLByRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "неверный идентификатор роли"})
		return
	}

	entries, err := h.aclService.GetACLByRole(c.Request.Context(), roleID)
	if err != nil {
		h.logger.Errorf("get acl by role: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrRoleNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	resp := make([]dto.ACLEntryResponse, 0, len(entries))
	for _, e := range entries {
		resp = append(resp, dto.ACLEntryResponse{
			RoleID:     e.RoleID,
			ResourceID: e.ResourceID,
			Resource: dto.ResourceResponse{
				ID:          e.Resource.ID,
				Name:        e.Resource.Name,
				Description: e.Resource.Description,
			},
			CanRead:   e.CanRead,
			CanWrite:  e.CanWrite,
			CanDelete: e.CanDelete,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateACLEntry обновляет одну запись ACL (сбрасывает кэш при необходимости).
func (h *ACLHandler) UpdateACLEntry(c *gin.Context) {
	var req dto.ACLEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("update acl: invalid body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "некорректные данные"})
		return
	}

	if err := h.aclService.UpdateACLEntry(c.Request.Context(), req.RoleID, req.ResourceID, req.CanRead, req.CanWrite, req.CanDelete); err != nil {
		h.logger.Errorf("update acl: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrRoleNotFound || err == service.ErrResourceNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "acl updated"})
}

func userToResponse(u entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.Name,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
