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

type Router struct {
	engine       *gin.Engine
	userHandler  *handler.UserHandler
	adminHandler *handler.AdminHandler
	aclHandler   *handler.ACLHandler
	userService  service.UserService
	jwtClient    *jwt.Manager
	logger       *logger.Logger
	corsOrigins  []string
}

func NewRouter(
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	aclHandler *handler.ACLHandler,
	userService service.UserService,
	jwtClient *jwt.Manager,
	logger *logger.Logger,
	corsOrigins []string,
) *Router {
	return &Router{
		engine:       gin.New(),
		userHandler:  userHandler,
		adminHandler: adminHandler,
		aclHandler:   aclHandler,
		userService:  userService,
		jwtClient:    jwtClient,
		logger:       logger,
		corsOrigins:  corsOrigins,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.LoggingMiddleware(r.logger))
	r.engine.Use(middleware.RecoveryMiddleware(r.logger))
	r.engine.Use(middleware.CORSMiddleware(r.corsOrigins))

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	public := r.engine.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", r.userHandler.RegisterUser)
			auth.POST("/login", r.userHandler.LoginUser)
		}
	}

	protected := r.engine.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(r.jwtClient, r.logger))
	{
		users := protected.Group("/users")
		{
			users.GET("/me", r.userHandler.GetMe)
		}
	}

	admin := r.engine.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(r.jwtClient, r.logger))
	admin.Use(middleware.ACLMiddleware(r.userService, "admin", entity.ActionRead, r.logger))
	{
		adminUsers := admin.Group("/users")
		{
			adminUsers.GET("", r.adminHandler.ListUsers)
			adminUsers.POST("", r.adminHandler.CreateUser)
			adminUsers.GET("/:id", r.adminHandler.GetUserByID)
			adminUsers.POST("/:id/approve", r.adminHandler.ApproveUser)
			adminUsers.PATCH("/:id/status", r.adminHandler.BlockUser)
			adminUsers.PATCH("/:id/role", r.adminHandler.ChangeUserRole)
		}

		admin.GET("/roles", r.aclHandler.GetRoles)

		acl := admin.Group("/acl")
		{
			acl.GET("/:role", r.aclHandler.GetACLByRole)
			acl.PATCH("", r.aclHandler.UpdateACLEntry)
		}
	}

	return r.engine
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
