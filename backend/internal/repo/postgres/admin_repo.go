package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminRepo — репозиторий для административных операций (статистика, пользователи, академические данные).
type AdminRepo struct {
	pool *pgxpool.Pool
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{pool: pool}
}

// =============================================================
// Stats
// =============================================================

type AdminStats struct {
	TotalUsers      int
	Students        int
	Teachers        int
	PendingApproval int
}

func (r *AdminRepo) GetStats(ctx context.Context) (*AdminStats, error) {
	var s AdminStats
	err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE ro.name = 'student'),
			COUNT(*) FILTER (WHERE ro.name = 'teacher'),
			COUNT(*) FILTER (WHERE u.status = 'pending')
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles ro ON ro.id = ur.role_id`,
	).Scan(&s.TotalUsers, &s.Students, &s.Teachers, &s.PendingApproval)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetStats: %w", err)
	}
	return &s, nil
}

// =============================================================
// Audit Log
// =============================================================

type AuditLogRow struct {
	ID        int64
	UserID    *int
	UserName  string
	Action    string
	CreatedAt time.Time
}

func (r *AdminRepo) GetAuditLogs(ctx context.Context, limit int) ([]AuditLogRow, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_log`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.GetAuditLogs count: %w", err)
	}

	if limit <= 0 || limit > total {
		limit = total
	}

	rows, err := r.pool.Query(ctx, `
		SELECT al.id, al.user_id, COALESCE(u.username, 'система'), al.action, al.created_at
		FROM audit_log al
		LEFT JOIN users u ON u.id = al.user_id
		ORDER BY al.created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.GetAuditLogs: %w", err)
	}
	defer rows.Close()

	logs := make([]AuditLogRow, 0)
	for rows.Next() {
		var row AuditLogRow
		if err := rows.Scan(&row.ID, &row.UserID, &row.UserName, &row.Action, &row.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("AdminRepo.GetAuditLogs scan: %w", err)
		}
		logs = append(logs, row)
	}
	return logs, total, rows.Err()
}

func (r *AdminRepo) AddAuditLog(ctx context.Context, userID *int, action string) {
	r.pool.Exec(ctx, `INSERT INTO audit_log (user_id, action) VALUES ($1, $2)`, userID, action)
}

// =============================================================
// Users (Admin)
// =============================================================

type AdminUserRow struct {
	ID        int
	Username  string
	Email     string
	Role      string
	Status    string
	CreatedAt time.Time
}

func (r *AdminRepo) ListUsers(ctx context.Context) ([]AdminUserRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.username, u.email, ro.name, u.status, u.created_at
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles ro ON ro.id = ur.role_id
		ORDER BY u.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.ListUsers: %w", err)
	}
	defer rows.Close()

	users := make([]AdminUserRow, 0)
	for rows.Next() {
		var u AdminUserRow
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.Status, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("AdminRepo.ListUsers scan: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *AdminRepo) UpdateUserByAdmin(ctx context.Context, id int, name, email, role, status string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("AdminRepo.UpdateUserByAdmin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE users SET
			username  = CASE WHEN $1 != '' THEN $1 ELSE username END,
			email     = CASE WHEN $2 != '' THEN $2 ELSE email END,
			status    = CASE WHEN $3 != '' THEN $3 ELSE status END,
			is_active = CASE WHEN $3 != '' THEN $3 = 'active' ELSE is_active END,
			updated_at = NOW()
		WHERE id = $4`,
		name, email, status, id)
	if err != nil {
		return fmt.Errorf("AdminRepo.UpdateUserByAdmin update: %w", err)
	}

	if role != "" {
		_, err = tx.Exec(ctx, `
			UPDATE user_roles SET role_id = (SELECT id FROM roles WHERE name = $1)
			WHERE user_id = $2`, role, id)
		if err != nil {
			return fmt.Errorf("AdminRepo.UpdateUserByAdmin role: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *AdminRepo) DeleteUser(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *AdminRepo) SetUserStatus(ctx context.Context, id int, status string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET status = $1, is_active = ($1 = 'active'), updated_at = NOW()
		WHERE id = $2`, status, id)
	return err
}

func (r *AdminRepo) GetUserByID(ctx context.Context, id int) (*AdminUserRow, error) {
	var u AdminUserRow
	err := r.pool.QueryRow(ctx, `
		SELECT u.id, u.username, u.email, ro.name, u.status, u.created_at
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles ro ON ro.id = ur.role_id
		WHERE u.id = $1`, id).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.Status, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("AdminRepo.GetUserByID: %w", err)
	}
	return &u, nil
}

// =============================================================
// Teacher Card
// =============================================================

type TeacherCardRow struct {
	ID             int
	UserID         int
	FullName       string
	Email          string
	Department     string
	AcademicDegree string
	EmployeeID     string
	Phone          string
}

func (r *AdminRepo) GetTeacherCardByUserID(ctx context.Context, userID int) (*TeacherCardRow, error) {
	var row TeacherCardRow
	err := r.pool.QueryRow(ctx, `
		SELECT p.id, p.user_id,
		       p.familiya || ' ' || p.imya || COALESCE(' ' || NULLIF(p.otchestvo, ''), '') AS full_name,
		       COALESCE(p.email, ''), COALESCE(k.nazvanie, ''),
		       COALESCE(p.uchenaya_stepen, ''), COALESCE(p.tabelnyy_nomer, ''), COALESCE(p.telefon, '')
		FROM prepodavateli p
		LEFT JOIN kafedra k ON k.id = p.id_kafedry
		WHERE p.user_id = $1`, userID,
	).Scan(&row.ID, &row.UserID, &row.FullName, &row.Email, &row.Department,
		&row.AcademicDegree, &row.EmployeeID, &row.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AdminRepo.GetTeacherCardByUserID: %w", err)
	}
	return &row, nil
}

func (r *AdminRepo) UpsertTeacherCard(ctx context.Context, userID int, card TeacherCardRow) error {
	parts := strings.Fields(card.FullName)
	var familiya, imya, otchestvo string
	if len(parts) >= 1 {
		familiya = parts[0]
	}
	if len(parts) >= 2 {
		imya = parts[1]
	}
	if len(parts) >= 3 {
		otchestvo = strings.Join(parts[2:], " ")
	}

	var kafedraID *int
	if card.Department != "" {
		var kid int
		if err := r.pool.QueryRow(ctx, `SELECT id FROM kafedra WHERE nazvanie = $1`, card.Department).Scan(&kid); err == nil {
			kafedraID = &kid
		}
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO prepodavateli (user_id, familiya, imya, otchestvo, email, telefon, uchenaya_stepen, tabelnyy_nomer, id_kafedry)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''), $9)
		ON CONFLICT (user_id) DO UPDATE SET
			familiya = EXCLUDED.familiya,
			imya = EXCLUDED.imya,
			otchestvo = EXCLUDED.otchestvo,
			email = EXCLUDED.email,
			telefon = EXCLUDED.telefon,
			uchenaya_stepen = EXCLUDED.uchenaya_stepen,
			tabelnyy_nomer = EXCLUDED.tabelnyy_nomer,
			id_kafedry = EXCLUDED.id_kafedry,
			updated_at = NOW()`,
		userID, familiya, imya, otchestvo, card.Email, card.Phone,
		card.AcademicDegree, card.EmployeeID, kafedraID,
	)
	return err
}

// =============================================================
// Employment History
// =============================================================

type EmploymentRow struct {
	ID             int
	Position       string
	EmploymentType string
	StartDate      string
	EndDate        *string
	Notes          string
}

func (r *AdminRepo) GetEmploymentByUserID(ctx context.Context, userID int) ([]EmploymentRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT it.id, COALESCE(it.nazvanie_dolzhnosti, ''), COALESCE(it.tip_zanyatosti, ''),
		       it.data_nachala::text, it.data_okonchaniya::text, COALESCE(it.primechaniya, '')
		FROM istoriya_trudoustrojstva it
		JOIN prepodavateli p ON p.id = it.id_prepodavatelya
		WHERE p.user_id = $1
		ORDER BY it.data_nachala DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetEmploymentByUserID: %w", err)
	}
	defer rows.Close()

	records := make([]EmploymentRow, 0)
	for rows.Next() {
		var rec EmploymentRow
		if err := rows.Scan(&rec.ID, &rec.Position, &rec.EmploymentType, &rec.StartDate, &rec.EndDate, &rec.Notes); err != nil {
			return nil, fmt.Errorf("AdminRepo.GetEmploymentByUserID scan: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

func (r *AdminRepo) AddEmployment(ctx context.Context, userID int, position, empType, startDate, endDate, notes string) (int, error) {
	var prepID int
	if err := r.pool.QueryRow(ctx, `SELECT id FROM prepodavateli WHERE user_id = $1`, userID).Scan(&prepID); err != nil {
		return 0, fmt.Errorf("teacher card not found for user_id %d: first create the teacher card", userID)
	}

	var endPtr *string
	if endDate != "" {
		endPtr = &endDate
	}

	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO istoriya_trudoustrojstva
		       (id_prepodavatelya, nazvanie_dolzhnosti, tip_zanyatosti, data_nachala, data_okonchaniya, primechaniya)
		VALUES ($1, $2, NULLIF($3,''), $4::date, $5::date, NULLIF($6,''))
		RETURNING id`,
		prepID, position, empType, startDate, endPtr, notes,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AdminRepo.AddEmployment: %w", err)
	}
	return id, nil
}

// =============================================================
// Teacher Disciplines
// =============================================================

type DisciplineRow struct {
	ID       int
	Name     string
	Group    string
	Semester int
}

func (r *AdminRepo) GetTeacherDisciplines(ctx context.Context, userID int) ([]DisciplineRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT d.id, d.nazvanie, COALESCE(g.nazvanie, ''), COALESCE(g.semestr, 0)
		FROM discipliny d
		LEFT JOIN zanyatiya z ON z.id_discipliny = d.id
		LEFT JOIN gruppa g ON g.id = z.id_gruppy
		WHERE d.id_prepodavatelya = (SELECT id FROM prepodavateli WHERE user_id = $1)
		   OR z.id_prepodavatelya_lekcij    = (SELECT id FROM prepodavateli WHERE user_id = $1)
		   OR z.id_prepodavatelya_seminarov = (SELECT id FROM prepodavateli WHERE user_id = $1)
		   OR z.id_prepodavatelya_lab_rabot = (SELECT id FROM prepodavateli WHERE user_id = $1)
		ORDER BY d.nazvanie`, userID)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetTeacherDisciplines: %w", err)
	}
	defer rows.Close()

	disciplines := make([]DisciplineRow, 0)
	for rows.Next() {
		var d DisciplineRow
		if err := rows.Scan(&d.ID, &d.Name, &d.Group, &d.Semester); err != nil {
			return nil, fmt.Errorf("AdminRepo.GetTeacherDisciplines scan: %w", err)
		}
		disciplines = append(disciplines, d)
	}
	return disciplines, rows.Err()
}

// =============================================================
// Students
// =============================================================

type StudentRow struct {
	ID               int
	FullName         string
	Group            string
	Course           int
	Status           string
	RegistrationDate string
}

func studentStatusToEn(ru string) string {
	switch ru {
	case "активный":
		return "active"
	case "отчислен":
		return "expelled"
	case "академотпуск":
		return "on_leave"
	default:
		return ru
	}
}

func studentStatusToRu(en string) string {
	switch en {
	case "active":
		return "активный"
	case "expelled":
		return "отчислен"
	case "on_leave":
		return "академотпуск"
	default:
		return en
	}
}

func (r *AdminRepo) ListStudents(ctx context.Context) ([]StudentRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id,
		       s.familiya || ' ' || s.imya || COALESCE(' ' || NULLIF(s.otchestvo, ''), '') AS full_name,
		       COALESCE(g.nazvanie, ''),
		       COALESCE(g.kurs, 0),
		       s.status,
		       s.created_at::date::text
		FROM student s
		LEFT JOIN dose d ON d.id_studenta = s.id
		LEFT JOIN gruppa g ON g.id = d.id_gruppy
		ORDER BY s.familiya, s.imya`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.ListStudents: %w", err)
	}
	defer rows.Close()

	students := make([]StudentRow, 0)
	for rows.Next() {
		var s StudentRow
		if err := rows.Scan(&s.ID, &s.FullName, &s.Group, &s.Course, &s.Status, &s.RegistrationDate); err != nil {
			return nil, fmt.Errorf("AdminRepo.ListStudents scan: %w", err)
		}
		s.Status = studentStatusToEn(s.Status)
		students = append(students, s)
	}
	return students, rows.Err()
}

func (r *AdminRepo) UpdateStudent(ctx context.Context, id int, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE student SET status = $1, updated_at = NOW() WHERE id = $2`, studentStatusToRu(status), id)
	return err
}

// =============================================================
// Groups
// =============================================================

type GroupRow struct {
	ID           int
	Name         string
	Course       int
	CuratorID    *int
	CuratorName  string
	StudentCount int
}

func (r *AdminRepo) ListGroups(ctx context.Context) ([]GroupRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.nazvanie, g.kurs, g.id_kuratora,
		       COALESCE(p.familiya || ' ' || p.imya, ''),
		       COUNT(d.id_studenta)
		FROM gruppa g
		LEFT JOIN prepodavateli p ON p.id = g.id_kuratora
		LEFT JOIN dose d ON d.id_gruppy = g.id
		GROUP BY g.id, g.nazvanie, g.kurs, g.id_kuratora, p.familiya, p.imya
		ORDER BY g.nazvanie`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.ListGroups: %w", err)
	}
	defer rows.Close()

	groups := make([]GroupRow, 0)
	for rows.Next() {
		var g GroupRow
		if err := rows.Scan(&g.ID, &g.Name, &g.Course, &g.CuratorID, &g.CuratorName, &g.StudentCount); err != nil {
			return nil, fmt.Errorf("AdminRepo.ListGroups scan: %w", err)
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *AdminRepo) findTeacherIDByName(ctx context.Context, name string) (*int, error) {
	if name == "" {
		return nil, nil
	}
	var id int
	err := r.pool.QueryRow(ctx,
		`SELECT id FROM prepodavateli WHERE familiya || ' ' || imya ILIKE $1 LIMIT 1`,
		strings.TrimSpace(name),
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *AdminRepo) findGroupIDByName(ctx context.Context, name string) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`SELECT id FROM gruppa WHERE nazvanie = $1`, name,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("group not found: %s", name)
	}
	return id, err
}

func (r *AdminRepo) findOrCreateDisciplineByName(ctx context.Context, name string) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO discipliny (nazvanie) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id`, name,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = r.pool.QueryRow(ctx,
			`SELECT id FROM discipliny WHERE nazvanie = $1 LIMIT 1`, name,
		).Scan(&id)
	}
	return id, err
}

func (r *AdminRepo) findOrCreateRoomByName(ctx context.Context, name string) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO kabinety (nazvanie) VALUES ($1) ON CONFLICT (nazvanie, korpus) DO NOTHING RETURNING id`, name,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = r.pool.QueryRow(ctx,
			`SELECT id FROM kabinety WHERE nazvanie = $1 LIMIT 1`, name,
		).Scan(&id)
	}
	return id, err
}

func (r *AdminRepo) CreateGroup(ctx context.Context, name string, course int, curatorName string) (int, error) {
	curatorID, err := r.findTeacherIDByName(ctx, curatorName)
	if err != nil {
		return 0, err
	}
	var id int
	err = r.pool.QueryRow(ctx, `
		INSERT INTO gruppa (nazvanie, kurs, semestr, forma_obucheniya, id_kuratora)
		VALUES ($1, $2, $3, 'очная', $4)
		RETURNING id`, name, course, course*2, curatorID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AdminRepo.CreateGroup: %w", err)
	}
	return id, nil
}

func (r *AdminRepo) UpdateGroup(ctx context.Context, id int, name string, course int, curatorName string) error {
	curatorID, err := r.findTeacherIDByName(ctx, curatorName)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE gruppa SET
			nazvanie    = CASE WHEN $1 != '' THEN $1 ELSE nazvanie END,
			kurs        = CASE WHEN $2 > 0   THEN $2 ELSE kurs    END,
			id_kuratora = COALESCE($3, id_kuratora),
			updated_at  = NOW()
		WHERE id = $4`, name, course, curatorID, id)
	return err
}

func (r *AdminRepo) DeleteGroup(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM gruppa WHERE id = $1`, id)
	return err
}

// =============================================================
// Exam Schedule
// =============================================================

type ExamScheduleRow struct {
	ID          int
	Discipline  string
	Date        string
	Time        string
	Room        string
	TeacherID   *int
	TeacherName string
}

func (r *AdminRepo) ListExamSchedule(ctx context.Context) ([]ExamScheduleRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT re.id,
		       d.nazvanie,
		       re.data_ekzamena::text,
		       re.vremya_ekzamena::text,
		       k.nazvanie,
		       re.id_prepodavatelya,
		       COALESCE(p.familiya || ' ' || p.imya, '')
		FROM raspisanie_ekzamenov re
		JOIN discipliny d ON d.id = re.id_discipliny
		JOIN kabinety k ON k.id = re.id_kabineta
		LEFT JOIN prepodavateli p ON p.id = re.id_prepodavatelya
		ORDER BY re.data_ekzamena, re.vremya_ekzamena`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.ListExamSchedule: %w", err)
	}
	defer rows.Close()

	exams := make([]ExamScheduleRow, 0)
	for rows.Next() {
		var e ExamScheduleRow
		if err := rows.Scan(&e.ID, &e.Discipline, &e.Date, &e.Time, &e.Room, &e.TeacherID, &e.TeacherName); err != nil {
			return nil, fmt.Errorf("AdminRepo.ListExamSchedule scan: %w", err)
		}
		exams = append(exams, e)
	}
	return exams, rows.Err()
}

func (r *AdminRepo) CreateExamSchedule(ctx context.Context, discipline, group, room, teacherName, date, examTime string) (int, error) {
	disciplineID, err := r.findOrCreateDisciplineByName(ctx, discipline)
	if err != nil {
		return 0, fmt.Errorf("discipline lookup: %w", err)
	}
	groupID, err := r.findGroupIDByName(ctx, group)
	if err != nil {
		return 0, fmt.Errorf("group lookup: %w", err)
	}
	roomID, err := r.findOrCreateRoomByName(ctx, room)
	if err != nil {
		return 0, fmt.Errorf("room lookup: %w", err)
	}
	teacherID, err := r.findTeacherIDByName(ctx, teacherName)
	if err != nil {
		return 0, err
	}

	var id int
	err = r.pool.QueryRow(ctx, `
		INSERT INTO raspisanie_ekzamenov
		       (id_gruppy, id_discipliny, id_prepodavatelya, id_kabineta, data_ekzamena, vremya_ekzamena)
		VALUES ($1, $2, $3, $4, $5::date, $6::time)
		RETURNING id`, groupID, disciplineID, teacherID, roomID, date, examTime,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AdminRepo.CreateExamSchedule: %w", err)
	}
	return id, nil
}

func (r *AdminRepo) UpdateExamSchedule(ctx context.Context, id int, discipline, group, room, teacherName, date, examTime string) error {
	var disciplineID, groupID, roomID int
	var teacherID *int
	var err error

	if discipline != "" {
		disciplineID, err = r.findOrCreateDisciplineByName(ctx, discipline)
		if err != nil {
			return fmt.Errorf("discipline lookup: %w", err)
		}
	}
	if group != "" {
		groupID, err = r.findGroupIDByName(ctx, group)
		if err != nil {
			return fmt.Errorf("group lookup: %w", err)
		}
	}
	if room != "" {
		roomID, err = r.findOrCreateRoomByName(ctx, room)
		if err != nil {
			return fmt.Errorf("room lookup: %w", err)
		}
	}
	if teacherName != "" {
		teacherID, err = r.findTeacherIDByName(ctx, teacherName)
		if err != nil {
			return err
		}
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE raspisanie_ekzamenov SET
			id_discipliny     = CASE WHEN $1 > 0 THEN $1 ELSE id_discipliny END,
			id_gruppy         = CASE WHEN $2 > 0 THEN $2 ELSE id_gruppy END,
			id_kabineta       = CASE WHEN $3 > 0 THEN $3 ELSE id_kabineta END,
			id_prepodavatelya = COALESCE($4, id_prepodavatelya),
			data_ekzamena     = CASE WHEN $5 != '' THEN $5::date ELSE data_ekzamena END,
			vremya_ekzamena   = CASE WHEN $6 != '' THEN $6::time ELSE vremya_ekzamena END,
			updated_at        = NOW()
		WHERE id = $7`, disciplineID, groupID, roomID, teacherID, date, examTime, id)
	return err
}

func (r *AdminRepo) DeleteExamSchedule(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM raspisanie_ekzamenov WHERE id = $1`, id)
	return err
}

// =============================================================
// Departments (kafedra)
// =============================================================

type DepartmentRow struct {
	ID       int
	Name     string
	Office   string
	HeadID   *int
	HeadName string
}

func (r *AdminRepo) ListDepartments(ctx context.Context) ([]DepartmentRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT k.id, k.nazvanie, COALESCE(k.kabinet, ''), k.id_zav_kafedroy,
		       COALESCE(p.familiya || ' ' || p.imya, '')
		FROM kafedra k
		LEFT JOIN prepodavateli p ON p.id = k.id_zav_kafedroy
		ORDER BY k.nazvanie`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.ListDepartments: %w", err)
	}
	defer rows.Close()

	depts := make([]DepartmentRow, 0)
	for rows.Next() {
		var d DepartmentRow
		if err := rows.Scan(&d.ID, &d.Name, &d.Office, &d.HeadID, &d.HeadName); err != nil {
			return nil, fmt.Errorf("AdminRepo.ListDepartments scan: %w", err)
		}
		depts = append(depts, d)
	}
	return depts, rows.Err()
}

func (r *AdminRepo) CreateDepartment(ctx context.Context, name, office, headName string) (int, error) {
	headID, err := r.findTeacherIDByName(ctx, headName)
	if err != nil {
		return 0, err
	}
	var id int
	err = r.pool.QueryRow(ctx, `
		INSERT INTO kafedra (nazvanie, kabinet, id_zav_kafedroy)
		VALUES ($1, NULLIF($2,''), $3)
		RETURNING id`, name, office, headID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AdminRepo.CreateDepartment: %w", err)
	}
	return id, nil
}

func (r *AdminRepo) UpdateDepartment(ctx context.Context, id int, name, office, headName string) error {
	headID, err := r.findTeacherIDByName(ctx, headName)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE kafedra SET
			nazvanie        = CASE WHEN $1 != '' THEN $1 ELSE nazvanie END,
			kabinet         = CASE WHEN $2 != '' THEN $2 ELSE kabinet END,
			id_zav_kafedroy = COALESCE($3, id_zav_kafedroy),
			updated_at      = NOW()
		WHERE id = $4`, name, office, headID, id)
	return err
}

func (r *AdminRepo) DeleteDepartment(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM kafedra WHERE id = $1`, id)
	return err
}

// =============================================================
// Reports
// =============================================================

type GradeReportRow struct {
	ID          int
	StudentName string
	Discipline  string
	Grade       float64
	Date        string
}

func (r *AdminRepo) GetGradeReports(ctx context.Context) ([]GradeReportRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT v.id,
		       s.familiya || ' ' || s.imya AS student_name,
		       d.nazvanie,
		       COALESCE((COALESCE(v.modul_1, 0) + COALESCE(v.modul_2, 0)) / 2.0, 0),
		       v.updated_at::date::text
		FROM vedomost v
		JOIN student s ON s.id = v.id_studenta
		JOIN discipliny d ON d.id = v.id_predmeta
		ORDER BY v.updated_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetGradeReports: %w", err)
	}
	defer rows.Close()

	reports := make([]GradeReportRow, 0)
	for rows.Next() {
		var g GradeReportRow
		if err := rows.Scan(&g.ID, &g.StudentName, &g.Discipline, &g.Grade, &g.Date); err != nil {
			return nil, fmt.Errorf("AdminRepo.GetGradeReports scan: %w", err)
		}
		reports = append(reports, g)
	}
	return reports, rows.Err()
}

type AttendanceReportRow struct {
	ID                int
	StudentName       string
	Group             string
	AttendancePercent float64
}

func (r *AdminRepo) GetAttendanceReports(ctx context.Context) ([]AttendanceReportRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id,
		       s.familiya || ' ' || s.imya AS student_name,
		       COALESCE(g.nazvanie, ''),
		       COALESCE(
		           ROUND(100.0 * COUNT(pos.id) FILTER (WHERE pos.prisutstvoval) / NULLIF(COUNT(pos.id), 0), 1),
		           0
		       )
		FROM student s
		LEFT JOIN dose d ON d.id_studenta = s.id
		LEFT JOIN gruppa g ON g.id = d.id_gruppy
		LEFT JOIN poseshchaemost pos ON pos.id_studenta = s.id
		GROUP BY s.id, s.familiya, s.imya, g.nazvanie
		ORDER BY s.familiya, s.imya`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetAttendanceReports: %w", err)
	}
	defer rows.Close()

	reports := make([]AttendanceReportRow, 0)
	for rows.Next() {
		var a AttendanceReportRow
		if err := rows.Scan(&a.ID, &a.StudentName, &a.Group, &a.AttendancePercent); err != nil {
			return nil, fmt.Errorf("AdminRepo.GetAttendanceReports scan: %w", err)
		}
		reports = append(reports, a)
	}
	return reports, rows.Err()
}

// SetPasswordHash напрямую обновляет хэш пароля пользователя.
func (r *AdminRepo) SetPasswordHash(ctx context.Context, userID int, hash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`, hash, userID)
	return err
}

// =============================================================
// Full ACL Matrix
// =============================================================

type ACLPermissionRow struct {
	Resource  string
	Role      string
	CanRead   bool
	CanWrite  bool
	CanDelete bool
}

func (r *AdminRepo) GetFullACL(ctx context.Context) ([]ACLPermissionRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT res.name, ro.name, ae.can_read, ae.can_write, ae.can_delete
		FROM acl_entries ae
		JOIN roles ro ON ro.id = ae.role_id
		JOIN resources res ON res.id = ae.resource_id
		ORDER BY res.name, ro.name`)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.GetFullACL: %w", err)
	}
	defer rows.Close()

	perms := make([]ACLPermissionRow, 0)
	for rows.Next() {
		var p ACLPermissionRow
		if err := rows.Scan(&p.Resource, &p.Role, &p.CanRead, &p.CanWrite, &p.CanDelete); err != nil {
			return nil, fmt.Errorf("AdminRepo.GetFullACL scan: %w", err)
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

func (r *AdminRepo) UpdateACLByName(ctx context.Context, resource, role, permission string, value bool) (*ACLPermissionRow, error) {
	colMap := map[string]string{
		"read": "can_read", "write": "can_write", "delete": "can_delete",
	}
	col, ok := colMap[permission]
	if !ok {
		return nil, fmt.Errorf("unknown permission: %s", permission)
	}

	// col is safe — whitelisted above
	query := fmt.Sprintf(`
		UPDATE acl_entries ae
		SET %s = $1
		FROM roles ro, resources res
		WHERE ae.role_id = ro.id AND ae.resource_id = res.id
		  AND res.name = $2 AND ro.name = $3`, col)

	if _, err := r.pool.Exec(ctx, query, value, resource, role); err != nil {
		return nil, fmt.Errorf("AdminRepo.UpdateACLByName: %w", err)
	}

	var p ACLPermissionRow
	err := r.pool.QueryRow(ctx, `
		SELECT res.name, ro.name, ae.can_read, ae.can_write, ae.can_delete
		FROM acl_entries ae
		JOIN roles ro ON ro.id = ae.role_id
		JOIN resources res ON res.id = ae.resource_id
		WHERE res.name = $1 AND ro.name = $2`, resource, role,
	).Scan(&p.Resource, &p.Role, &p.CanRead, &p.CanWrite, &p.CanDelete)
	if err != nil {
		return nil, fmt.Errorf("AdminRepo.UpdateACLByName read back: %w", err)
	}
	return &p, nil
}
