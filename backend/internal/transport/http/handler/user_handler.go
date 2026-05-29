// internal/transport/http/handler/user_handler.go
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/httputil"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// UserHandler обрабатывает HTTP-запросы, связанные с пользователями.
type UserHandler struct {
	svc       service.UserService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewUserHandler создаёт новый экземпляр UserHandler.
func NewUserHandler(svc service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		svc:       svc,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterUser обрабатывает POST /api/v1/users/register.
// Регистрирует нового пользователя в системе.
// Новому пользователю по умолчанию назначается роль "student" если не указана.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest

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

	// Вызываем сервис для регистрации
	user, err := h.svc.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleName: req.Role,
	})

	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, dto.ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		h.logger.Error("register error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	// Возвращаем успешный ответ (201 Created)
	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
	h.logger.Info("user registered successfully", map[string]interface{}{"user_id": user.ID})
}

// LoginUser обрабатывает POST /api/v1/users/login.
// Аутентифицирует пользователя и возвращает JWT access-токен.
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest

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
			Message: "validation failed",
		})
		return
	}

	// Вызываем сервис для логина
	authResult, err := h.svc.Login(c.Request.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, dto.ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		h.logger.Error("login error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	// Возвращаем успешный ответ (200 OK)
	response := dto.AuthResponse{
		AccessToken: authResult.Token,
		ExpiresAt:   authResult.ExpiresAt,
		User: dto.UserResponse{
			ID:        authResult.User.ID,
			Username:  authResult.User.Username,
			Email:     authResult.User.Email,
			Role:      authResult.User.Role.Name,
			IsActive:  authResult.User.IsActive,
			CreatedAt: authResult.User.CreatedAt,
			UpdatedAt: authResult.User.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
	h.logger.Info("user logged in successfully", map[string]interface{}{"user_id": authResult.User.ID})
}

// GetUser обрабатывает GET /api/v1/users/:id.
// Возвращает профиль пользователя по ID.
func (h *UserHandler) GetUser(c *gin.Context) {
	// Получаем ID из URL параметра
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid user id",
		})
		return
	}

	// Получаем профиль пользователя из сервиса
	profile, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, dto.ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		h.logger.Error("get user error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	// Возвращаем профиль пользователя
	response := dto.UserResponse{
		ID:       profile.ID,
		Username: profile.Username,
		Email:    profile.Email,
		Role:     profile.Role,
		IsActive: profile.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser обрабатывает PATCH /api/v1/users/:id.
// Обновляет профиль пользователя.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, httputil.ErrorResponse{Code: http.StatusNotImplemented, Message: "not implemented"})
}

// DeleteUser обрабатывает DELETE /api/v1/users/:id.
// Удаляет пользователя.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, httputil.ErrorResponse{Code: http.StatusNotImplemented, Message: "not implemented"})
}

// LoginPublic обрабатывает POST /auth/login.
// Возвращает ответ в формате, ожидаемом фронтендом: { access_token, user: { id, name, email, role, status } }.
func (h *UserHandler) LoginPublic(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "invalid request format"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "validation failed"})
		return
	}

	authResult, err := h.svc.Login(c.Request.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, dto.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
			return
		}
		h.logger.Error("login error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	status := authResult.User.Status
	if status == "" {
		if authResult.User.IsActive {
			status = "active"
		} else {
			status = "blocked"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": authResult.Token,
		"user": gin.H{
			"id":         authResult.User.ID,
			"name":       authResult.User.Username,
			"email":      authResult.User.Email,
			"role":       authResult.User.Role.Name,
			"status":     status,
			"created_at": authResult.User.CreatedAt,
		},
	})
}
