// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  user: User;
}

// User types
export interface User {
  id: number;
  name: string;
  email: string;
  role: 'student' | 'teacher' | 'dean' | 'admin';
  status: 'active' | 'blocked' | 'pending';
  created_at: string;
}

export interface UserCreateRequest {
  name: string;
  email: string;
  password: string;
  role: 'student' | 'teacher' | 'dean' | 'admin';
}

export interface UserUpdateRequest {
  name?: string;
  email?: string;
  role?: string;
}

// Stats
export interface AdminStats {
  total_users: number;
  students: number;
  teachers: number;
  pending_approval: number;
}

// Audit log
export interface AuditLog {
  id: number;
  user_id: number | null;
  user_name: string;
  action: string;
  timestamp: string;
}

export interface AuditResponse {
  logs: AuditLog[];
  total: number;
}

// Teacher
export interface Teacher {
  id: number;
  user_id: number;
  department: string;
  academic_degree: string;
  employee_id: string;
  phone: string;
  name: string;
  email: string;
}

export interface EmploymentRecord {
  id: number;
  position: string;
  employment_type: string;
  start_date: string;
  end_date?: string;
  notes?: string;
}

// Student — backend normalises DB Russian statuses to English
export interface Student {
  id: number;
  name: string;
  group: string;
  course: number;
  status: 'active' | 'expelled' | 'on_leave' | string;
  registration_date: string;
}

// Group
export interface Group {
  id: number;
  name: string;
  course: number;
  curator_id: number | null;
  curator_name: string;
  student_count: number;
}

// Exam schedule
export interface ExamSchedule {
  id: number;
  discipline: string;
  date: string;
  time: string;
  room: string;
  teacher_id: number | null;
  teacher_name: string;
}

// Department
export interface Department {
  id: number;
  name: string;
  office: string;
  head_id: number | null;
  head_name: string;
}

// Reports
export interface GradeReport {
  id: number;
  student_name: string;
  discipline: string;
  grade: number;
  date: string;
}

export interface AttendanceReport {
  id: number;
  student_name: string;
  group: string;
  attendance_percent: number;
}

// ACL
export interface ACLPermission {
  resource: string;
  role: string;
  read: boolean;
  write: boolean;
  delete: boolean;
}

export interface ACLResponse {
  permissions: ACLPermission[];
}
