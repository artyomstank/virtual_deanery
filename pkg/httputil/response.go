// pkg/httputil/response.go
package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is generic JSON response.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PaginationResponse wraps paginated data.
type PaginationResponse struct {
	Total  int64       `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Items  interface{} `json:"items"`
}

// JSON sends JSON response with status code.
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	})
}

// OK sends 200 OK with data.
func OK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, data)
}

// Created sends 201 Created with data.
func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, data)
}

// NoContent sends 204 No Content.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends error response.
// TODO: Map apperror codes to HTTP status codes
// NOT_FOUND → 404
// CONFLICT → 409
// INVALID_CREDENTIALS → 401
// FORBIDDEN → 403
// BAD_REQUEST → 400
// INTERNAL_SERVER → 500
func Error(c *gin.Context, err error) {
	// TODO: Check if err is apperror, extract code

	// TODO: Determine HTTP status code

	// TODO: Send error response

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER",
			Message: "internal server error",
		},
	})
}

// Paginated sends paginated response.
func Paginated(c *gin.Context, items interface{}, total int64, limit, offset int) {
	JSON(c, http.StatusOK, PaginationResponse{
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Items:  items,
	})
}
