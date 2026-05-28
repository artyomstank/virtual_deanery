// internal/transport/http/handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/middleware"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// UserHandler обрабатывает запросы, связанные с аутентификацией и собственным профилем.
type UserHandler struct {
	userService service.UserService
	jwtManager  *jwt.Manager
	logger      *logger.Logger
}

func NewUserHandler(userService service.UserService, jwtManager *jwt.Manager, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtManager:  jwtManager,
		logger:      log,
	}
}

// RegisterUser обрабатывает публичную регистрацию.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("register: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "некорректные данные"})
		return
	}

	authResult, err := h.userService.Register(c.Request.Context(), req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		h.logger.Errorf("register: %v", err)
		status := http.StatusInternalServerError
		if err == service.ErrUserAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	resp := dto.AuthResponse{
		AccessToken: authResult.AccessToken,
		ExpiresAt:   authResult.ExpiresAt,
		User:        userToResponse(*authResult.User),
	}
	c.JSON(http.StatusCreated, resp)
}

// LoginUser обрабатывает вход.
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("login: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: http.StatusBadRequest, Message: "некорректные данные"})
		return
	}

	authResult, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Errorf("login: %v", err)
		status := http.StatusUnauthorized
		if err == service.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		} else if err == service.ErrUserInactive {
			status = http.StatusForbidden
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	resp := dto.AuthResponse{
		AccessToken: authResult.AccessToken,
		ExpiresAt:   authResult.ExpiresAt,
		User:        userToResponse(*authResult.User),
	}
	c.JSON(http.StatusOK, resp)
}

// GetMe возвращает профиль текущего пользователя.
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		h.logger.Error("GetMe: userID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: http.StatusUnauthorized, Message: "не авторизован"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		h.logger.Error("GetMe: invalid userID type", nil, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: http.StatusInternalServerError, Message: "внутренняя ошибка"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		h.logger.Errorf("GetMe: %v", err)
		status := http.StatusNotFound
		if err != service.ErrUserNotFound {
			status = http.StatusInternalServerError
		}
		c.JSON(status, dto.ErrorResponse{Code: status, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userToResponse(*user))
}
