package dto

import "time"

// --- Auth (frontend-compatible format) ---

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthUserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthLoginResponse struct {
	AccessToken string           `json:"access_token"`
	User        AuthUserResponse `json:"user"`
}

// --- Stats ---

type StatsResponse struct {
	TotalUsers      int `json:"total_users"`
	Students        int `json:"students"`
	Teachers        int `json:"teachers"`
	PendingApproval int `json:"pending_approval"`
}

// --- Audit Log ---

type AuditLogEntry struct {
	ID        int64     `json:"id"`
	UserID    *int      `json:"user_id"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type AuditLogsResponse struct {
	Logs  []AuditLogEntry `json:"logs"`
	Total int             `json:"total"`
}

// --- Admin Users ---

type AdminUserItem struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUsersResponse struct {
	Users []AdminUserItem `json:"users"`
}

type CreateAdminUserRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=admin teacher student dean"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateAdminUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

// --- Teacher Card ---

type TeacherCardResponse struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Department     string `json:"department"`
	AcademicDegree string `json:"academic_degree"`
	EmployeeID     string `json:"employee_id"`
	Phone          string `json:"phone"`
}

type UpdateTeacherCardRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Department     string `json:"department"`
	AcademicDegree string `json:"academic_degree"`
	EmployeeID     string `json:"employee_id"`
	Phone          string `json:"phone"`
}

// --- Employment History ---

type EmploymentRecordItem struct {
	ID             int    `json:"id"`
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date,omitempty"`
	Notes          string `json:"notes"`
}

type EmploymentRecordsResponse struct {
	Records []EmploymentRecordItem `json:"records"`
}

type AddEmploymentRequest struct {
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
	StartDate      string `json:"start_date" validate:"required"`
	EndDate        string `json:"end_date"`
	Notes          string `json:"notes"`
}

// --- Teacher Disciplines ---

type TeacherDisciplineItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Semester string `json:"semester"`
}

type TeacherDisciplinesResponse struct {
	Disciplines []TeacherDisciplineItem `json:"disciplines"`
}

// --- Students ---

type StudentItem struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Group            string `json:"group"`
	Course           int    `json:"course"`
	Status           string `json:"status"`
	RegistrationDate string `json:"registration_date"`
}

type StudentsResponse struct {
	Students []StudentItem `json:"students"`
}

type UpdateStudentRequest struct {
	Status string `json:"status"`
}

// --- Groups ---

type GroupItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Course       int    `json:"course"`
	CuratorID    *int   `json:"curator_id"`
	CuratorName  string `json:"curator_name"`
	StudentCount int    `json:"student_count"`
}

type GroupsResponse struct {
	Groups []GroupItem `json:"groups"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required"`
	Course      int    `json:"course" validate:"required,min=1,max=6"`
	CuratorName string `json:"curator_name"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name"`
	Course      int    `json:"course"`
	CuratorName string `json:"curator_name"`
}

// --- Exam Schedule ---

type ExamScheduleItem struct {
	ID          int    `json:"id"`
	Discipline  string `json:"discipline"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Room        string `json:"room"`
	TeacherID   *int   `json:"teacher_id"`
	TeacherName string `json:"teacher_name"`
}

type ExamSchedulesResponse struct {
	Exams []ExamScheduleItem `json:"exams"`
}

type CreateExamScheduleRequest struct {
	Discipline  string `json:"discipline" validate:"required"`
	Group       string `json:"group" validate:"required"`
	Room        string `json:"room" validate:"required"`
	TeacherName string `json:"teacher_name"`
	Date        string `json:"date" validate:"required"`
	Time        string `json:"time" validate:"required"`
}

type UpdateExamScheduleRequest struct {
	Discipline  string `json:"discipline"`
	Group       string `json:"group"`
	Room        string `json:"room"`
	TeacherName string `json:"teacher_name"`
	Date        string `json:"date"`
	Time        string `json:"time"`
}

// --- Departments ---

type DepartmentItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Office   string `json:"office"`
	HeadID   *int   `json:"head_id"`
	HeadName string `json:"head_name"`
}

type DepartmentsResponse struct {
	Departments []DepartmentItem `json:"departments"`
}

type CreateDepartmentRequest struct {
	Name     string `json:"name" validate:"required"`
	Office   string `json:"office"`
	HeadName string `json:"head_name"`
}

type UpdateDepartmentRequest struct {
	Name     string `json:"name"`
	Office   string `json:"office"`
	HeadName string `json:"head_name"`
}

// --- Reports ---

type GradeReportItem struct {
	ID          int     `json:"id"`
	StudentName string  `json:"student_name"`
	Discipline  string  `json:"discipline"`
	Grade       float64 `json:"grade"`
	Date        string  `json:"date"`
}

type GradeReportsResponse struct {
	Reports []GradeReportItem `json:"reports"`
}

type AttendanceReportItem struct {
	ID                int     `json:"id"`
	StudentName       string  `json:"student_name"`
	Group             string  `json:"group"`
	AttendancePercent float64 `json:"attendance_percent"`
}

type AttendanceReportsResponse struct {
	Reports []AttendanceReportItem `json:"reports"`
}

// --- Full ACL Matrix ---

type ACLPermissionItem struct {
	Resource string `json:"resource"`
	Role     string `json:"role"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
	Delete   bool   `json:"delete"`
}

type FullACLResponse struct {
	Permissions []ACLPermissionItem `json:"permissions"`
}

type UpdateACLByNameRequest struct {
	Resource   string `json:"resource" validate:"required"`
	Role       string `json:"role" validate:"required"`
	Permission string `json:"permission" validate:"required,oneof=read write delete"`
	Value      bool   `json:"value"`
}
