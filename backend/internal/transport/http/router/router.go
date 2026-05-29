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
	engine       *gin.Engine
	userHandler  *handler.UserHandler
	aclHandler   *handler.ACLHandler
	adminHandler *handler.AdminHandler
	userService  service.UserService
	jwtClient    *jwt.Manager
	logger       *logger.Logger
	corsOrigins  []string
}

// NewRouter создаёт новый роутер.
func NewRouter(
	userHandler *handler.UserHandler,
	aclHandler *handler.ACLHandler,
	adminHandler *handler.AdminHandler,
	userService service.UserService,
	jwtClient *jwt.Manager,
	logger *logger.Logger,
	corsOrigins []string,
) *Router {
	return &Router{
		engine:       gin.New(),
		userHandler:  userHandler,
		aclHandler:   aclHandler,
		adminHandler: adminHandler,
		userService:  userService,
		jwtClient:    jwtClient,
		logger:       logger,
		corsOrigins:  corsOrigins,
	}
}

// Setup регистрирует все маршруты и middleware.
func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.LoggingMiddleware(r.logger))
	r.engine.Use(middleware.RecoveryMiddleware(r.logger))
	r.engine.Use(middleware.CORSMiddleware(r.corsOrigins))

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ==========================================================
	// Публичные маршруты (без аутентификации)
	// ==========================================================

	// Устаревший /api/v1 вариант (совместимость)
	public := r.engine.Group("/api/v1")
	{
		users := public.Group("/users")
		{
			users.POST("/register", r.userHandler.RegisterUser)
			users.POST("/login", r.userHandler.LoginUser)
		}
	}

	// Фронтенд-совместимый логин
	r.engine.POST("/auth/login", r.userHandler.LoginPublic)

	// ==========================================================
	// Защищённые маршруты (требуют JWT)
	// ==========================================================
	authMW := middleware.AuthMiddleware(r.jwtClient, r.logger)

	// --- Старые /api/v1 маршруты ---
	protected := r.engine.Group("/api/v1")
	protected.Use(authMW)
	{
		users := protected.Group("/users")
		{
			users.GET("/:id", r.userHandler.GetUser)
			users.PATCH("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.ACLMiddleware(r.userService, "admin", entity.ActionRead, r.logger))
		{
			admin.GET("/roles", r.aclHandler.GetRoles)

			acl := admin.Group("/acl")
			{
				acl.PATCH("", r.aclHandler.UpdateACLEntry)
				acl.GET("/:role", r.aclHandler.GetACLByRole)
			}
		}
	}

	// --- Новые маршруты для фронтенда (без /api/v1 префикса) ---
	fe := r.engine.Group("")
	fe.Use(authMW)
	{
		// Admin: статистика и аудит
		adm := fe.Group("/admin")
		{
			adm.GET("/stats", r.adminHandler.GetStats)
			adm.GET("/audit", r.adminHandler.GetAuditLogs)

			// ACL матрица (фронтенд-формат)
			adm.GET("/acl", r.aclHandler.GetFullACL)
			adm.PATCH("/acl", r.aclHandler.UpdateACLByName)

			// Управление пользователями
			adm.GET("/users", r.adminHandler.ListUsers)
			adm.POST("/users", r.adminHandler.CreateUser)
			adm.PUT("/users/:id", r.adminHandler.UpdateUser)
			adm.DELETE("/users/:id", r.adminHandler.DeleteUser)
			adm.PATCH("/users/:id/block", r.adminHandler.BlockUser)
			adm.PATCH("/users/:id/approve", r.adminHandler.ApproveUser)
			adm.PATCH("/users/:id/reset-password", r.adminHandler.ResetPassword)

			// Карточка преподавателя
			adm.GET("/users/:id/teacher", r.adminHandler.GetTeacherCard)
			adm.PUT("/users/:id/teacher", r.adminHandler.UpdateTeacherCard)
			adm.GET("/users/:id/employment", r.adminHandler.GetEmploymentHistory)
			adm.POST("/users/:id/employment", r.adminHandler.AddEmploymentRecord)
			adm.GET("/users/:id/disciplines", r.adminHandler.GetTeacherDisciplines)

			// Студенты
			adm.GET("/students", r.adminHandler.ListStudents)
			adm.PATCH("/students/:id", r.adminHandler.UpdateStudent)

			// Группы
			adm.GET("/groups", r.adminHandler.ListGroups)
			adm.POST("/groups", r.adminHandler.CreateGroup)
			adm.PUT("/groups/:id", r.adminHandler.UpdateGroup)
			adm.DELETE("/groups/:id", r.adminHandler.DeleteGroup)

			// Расписание экзаменов
			adm.GET("/exam-schedule", r.adminHandler.ListExamSchedule)
			adm.POST("/exam-schedule", r.adminHandler.CreateExamSchedule)
			adm.PUT("/exam-schedule/:id", r.adminHandler.UpdateExamSchedule)
			adm.DELETE("/exam-schedule/:id", r.adminHandler.DeleteExamSchedule)

			// Кафедры
			adm.GET("/departments", r.adminHandler.ListDepartments)
			adm.POST("/departments", r.adminHandler.CreateDepartment)
			adm.PUT("/departments/:id", r.adminHandler.UpdateDepartment)
			adm.DELETE("/departments/:id", r.adminHandler.DeleteDepartment)
		}

		// Отчёты
		rep := fe.Group("/reports")
		{
			rep.GET("/grades", r.adminHandler.GetGradeReports)
			rep.GET("/attendance", r.adminHandler.GetAttendanceReports)
		}
	}

	return r.engine
}

// Engine возвращает базовый gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
