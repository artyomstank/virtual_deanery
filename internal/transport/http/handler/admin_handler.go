// internal/transport/http/handler/admin_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// AdminHandler обрабатывает административные запросы по управлению пользователями.
type AdminHandler struct {
	adminService service.AdminService
	logger       *logger.Logger
}

func NewAdminHandler(adminService service.AdminService, log *logger.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		logger:       log,
	}
}

// CreateUser создаёт пользователя с любой ролью и сразу активным.
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req dto.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("admin create user: invalid body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "некорректные данные"})
		return
	}

	// Администратор сам решает, активен ли пользователь; по умолчанию true
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user, err := h.adminService.AdminCreateUser(c.Request.Context(), req.Username, req.Email, req.Password, req.Role, isActive)
	if err != nil {
		h.logger.Errorf("admin create user: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrUserAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userToResponse(*user))
}

// ApproveUser подтверждает регистрацию (is_active = true).
func (h *AdminHandler) ApproveUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "неверный id"})
		return
	}

	if err := h.adminService.ApproveUser(c.Request.Context(), id); err != nil {
		h.logger.Errorf("approve user: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

// BlockUser изменяет статус активности (блокировка/разблокировка).
func (h *AdminHandler) BlockUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "неверный id"})
		return
	}

	var req dto.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("block user: invalid body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "некорректные данные"})
		return
	}

	if err := h.adminService.ChangeUserStatus(c.Request.Context(), id, req.IsActive); err != nil {
		h.logger.Errorf("block user: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// ChangeUserRole меняет роль пользователя.
func (h *AdminHandler) ChangeUserRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "неверный id"})
		return
	}

	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("change role: invalid body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "некорректные данные"})
		return
	}

	if err := h.adminService.ChangeUserRole(c.Request.Context(), id, req.Role); err != nil {
		h.logger.Errorf("change role: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "role updated"})
}

// ListUsers возвращает список пользователей с возможностью фильтрации (напр. ?is_active=false).
func (h *AdminHandler) ListUsers(c *gin.Context) {
	filter := service.UserFilter{
		IsActive: nil,
		Role:     c.Query("role"),
	}
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	users, err := h.adminService.ListUsers(c.Request.Context(), filter)
	if err != nil {
		h.logger.Errorf("list users: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: http.StatusInternalServerError, Message: "ошибка получения списка"})
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, userToResponse(u))
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByID возвращает пользователя по ID.
func (h *AdminHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "неверный id"})
		return
	}

	user, err := h.adminService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorf("get user by id: %v", err)
		status := http.StatusNotFound
		if err != service.ErrUserNotFound {
			status = http.StatusInternalServerError
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userToResponse(*user))
}
