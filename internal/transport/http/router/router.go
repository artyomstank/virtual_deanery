// internal/transport/http/router/router.go
package router

import (
	"github.com/artyomstank/virtual_deanery/internal/transport/http/handler"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Router sets up all HTTP routes.
type Router struct {
	engine      *gin.Engine
	userHandler *handler.UserHandler
	jwtClient   jwt.TokenClient
	logger      logger.Logger
	corsOrigins []string
}

// NewRouter creates new router.
func NewRouter(
	userHandler *handler.UserHandler,
	jwtClient jwt.TokenClient,
	logger logger.Logger,
	corsOrigins []string,
) *Router {
	return &Router{
		engine:      gin.New(),
		userHandler: userHandler,
		jwtClient:   jwtClient,
		logger:      logger,
		corsOrigins: corsOrigins,
	}
}

// Setup registers all routes and middleware.
func (r *Router) Setup() *gin.Engine {
	// TODO: Register global middleware
	// - LoggingMiddleware
	// - RecoveryMiddleware
	// - CORSMiddleware

	// TODO: Setup public routes (no auth)
	// POST   /api/v1/users/register
	// POST   /api/v1/users/login
	// POST   /api/v1/users/refresh

	// TODO: Setup protected routes group with AuthMiddleware
	// GET    /api/v1/users
	// GET    /api/v1/users/:id
	// PATCH  /api/v1/users/:id
	// DELETE /api/v1/users/:id

	// TODO: Setup health check
	// GET    /health

	return r.engine
}

// Engine returns underlying gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
