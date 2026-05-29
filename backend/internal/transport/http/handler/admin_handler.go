package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/artyomstank/virtual_deanery/apperror"
	postgres_repo "github.com/artyomstank/virtual_deanery/internal/repo/postgres"
	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/dto"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// AdminHandler обрабатывает все административные HTTP-запросы.
type AdminHandler struct {
	adminRepo  *postgres_repo.AdminRepo
	userSvc    service.UserService
	logger     *logger.Logger
	validator  *validator.Validate
	bcryptCost int
}

func NewAdminHandler(
	adminRepo *postgres_repo.AdminRepo,
	userSvc service.UserService,
	logger *logger.Logger,
	bcryptCost int,
) *AdminHandler {
	return &AdminHandler{
		adminRepo:  adminRepo,
		userSvc:    userSvc,
		logger:     logger,
		validator:  validator.New(),
		bcryptCost: bcryptCost,
	}
}

// =============================================================
// Stats
// =============================================================

// GetStats — GET /admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminRepo.GetStats(c.Request.Context())
	if err != nil {
		h.logger.Error("get stats error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}
	c.JSON(http.StatusOK, dto.StatsResponse{
		TotalUsers:      stats.TotalUsers,
		Students:        stats.Students,
		Teachers:        stats.Teachers,
		PendingApproval: stats.PendingApproval,
	})
}

// =============================================================
// Audit Logs
// =============================================================

// GetAuditLogs — GET /admin/audit?limit=N
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))

	logs, total, err := h.adminRepo.GetAuditLogs(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("get audit logs error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	entries := make([]dto.AuditLogEntry, len(logs))
	for i, l := range logs {
		entries[i] = dto.AuditLogEntry{
			ID:        l.ID,
			UserID:    l.UserID,
			UserName:  l.UserName,
			Action:    l.Action,
			Timestamp: l.CreatedAt,
		}
	}
	c.JSON(http.StatusOK, dto.AuditLogsResponse{Logs: entries, Total: total})
}

// =============================================================
// Users (Admin)
// =============================================================

// ListUsers — GET /admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.adminRepo.ListUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("list users error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.AdminUserItem, len(users))
	for i, u := range users {
		items[i] = dto.AdminUserItem{
			ID:        u.ID,
			Name:      u.Username,
			Email:     u.Email,
			Role:      u.Role,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		}
	}
	c.JSON(http.StatusOK, dto.AdminUsersResponse{Users: items})
}

// CreateUser — POST /admin/users
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req dto.CreateAdminUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
		return
	}

	user, err := h.userSvc.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleName: req.Role,
	})
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, dto.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
			return
		}
		h.logger.Error("create user error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	userID := user.ID
	h.adminRepo.AddAuditLog(c.Request.Context(), &userID, fmt.Sprintf("Создал пользователя %s", user.Username))

	c.JSON(http.StatusCreated, dto.AdminUserItem{
		ID:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		Status:    "active",
		CreatedAt: user.CreatedAt,
	})
}

// UpdateUser — PUT /admin/users/:id
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.UpdateAdminUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}

	if err := h.adminRepo.UpdateUserByAdmin(c.Request.Context(), id, req.Name, req.Email, req.Role, req.Status); err != nil {
		h.logger.Error("update user error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	user, err := h.adminRepo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "user not found"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Обновил пользователя %s", user.Username))

	c.JSON(http.StatusOK, dto.AdminUserItem{
		ID: user.ID, Name: user.Username, Email: user.Email,
		Role: user.Role, Status: user.Status, CreatedAt: user.CreatedAt,
	})
}

// DeleteUser — DELETE /admin/users/:id
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	user, _ := h.adminRepo.GetUserByID(c.Request.Context(), id)

	if err := h.adminRepo.DeleteUser(c.Request.Context(), id); err != nil {
		h.logger.Error("delete user error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	if user != nil {
		h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Удалил пользователя %s", user.Username))
	}

	c.JSON(http.StatusNoContent, nil)
}

// BlockUser — PATCH /admin/users/:id/block
func (h *AdminHandler) BlockUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	user, err := h.adminRepo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "user not found"})
		return
	}

	newStatus := "blocked"
	if user.Status == "blocked" {
		newStatus = "active"
	}

	if err := h.adminRepo.SetUserStatus(c.Request.Context(), id, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	user.Status = newStatus
	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Изменил статус пользователя %s -> %s", user.Username, newStatus))

	c.JSON(http.StatusOK, dto.AdminUserItem{
		ID: user.ID, Name: user.Username, Email: user.Email,
		Role: user.Role, Status: user.Status, CreatedAt: user.CreatedAt,
	})
}

// ApproveUser — PATCH /admin/users/:id/approve
func (h *AdminHandler) ApproveUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	if err := h.adminRepo.SetUserStatus(c.Request.Context(), id, "active"); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	user, err := h.adminRepo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "user not found"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Подтвердил пользователя %s", user.Username))

	c.JSON(http.StatusOK, dto.AdminUserItem{
		ID: user.ID, Name: user.Username, Email: user.Email,
		Role: user.Role, Status: user.Status, CreatedAt: user.CreatedAt,
	})
}

// ResetPassword — PATCH /admin/users/:id/reset-password
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	user, err := h.adminRepo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "user not found"})
		return
	}

	tempPassword := "Password123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), h.bcryptCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	if err := h.adminRepo.SetPasswordHash(c.Request.Context(), id, string(hash)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Сбросил пароль пользователя %s", user.Username))

	c.JSON(http.StatusOK, dto.AdminUserItem{
		ID: user.ID, Name: user.Username, Email: user.Email,
		Role: user.Role, Status: user.Status, CreatedAt: user.CreatedAt,
	})
}

// =============================================================
// Teacher Card
// =============================================================

// GetTeacherCard — GET /admin/users/:id/teacher
func (h *AdminHandler) GetTeacherCard(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	card, err := h.adminRepo.GetTeacherCardByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("get teacher card error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}
	if card == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "teacher card not found"})
		return
	}

	c.JSON(http.StatusOK, dto.TeacherCardResponse{
		ID: card.ID, UserID: card.UserID, Name: card.FullName,
		Email: card.Email, Department: card.Department,
		AcademicDegree: card.AcademicDegree, EmployeeID: card.EmployeeID, Phone: card.Phone,
	})
}

// UpdateTeacherCard — PUT /admin/users/:id/teacher
func (h *AdminHandler) UpdateTeacherCard(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.UpdateTeacherCardRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}

	card := postgres_repo.TeacherCardRow{
		UserID: userID, FullName: req.Name, Email: req.Email,
		Department: req.Department, AcademicDegree: req.AcademicDegree,
		EmployeeID: req.EmployeeID, Phone: req.Phone,
	}

	if err := h.adminRepo.UpsertTeacherCard(c.Request.Context(), userID, card); err != nil {
		h.logger.Error("upsert teacher card error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	updated, _ := h.adminRepo.GetTeacherCardByUserID(c.Request.Context(), userID)
	if updated == nil {
		c.JSON(http.StatusOK, card)
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Обновил карточку преподавателя %s", updated.FullName))

	c.JSON(http.StatusOK, dto.TeacherCardResponse{
		ID: updated.ID, UserID: updated.UserID, Name: updated.FullName,
		Email: updated.Email, Department: updated.Department,
		AcademicDegree: updated.AcademicDegree, EmployeeID: updated.EmployeeID, Phone: updated.Phone,
	})
}

// =============================================================
// Employment History
// =============================================================

// GetEmploymentHistory — GET /admin/users/:id/employment
func (h *AdminHandler) GetEmploymentHistory(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	records, err := h.adminRepo.GetEmploymentByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("get employment error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.EmploymentRecordItem, len(records))
	for i, rec := range records {
		endDate := ""
		if rec.EndDate != nil {
			endDate = *rec.EndDate
		}
		items[i] = dto.EmploymentRecordItem{
			ID: rec.ID, Position: rec.Position, EmploymentType: rec.EmploymentType,
			StartDate: rec.StartDate, EndDate: endDate, Notes: rec.Notes,
		}
	}
	c.JSON(http.StatusOK, dto.EmploymentRecordsResponse{Records: items})
}

// AddEmploymentRecord — POST /admin/users/:id/employment
func (h *AdminHandler) AddEmploymentRecord(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.AddEmploymentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
		return
	}

	id, err := h.adminRepo.AddEmployment(c.Request.Context(), userID, req.Position, req.EmploymentType, req.StartDate, req.EndDate, req.Notes)
	if err != nil {
		h.logger.Error("add employment error", err, nil)
		statusCode := http.StatusInternalServerError
		msg := "internal server error"
		if strings.Contains(err.Error(), "teacher card not found") {
			statusCode = http.StatusBadRequest
			msg = err.Error()
		}
		c.JSON(statusCode, dto.ErrorResponse{Code: statusCode, Message: msg})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, "Добавил запись трудоустройства")

	c.JSON(http.StatusCreated, dto.EmploymentRecordItem{
		ID: id, Position: req.Position, EmploymentType: req.EmploymentType,
		StartDate: req.StartDate, EndDate: req.EndDate, Notes: req.Notes,
	})
}

// GetTeacherDisciplines — GET /admin/users/:id/disciplines
func (h *AdminHandler) GetTeacherDisciplines(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	disciplines, err := h.adminRepo.GetTeacherDisciplines(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("get disciplines error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.TeacherDisciplineItem, len(disciplines))
	for i, d := range disciplines {
		items[i] = dto.TeacherDisciplineItem{
			ID: d.ID, Name: d.Name, Group: d.Group,
			Semester: strconv.Itoa(d.Semester),
		}
	}
	c.JSON(http.StatusOK, dto.TeacherDisciplinesResponse{Disciplines: items})
}

// =============================================================
// Students
// =============================================================

// ListStudents — GET /admin/students
func (h *AdminHandler) ListStudents(c *gin.Context) {
	students, err := h.adminRepo.ListStudents(c.Request.Context())
	if err != nil {
		h.logger.Error("list students error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.StudentItem, len(students))
	for i, s := range students {
		items[i] = dto.StudentItem{
			ID: s.ID, Name: s.FullName, Group: s.Group,
			Course: s.Course, Status: s.Status, RegistrationDate: s.RegistrationDate,
		}
	}
	c.JSON(http.StatusOK, dto.StudentsResponse{Students: items})
}

// UpdateStudent — PATCH /admin/students/:id
func (h *AdminHandler) UpdateStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}

	if err := h.adminRepo.UpdateStudent(c.Request.Context(), id, req.Status); err != nil {
		h.logger.Error("update student error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Обновил студента ID=%d", id))

	c.JSON(http.StatusOK, gin.H{"id": id, "status": req.Status})
}

// =============================================================
// Groups
// =============================================================

// ListGroups — GET /admin/groups
func (h *AdminHandler) ListGroups(c *gin.Context) {
	groups, err := h.adminRepo.ListGroups(c.Request.Context())
	if err != nil {
		h.logger.Error("list groups error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.GroupItem, len(groups))
	for i, g := range groups {
		items[i] = dto.GroupItem{
			ID: g.ID, Name: g.Name, Course: g.Course,
			CuratorID: g.CuratorID, CuratorName: g.CuratorName, StudentCount: g.StudentCount,
		}
	}
	c.JSON(http.StatusOK, dto.GroupsResponse{Groups: items})
}

// CreateGroup — POST /admin/groups
func (h *AdminHandler) CreateGroup(c *gin.Context) {
	var req dto.CreateGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
		return
	}

	id, err := h.adminRepo.CreateGroup(c.Request.Context(), req.Name, req.Course, req.CuratorName)
	if err != nil {
		h.logger.Error("create group error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Создал группу %s", req.Name))

	c.JSON(http.StatusCreated, dto.GroupItem{
		ID: id, Name: req.Name, Course: req.Course,
		CuratorName: req.CuratorName, StudentCount: 0,
	})
}

// UpdateGroup — PUT /admin/groups/:id
func (h *AdminHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.UpdateGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}

	if err := h.adminRepo.UpdateGroup(c.Request.Context(), id, req.Name, req.Course, req.CuratorName); err != nil {
		h.logger.Error("update group error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Обновил группу ID=%d", id))

	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name, "course": req.Course})
}

// DeleteGroup — DELETE /admin/groups/:id
func (h *AdminHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	if err := h.adminRepo.DeleteGroup(c.Request.Context(), id); err != nil {
		h.logger.Error("delete group error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Удалил группу ID=%d", id))
	c.JSON(http.StatusNoContent, nil)
}

// =============================================================
// Exam Schedule
// =============================================================

// ListExamSchedule — GET /admin/exam-schedule
func (h *AdminHandler) ListExamSchedule(c *gin.Context) {
	exams, err := h.adminRepo.ListExamSchedule(c.Request.Context())
	if err != nil {
		h.logger.Error("list exam schedule error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.ExamScheduleItem, len(exams))
	for i, e := range exams {
		items[i] = dto.ExamScheduleItem{
			ID: e.ID, Discipline: e.Discipline, Date: e.Date, Time: e.Time,
			Room: e.Room, TeacherID: e.TeacherID, TeacherName: e.TeacherName,
		}
	}
	c.JSON(http.StatusOK, dto.ExamSchedulesResponse{Exams: items})
}

// CreateExamSchedule — POST /admin/exam-schedule
func (h *AdminHandler) CreateExamSchedule(c *gin.Context) {
	var req dto.CreateExamScheduleRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
		return
	}

	id, err := h.adminRepo.CreateExamSchedule(c.Request.Context(), req.Discipline, req.Group, req.Room, req.TeacherName, req.Date, req.Time)
	if err != nil {
		h.logger.Error("create exam schedule error", err, nil)
		msg := "internal server error"
		if strings.Contains(err.Error(), "group not found") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: msg})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, "Создал расписание экзамена")

	c.JSON(http.StatusCreated, dto.ExamScheduleItem{
		ID: id, Discipline: req.Discipline, Date: req.Date, Time: req.Time,
		Room: req.Room, TeacherName: req.TeacherName,
	})
}

// UpdateExamSchedule — PUT /admin/exam-schedule/:id
func (h *AdminHandler) UpdateExamSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.UpdateExamScheduleRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}

	if err := h.adminRepo.UpdateExamSchedule(c.Request.Context(), id, req.Discipline, req.Group, req.Room, req.TeacherName, req.Date, req.Time); err != nil {
		h.logger.Error("update exam schedule error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Обновил расписание экзаменов ID=%d", id))
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// DeleteExamSchedule — DELETE /admin/exam-schedule/:id
func (h *AdminHandler) DeleteExamSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	if err := h.adminRepo.DeleteExamSchedule(c.Request.Context(), id); err != nil {
		h.logger.Error("delete exam schedule error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Удалил расписание экзамена ID=%d", id))
	c.JSON(http.StatusNoContent, nil)
}

// =============================================================
// Departments
// =============================================================

// ListDepartments — GET /admin/departments
func (h *AdminHandler) ListDepartments(c *gin.Context) {
	depts, err := h.adminRepo.ListDepartments(c.Request.Context())
	if err != nil {
		h.logger.Error("list departments error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.DepartmentItem, len(depts))
	for i, d := range depts {
		items[i] = dto.DepartmentItem{
			ID: d.ID, Name: d.Name, Office: d.Office, HeadID: d.HeadID, HeadName: d.HeadName,
		}
	}
	c.JSON(http.StatusOK, dto.DepartmentsResponse{Departments: items})
}

// CreateDepartment — POST /admin/departments
func (h *AdminHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: err.Error()})
		return
	}

	id, err := h.adminRepo.CreateDepartment(c.Request.Context(), req.Name, req.Office, req.HeadName)
	if err != nil {
		h.logger.Error("create department error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Создал кафедру %s", req.Name))

	c.JSON(http.StatusCreated, dto.DepartmentItem{
		ID: id, Name: req.Name, Office: req.Office, HeadName: req.HeadName,
	})
}

// UpdateDepartment — PUT /admin/departments/:id
func (h *AdminHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request"})
		return
	}

	if err := h.adminRepo.UpdateDepartment(c.Request.Context(), id, req.Name, req.Office, req.HeadName); err != nil {
		h.logger.Error("update department error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Обновил кафедру ID=%d", id))
	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name})
}

// DeleteDepartment — DELETE /admin/departments/:id
func (h *AdminHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid id"})
		return
	}

	if err := h.adminRepo.DeleteDepartment(c.Request.Context(), id); err != nil {
		h.logger.Error("delete department error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	callerID := callerUserID(c)
	h.adminRepo.AddAuditLog(c.Request.Context(), callerID, fmt.Sprintf("Удалил кафедру ID=%d", id))
	c.JSON(http.StatusNoContent, nil)
}

// =============================================================
// Reports
// =============================================================

// GetGradeReports — GET /reports/grades
func (h *AdminHandler) GetGradeReports(c *gin.Context) {
	reports, err := h.adminRepo.GetGradeReports(c.Request.Context())
	if err != nil {
		h.logger.Error("get grade reports error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.GradeReportItem, len(reports))
	for i, r := range reports {
		items[i] = dto.GradeReportItem{
			ID: r.ID, StudentName: r.StudentName, Discipline: r.Discipline,
			Grade: r.Grade, Date: r.Date,
		}
	}
	c.JSON(http.StatusOK, dto.GradeReportsResponse{Reports: items})
}

// GetAttendanceReports — GET /reports/attendance
func (h *AdminHandler) GetAttendanceReports(c *gin.Context) {
	reports, err := h.adminRepo.GetAttendanceReports(c.Request.Context())
	if err != nil {
		h.logger.Error("get attendance reports error", err, nil)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal server error"})
		return
	}

	items := make([]dto.AttendanceReportItem, len(reports))
	for i, r := range reports {
		items[i] = dto.AttendanceReportItem{
			ID: r.ID, StudentName: r.StudentName, Group: r.Group,
			AttendancePercent: r.AttendancePercent,
		}
	}
	c.JSON(http.StatusOK, dto.AttendanceReportsResponse{Reports: items})
}

// =============================================================
// Helpers
// =============================================================

func callerUserID(c *gin.Context) *int {
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(int); ok {
			return &id
		}
	}
	return nil
}

// Suppress unused import warnings
var _ = errors.New
var _ = fmt.Sprintf
var _ = strings.Contains
