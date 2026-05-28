// internal/transport/http/middleware/acl.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// ACLMiddleware создаёт middleware для проверки ACL прав на выполнение действия над ресурсом.
// Требует, чтобы user_id был установлен в контекст (обычно через AuthMiddleware).
// resource — название ресурса для проверки (например, "students", "grades")
// action — тип действия (read, write, delete)
func ACLMiddleware(userService service.UserService, resource string, action entity.Action, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем user_id из контекста (установленного AuthMiddleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			logger.Warn("user_id not found in context", map[string]interface{}{"path": c.Request.URL.Path})
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "user_id not found in context",
			})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(int)
		if !ok {
			logger.Warn("invalid user_id type in context", map[string]interface{}{"path": c.Request.URL.Path})
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid user_id",
			})
			c.Abort()
			return
		}

		// Проверяем разрешение
		err := userService.CheckPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			logger.Warn("access denied by ACL", map[string]interface{}{
				"user_id":  userID,
				"resource": resource,
				"action":   action,
				"path":     c.Request.URL.Path,
			})
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "access denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
