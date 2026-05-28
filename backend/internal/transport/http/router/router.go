// internal/transport/http/router/router.go
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/artyomstank/virtual_deanery/internal/domain/entity"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/handler"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/middleware"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// Router устанавливает все HTTP-маршруты.
type Router struct {
	engine      *gin.Engine
	userHandler *handler.UserHandler
	aclHandler  *handler.ACLHandler
	userService service.UserService
	jwtClient   *jwt.Manager
	logger      *logger.Logger
	corsOrigins []string
}

// NewRouter создаёт новый роутер.
func NewRouter(
	userHandler *handler.UserHandler,
	aclHandler *handler.ACLHandler,
	userService service.UserService,
	jwtClient *jwt.Manager,
	logger *logger.Logger,
	corsOrigins []string,
) *Router {
	return &Router{
		engine:      gin.New(),
		userHandler: userHandler,
		aclHandler:  aclHandler,
		userService: userService,
		jwtClient:   jwtClient,
		logger:      logger,
		corsOrigins: corsOrigins,
	}
}

// Setup регистрирует все маршруты и middleware.
func (r *Router) Setup() *gin.Engine {
	// Регистрируем глобальный middleware
	r.engine.Use(middleware.LoggingMiddleware(r.logger))
	r.engine.Use(middleware.RecoveryMiddleware(r.logger))
	r.engine.Use(middleware.CORSMiddleware(r.corsOrigins))

	// Health check endpoint
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Публичные маршруты (без аутентификации)
	public := r.engine.Group("/api/v1")
	{
		users := public.Group("/users")
		{
			users.POST("/register", r.userHandler.RegisterUser)
			users.POST("/login", r.userHandler.LoginUser)
		}
	}

	// Защищённые маршруты (требуют аутентификации)
	protected := r.engine.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(r.jwtClient, r.logger))
	{
		users := protected.Group("/users")
		{
			users.GET("/:id", r.userHandler.GetUser)
			users.PATCH("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}

		// Администраторские маршруты (требуют ACL проверки на чтение ресурса "admin")
		admin := protected.Group("/admin")
		admin.Use(middleware.ACLMiddleware(r.userService, "admin", entity.ActionRead, r.logger))
		{
			// Роли
			admin.GET("/roles", r.aclHandler.GetRoles)

			// ACL управление
			acl := admin.Group("/acl")
			{
				acl.PATCH("", r.aclHandler.UpdateACLEntry)
				acl.GET("/:role", r.aclHandler.GetACLByRole)
			}
		}
	}

	return r.engine
}

// Engine возвращает базовый gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
