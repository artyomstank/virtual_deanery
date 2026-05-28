// internal/transport/http/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// AuthMiddleware валидирует JWT-токен из заголовка Authorization.
// Ожидает формат: "Authorization: Bearer <token>"
// При успешной валидации сохраняет user_id и role в контекст запроса.
func AuthMiddleware(jwtClient *jwt.Manager, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("missing authorization header", map[string]interface{}{"path": c.Request.URL.Path})
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "missing authorization header",
			})
			c.Abort()
			return
		}

		// Извлекаем токен из заголовка
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("invalid authorization header format", map[string]interface{}{"path": c.Request.URL.Path})
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Валидируем токен
		claims, err := jwtClient.Parse(token)
		if err != nil {
			logger.Warn("invalid or expired token", map[string]interface{}{"error": err})
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Сохраняем в контекст
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// LoggingMiddleware логирует HTTP-запросы и ответы.
func LoggingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("http request", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})

		c.Next()

		logger.Info("http response", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
		})
	}
}

// RecoveryMiddleware восстанавливается после паник.
func RecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", nil, map[string]interface{}{
					"panic": err,
					"path":  c.Request.URL.Path,
				})
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORSMiddleware обрабатывает CORS.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверяем, разрешено ли это происхождение
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
