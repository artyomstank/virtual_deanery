package httputil

import (
	"errors"
	"net/http"

	"github.com/artyomstank/virtual_deanery/apperror"
	"github.com/gin-gonic/gin"
)

// SuccessResponse представляет тело успешного JSON-ответа.
type SuccessResponse struct {
	Data any `json:"data"`
}

// ErrorResponse представляет тело JSON-ответа с ошибкой.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// OK отправляет успешный ответ с HTTP-статусом 200 и переданными данными.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

// Created отправляет ответ с HTTP-статусом 201 и переданными данными.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{Data: data})
}

// NoContent отправляет ответ с HTTP-статусом 204 без тела.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error обрабатывает переданную ошибку и отправляет соответствующий JSON-ответ.
// Если ошибка является AppError, используются её код и сообщение.
// В противном случае возвращается статус 500 с текстом "internal server error".
// После отправки ответа вызывается c.Abort() для предотвращения дальнейшей обработки запроса.
func Error(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	c.Abort()
}

// BadRequest отправляет ответ с HTTP-статусом 400 и переданным сообщением об ошибке.
func BadRequest(c *gin.Context, message string) {
	err := &apperror.AppError{Code: http.StatusBadRequest, Message: message}
	Error(c, err)
}

// Unauthorized отправляет ответ с HTTP-статусом 401 и сообщением "invalid credentials".
func Unauthorized(c *gin.Context) {
	err := &apperror.AppError{Code: http.StatusUnauthorized, Message: "invalid credentials"}
	Error(c, err)
}

// Forbidden отправляет ответ с HTTP-статусом 403 и сообщением "permission denied".
func Forbidden(c *gin.Context) {
	err := &apperror.AppError{Code: http.StatusForbidden, Message: "permission denied"}
	Error(c, err)
}

// NotFound отправляет ответ с HTTP-статусом 404 и сообщением "not found".
func NotFound(c *gin.Context) {
	err := &apperror.AppError{Code: http.StatusNotFound, Message: "not found"}
	Error(c, err)
}
