import {
  ACLPermission,
  AttendanceReport,
  AuditLog,
  Department,
  EmploymentRecord,
  ExamSchedule,
  GradeReport,
  Group,
  Student,
  Teacher,
  User,
} from "../types/api";

const API_BASE = "http://localhost:8088";

const mockUser: User = {
  id: 1,
  name: "Анна Смирнова",
  email: "admin@dekanat.local",
  role: "admin",
  status: "active",
  created_at: "2026-01-10T09:00:00.000Z",
};

const mockState = {
  users: [
    mockUser,
    { id: 2, name: "Мария Иванова", email: "ivanova@dekanat.local", role: "teacher", status: "active", created_at: "2026-01-12T11:20:00.000Z" },
    { id: 3, name: "Дмитрий Соколов", email: "sokolov@dekanat.local", role: "teacher", status: "pending", created_at: "2026-02-04T10:15:00.000Z" },
    { id: 4, name: "Иван Петров", email: "petrov@student.local", role: "student", status: "active", created_at: "2026-01-18T08:30:00.000Z" },
    { id: 5, name: "Елена Кузнецова", email: "kuznetsova@student.local", role: "student", status: "blocked", created_at: "2026-02-01T13:45:00.000Z" },
    { id: 6, name: "Олег Васильев", email: "dean@dekanat.local", role: "dean", status: "active", created_at: "2026-01-09T14:00:00.000Z" },
  ] as User[],
  audit: [
    { id: 1, user_id: 1, user_name: "Анна Смирнова", action: "Создала пользователя", timestamp: "2026-05-29T08:20:00.000Z" },
    { id: 2, user_id: 1, user_name: "Анна Смирнова", action: "Обновила расписание экзаменов", timestamp: "2026-05-29T08:05:00.000Z" },
    { id: 3, user_id: 6, user_name: "Олег Васильев", action: "Сформировал отчёт посещаемости", timestamp: "2026-05-28T16:40:00.000Z" },
    { id: 4, user_id: 1, user_name: "Анна Смирнова", action: "Подтвердила преподавателя", timestamp: "2026-05-28T14:15:00.000Z" },
  ] as AuditLog[],
  teachers: {
    "2": { id: 1, user_id: 2, name: "Мария Иванова", email: "ivanova@dekanat.local", department: "Прикладная информатика", academic_degree: "к.т.н.", employee_id: "T-1024", phone: "+7 900 123-45-67" },
    "3": { id: 2, user_id: 3, name: "Дмитрий Соколов", email: "sokolov@dekanat.local", department: "Математический анализ", academic_degree: "к.ф.-м.н.", employee_id: "T-1188", phone: "+7 900 765-43-21" },
  } as Record<string, Teacher>,
  employment: {
    "2": [
      { id: 1, position: "Доцент", employment_type: "Основная", start_date: "2023-09-01", notes: "Кафедра ПИ" },
      { id: 2, position: "Старший преподаватель", employment_type: "Основная", start_date: "2021-09-01", end_date: "2023-08-31", notes: "Повышение" },
    ],
  } as Record<string, EmploymentRecord[]>,
  disciplines: {
    "2": [
      { id: 1, name: "Базы данных", group: "ИС-21", semester: "6" },
      { id: 2, name: "Веб-разработка", group: "ИС-22", semester: "4" },
    ],
  },
  students: [
    { id: 1, name: "Иван Петров", group: "ИС-21", course: 3, status: "active", registration_date: "2023-09-01" },
    { id: 2, name: "Елена Кузнецова", group: "ИС-22", course: 2, status: "active", registration_date: "2024-09-01" },
    { id: 3, name: "Павел Орлов", group: "ПИ-11", course: 1, status: "expelled", registration_date: "2025-09-01" },
  ] as Student[],
  groups: [
    { id: 1, name: "ИС-21", course: 3, curator_id: 2, curator_name: "Мария Иванова", student_count: 24 },
    { id: 2, name: "ИС-22", course: 2, curator_id: 3, curator_name: "Дмитрий Соколов", student_count: 27 },
    { id: 3, name: "ПИ-11", course: 1, curator_id: 2, curator_name: "Мария Иванова", student_count: 30 },
  ] as Group[],
  exams: [
    { id: 1, discipline: "Базы данных", date: "2026-06-10", time: "10:00", room: "305", teacher_id: 2, teacher_name: "Мария Иванова" },
    { id: 2, discipline: "Математический анализ", date: "2026-06-12", time: "12:00", room: "201", teacher_id: 3, teacher_name: "Дмитрий Соколов" },
  ] as ExamSchedule[],
  departments: [
    { id: 1, name: "Прикладная информатика", office: "305", head_id: 6, head_name: "Олег Васильев" },
    { id: 2, name: "Математический анализ", office: "201", head_id: 3, head_name: "Дмитрий Соколов" },
  ] as Department[],
  grades: [
    { id: 1, student_name: "Иван Петров", discipline: "Базы данных", grade: 5, date: "2026-05-20" },
    { id: 2, student_name: "Елена Кузнецова", discipline: "Веб-разработка", grade: 4, date: "2026-05-21" },
  ] as GradeReport[],
  attendance: [
    { id: 1, student_name: "Иван Петров", group: "ИС-21", attendance_percent: 92 },
    { id: 2, student_name: "Елена Кузнецова", group: "ИС-22", attendance_percent: 86 },
  ] as AttendanceReport[],
  acl: [] as ACLPermission[],
};

mockState.acl = ["student", "teacher", "dean", "admin"].flatMap((role) =>
  ["users", "students", "grades", "attendance", "exam_schedule", "debts", "reports", "profile"].map((resource) => ({
    resource,
    role,
    read: role === "admin" || role === "dean" || resource === "profile",
    write: role === "admin" || (role === "dean" && resource !== "users"),
    delete: role === "admin",
  }))
);

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value));
}

function addAudit(action: string) {
  mockState.audit.unshift({
    id: Date.now(),
    user_id: mockUser.id,
    user_name: mockUser.name,
    action,
    timestamp: new Date().toISOString(),
  });
}

function nextId(_prefix: string): number {
  return Date.now();
}

function mockResponse(method: string, path: string, data?: any) {
  const url = new URL(path, "http://mock.local");
  const pathname = url.pathname;

  if (method === "POST" && pathname === "/auth/login") {
    if (data?.email === "admin@dekanat.local" && data?.password === "admin123") {
      return { access_token: "mock-jwt-token", user: clone(mockUser) };
    }
    throw { status: 401, message: "Неверный email или пароль" };
  }

  if (method === "GET" && pathname === "/admin/stats") {
    return {
      total_users: mockState.users.length,
      students: mockState.users.filter((user) => user.role === "student").length,
      teachers: mockState.users.filter((user) => user.role === "teacher").length,
      pending_approval: mockState.users.filter((user) => user.status === "pending").length,
    };
  }

  if (method === "GET" && pathname === "/admin/audit") {
    const limit = Number(url.searchParams.get("limit") || mockState.audit.length);
    return { logs: clone(mockState.audit.slice(0, limit)), total: mockState.audit.length };
  }

  if (method === "GET" && pathname === "/admin/users") return { users: clone(mockState.users) };
  if (method === "POST" && pathname === "/admin/users") {
    if (mockState.users.some((user) => user.email === data.email)) throw { status: 409, message: "Email уже занят" };
    const user: User = { id: nextId("u"), name: data.name, email: data.email, role: data.role, status: "active", created_at: new Date().toISOString() };
    mockState.users.unshift(user);
    addAudit(`Создала пользователя ${user.name}`);
    return clone(user);
  }


  const userMatch = pathname.match(/^\/admin\/users\/([^/]+)$/);
  if (userMatch && method === "PUT") {
    const uid = Number(userMatch[1]);
    const user = mockState.users.find((item) => item.id === uid);
    if (!user) throw { status: 404, message: "Пользователь не найден" };
    Object.assign(user, data);
    if (mockState.teachers[String(user.id)]) Object.assign(mockState.teachers[String(user.id)], { name: user.name, email: user.email });
    addAudit(`Обновила пользователя ${user.name}`);
    return clone(user);
  }
  if (userMatch && method === "DELETE") {
    const uid = Number(userMatch[1]);
    const user = mockState.users.find((item) => item.id === uid);
    mockState.users = mockState.users.filter((item) => item.id !== uid);
    addAudit(`Удалила пользователя ${user?.name || userMatch[1]}`);
    return null;
  }

  const userActionMatch = pathname.match(/^\/admin\/users\/([^/]+)\/(block|approve|reset-password)$/);
  if (userActionMatch && method === "PATCH") {
    const uid = Number(userActionMatch[1]);
    const user = mockState.users.find((item) => item.id === uid);
    if (!user) throw { status: 404, message: "Пользователь не найден" };
    if (userActionMatch[2] === "block") user.status = user.status === "blocked" ? "active" : "blocked";
    if (userActionMatch[2] === "approve") user.status = "active";
    addAudit(userActionMatch[2] === "reset-password" ? `Сбросила пароль ${user.name}` : `Изменила статус ${user.name}`);
    return clone(user);
  }

  const teacherMatch = pathname.match(/^\/admin\/users\/([^/]+)\/teacher$/);
  if (teacherMatch && method === "GET") return clone(mockState.teachers[teacherMatch[1]] || null);
  if (teacherMatch && method === "PUT") {
    const existing = mockState.teachers[teacherMatch[1]] || { id: nextId("t"), user_id: Number(teacherMatch[1]) };
    mockState.teachers[teacherMatch[1]] = { ...existing, ...data };
    addAudit(`Обновила карточку преподавателя ${data.name}`);
    return clone(mockState.teachers[teacherMatch[1]]);
  }

  const employmentMatch = pathname.match(/^\/admin\/users\/([^/]+)\/employment$/);
  if (employmentMatch && method === "GET") return { records: clone(mockState.employment[employmentMatch[1]] || []) };
  if (employmentMatch && method === "POST") {
    const record = { id: nextId("emp"), ...data };
    mockState.employment[employmentMatch[1]] = [record, ...(mockState.employment[employmentMatch[1]] || [])];
    addAudit("Добавила запись трудоустройства");
    return clone(record);
  }

  const disciplinesMatch = pathname.match(/^\/admin\/users\/([^/]+)\/disciplines$/);
  if (disciplinesMatch && method === "GET") return { disciplines: clone((mockState.disciplines as any)[disciplinesMatch[1]] || []) };

  if (method === "GET" && pathname === "/admin/students") return { students: clone(mockState.students) };
  const studentMatch = pathname.match(/^\/admin\/students\/([^/]+)$/);
  if (studentMatch && method === "PATCH") {
    const sid = Number(studentMatch[1]);
    const student = mockState.students.find((item) => item.id === sid);
    if (!student) throw { status: 404, message: "Студент не найден" };
    Object.assign(student, data);
    addAudit(`Обновила студента ${student.name}`);
    return clone(student);
  }

  const collections: Record<string, { list: any[]; key: string }> = {
    "/admin/groups": { list: mockState.groups, key: "groups" },
    "/admin/exam-schedule": { list: mockState.exams, key: "exams" },
    "/admin/departments": { list: mockState.departments, key: "departments" },
  };
  for (const [basePath, collection] of Object.entries(collections)) {
    if (pathname === basePath && method === "GET") return { [collection.key]: clone(collection.list) };
    if (pathname === basePath && method === "POST") {
      const item = { id: nextId(collection.key), ...data };
      collection.list.unshift(item);
      addAudit("Создала запись учебных данных");
      return clone(item);
    }
    const itemMatch = pathname.match(new RegExp(`^${basePath}/([^/]+)$`));
    if (itemMatch && method === "PUT") {
      const iid = Number(itemMatch[1]);
      const item = collection.list.find((entry) => entry.id === iid);
      if (!item) throw { status: 404, message: "Запись не найдена" };
      Object.assign(item, data);
      addAudit("Обновила запись учебных данных");
      return clone(item);
    }
    if (itemMatch && method === "DELETE") {
      const iid = Number(itemMatch[1]);
      collection.list.splice(collection.list.findIndex((entry) => entry.id === iid), 1);
      addAudit("Удалила запись учебных данных");
      return null;
    }
  }

  if (method === "GET" && pathname === "/reports/grades") return { reports: clone(mockState.grades) };
  if (method === "GET" && pathname === "/reports/attendance") return { reports: clone(mockState.attendance) };
  if (method === "GET" && pathname === "/admin/acl") return { permissions: clone(mockState.acl) };
  if (method === "PATCH" && pathname === "/admin/acl") {
    const permission = mockState.acl.find((item) => item.resource === data.resource && item.role === data.role);
    if (permission) (permission as any)[data.permission] = data.value;
    addAudit(`Изменила право ${data.role}/${data.resource}/${data.permission}`);
    return clone(permission);
  }

  throw { status: 404, message: `Mock endpoint not found: ${method} ${pathname}` };
}

export class ApiClient {
  private token: string | null = null;
  private useMock = false;

  setToken(token: string) {
    this.token = token;
  }

  getToken(): string | null {
    return this.token;
  }

  clearToken() {
    this.token = null;
  }

  private async request(method: string, path: string, data?: any, headers?: Record<string, string>) {
    if (this.useMock) return mockResponse(method, path, data);

    const requestHeaders: Record<string, string> = {
      "Content-Type": "application/json",
      ...headers,
    };

    if (this.token) requestHeaders.Authorization = `Bearer ${this.token}`;

    try {
      const response = await fetch(`${API_BASE}${path}`, {
        method,
        headers: requestHeaders,
        body: data && ["POST", "PUT", "PATCH"].includes(method) ? JSON.stringify(data) : undefined,
      });

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw { status: response.status, message: error.message || response.statusText, data: error };
      }

      if (response.status === 204) return null;
      return response.json().catch(() => null);
    } catch (error: any) {
      if (error?.status) throw error;
      throw { status: 503, message: "Сервер недоступен. Проверьте подключение к бекенду." };
    }
  }

  async get(path: string) {
    return this.request("GET", path);
  }

  async post(path: string, data: any) {
    return this.request("POST", path, data);
  }

  async put(path: string, data: any) {
    return this.request("PUT", path, data);
  }

  async patch(path: string, data: any = {}) {
    return this.request("PATCH", path, data);
  }

  async delete(path: string) {
    return this.request("DELETE", path);
  }

  async download(path: string) {
    if (this.useMock) {
      return new Blob(["Демо-экспорт;данные\n"], {
        type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      });
    }

    try {
      const response = await fetch(`${API_BASE}${path}`, {
        headers: this.token ? { Authorization: `Bearer ${this.token}` } : {},
      });

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw { status: response.status, message: error.message || response.statusText, data: error };
      }

      return response.blob();
    } catch (error: any) {
      if (error?.status) throw error;
      throw { status: 503, message: "Сервер недоступен." };
    }
  }
}

export const apiClient = new ApiClient();
