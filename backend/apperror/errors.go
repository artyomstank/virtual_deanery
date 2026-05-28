// Package apperror содержит определения ошибок приложения, используемых во всех слоях.
package apperror

import (
	"errors"
	"net/http"
)

// AppError представляет ошибку с HTTP-статусом, сообщением для клиента и внутренней ошибкой.
type AppError struct {
	Code    int    // HTTP-статус код
	Message string // сообщение для клиента
	Err     error  // оригинальная внутренняя ошибка (может быть nil)
}

// Error реализует интерфейс error.
// Возвращает Message, если внутренняя ошибка отсутствует, иначе её текст.
func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap возвращает внутреннюю ошибку для поддержки errors.Is и errors.As.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New создаёт новый экземпляр AppError с заданными параметрами.
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Стандартные ошибки приложения.
var (
	// ErrNotFound — ресурс не найден.
	ErrNotFound = New(http.StatusNotFound, "not found", nil)

	// ErrUnauthorized — неверные учётные данные.
	ErrUnauthorized = New(http.StatusUnauthorized, "invalid credentials", nil)

	// ErrForbidden — недостаточно прав.
	ErrForbidden = New(http.StatusForbidden, "permission denied", nil)

	// ErrConflict — конфликт (например, дубликат записи).
	ErrConflict = New(http.StatusConflict, "already exists", nil)

	// ErrInvalidInput — некорректные входные данные.
	ErrInvalidInput = New(http.StatusBadRequest, "invalid input", nil)

	// ErrInternal — внутренняя ошибка сервера.
	ErrInternal = New(http.StatusInternalServerError, "internal server error", nil)
)

// Is проверяет, является ли ошибка err ошибкой типа *AppError с кодом, равным коду target.
// Использует errors.As для извлечения *AppError из цепочки ошибок.
func Is(err error, target *AppError) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == target.Code
	}
	return false
}
