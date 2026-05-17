// internal/transport/http/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from Authorization header.
func AuthMiddleware(jwtClient jwt.TokenClient, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Extract token from "Authorization: Bearer <token>" header

		// TODO: Validate token using jwtClient.ValidateToken()

		// TODO: Extract claims from token

		// TODO: Put user ID in context (c.Set("user_id", claims.UserID))

		// TODO: Handle invalid/expired token → return 401

		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests and responses.
func LoggingMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Log request method, path, query params
		c.Next()
		// TODO: Log response status code, latency
	}
}

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement panic recovery
		// TODO: Log panic details
		// TODO: Return 500 Internal Server Error
		c.Next()
	}
}

// CORSMiddleware handles CORS.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get origin from request

		// TODO: Check if origin in allowedOrigins

		// TODO: Set CORS headers (Access-Control-Allow-Origin, etc.)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
