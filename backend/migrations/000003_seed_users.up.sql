-- =============================================================
-- SEED DATA: Users, Roles, ACL, and Audit Logs
-- =============================================================

BEGIN;

-- =============================================================
-- 1. INSERT ROLES (if not already present)
-- =============================================================

INSERT INTO roles (name, description) VALUES
    ('admin', 'Администратор системы'),
    ('dean', 'Декан факультета'),
    ('teacher', 'Преподаватель'),
    ('student', 'Студент')
ON CONFLICT DO NOTHING;

-- =============================================================
-- 2. INSERT RESOURCES
-- =============================================================

INSERT INTO resources (name, description) VALUES
    ('users', 'Управление пользователями'),
    ('students', 'Управление студентами'),
    ('teachers', 'Управление преподавателями'),
    ('grades', 'Просмотр и ввод оценок'),
    ('attendance', 'Управление посещаемостью'),
    ('exam_schedule', 'Расписание экзаменов'),
    ('reports', 'Формирование отчётов'),
    ('audit', 'Просмотр журнала аудита'),
    ('acl', 'Управление правами доступа'),
    ('profile', 'Редактирование профиля')
ON CONFLICT DO NOTHING;

-- =============================================================
-- 3. INSERT ACL ENTRIES
-- =============================================================

-- ADMIN: Full access to all resources
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete)
SELECT r.id, res.id, TRUE, TRUE, TRUE
FROM roles r, resources res
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- DEAN: Read/Write access (no delete on users)
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete)
SELECT r.id, res.id, TRUE, TRUE, FALSE
FROM roles r, resources res
WHERE r.name = 'dean' 
  AND res.name NOT IN ('audit', 'acl', 'users')
ON CONFLICT DO NOTHING;

-- DEAN: Read-only access to audit and users
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete)
SELECT r.id, res.id, TRUE, FALSE, FALSE
FROM roles r, resources res
WHERE r.name = 'dean' 
  AND res.name IN ('audit', 'users')
ON CONFLICT DO NOTHING;

-- TEACHER: Limited access
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete)
SELECT r.id, res.id, TRUE, TRUE, FALSE
FROM roles r, resources res
WHERE r.name = 'teacher' 
  AND res.name IN ('students', 'grades', 'attendance', 'exam_schedule', 'profile')
ON CONFLICT DO NOTHING;

-- STUDENT: Very limited access
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete)
SELECT r.id, res.id, TRUE, FALSE, FALSE
FROM roles r, resources res
WHERE r.name = 'student' 
  AND res.name IN ('grades', 'attendance', 'exam_schedule', 'profile')
ON CONFLICT DO NOTHING;

-- =============================================================
-- 4. INSERT TEST USERS (with various roles and statuses)
-- =============================================================

-- Admin user (already created in 000002_seed_admin.up.sql)
-- Password: admin123 (hash: $2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e)

-- DEAN: Олег Васильев
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('vasil_oleg', 'dean@dekanat.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'active', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'dean@dekanat.local' AND r.name = 'dean'
ON CONFLICT DO NOTHING;

-- TEACHER 1: Мария Иванова (active)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('ivanova_maria', 'ivanova@dekanat.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'active', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'ivanova@dekanat.local' AND r.name = 'teacher'
ON CONFLICT DO NOTHING;

-- TEACHER 2: Дмитрий Соколов (pending approval)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('sokolov_dmitry', 'sokolov@dekanat.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'pending', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'sokolov@dekanat.local' AND r.name = 'teacher'
ON CONFLICT DO NOTHING;

-- TEACHER 3: Петр Коваленко (blocked)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('kovalenko_petr', 'kovalenko@dekanat.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'blocked', false)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'kovalenko@dekanat.local' AND r.name = 'teacher'
ON CONFLICT DO NOTHING;

-- STUDENT 1: Иван Петров (active)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('petrov_ivan', 'petrov@student.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'active', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'petrov@student.local' AND r.name = 'student'
ON CONFLICT DO NOTHING;

-- STUDENT 2: Елена Кузнецова (pending approval)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('kuznetsova_elena', 'kuznetsova@student.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'pending', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'kuznetsova@student.local' AND r.name = 'student'
ON CONFLICT DO NOTHING;

-- STUDENT 3: Павел Орлов (pending approval)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('orlov_pavel', 'orlov@student.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'pending', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'orlov@student.local' AND r.name = 'student'
ON CONFLICT DO NOTHING;

-- STUDENT 4: Анна Смирнова (pending approval)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('smirnova_anna', 'smirnova@student.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'pending', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'smirnova@student.local' AND r.name = 'student'
ON CONFLICT DO NOTHING;

-- STUDENT 5: Максим Волков (pending approval)
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('volkov_maksim', 'volkov@student.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'pending', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'volkov@student.local' AND r.name = 'student'
ON CONFLICT DO NOTHING;

-- =============================================================
-- 5. INSERT AUDIT LOG ENTRIES (admin actions)
-- =============================================================

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '2'::TEXT, '{"role":"dean","email":"dean@dekanat.local"}'::JSONB, NOW() - INTERVAL '48 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '3'::TEXT, '{"role":"teacher","email":"ivanova@dekanat.local"}'::JSONB, NOW() - INTERVAL '47 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '4'::TEXT, '{"role":"teacher","email":"sokolov@dekanat.local"}'::JSONB, NOW() - INTERVAL '46 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Updated user status', 'users', '4'::TEXT, '{"status":"pending","reason":"approval pending"}'::JSONB, NOW() - INTERVAL '45 hours 30 minutes'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '5'::TEXT, '{"role":"teacher","email":"kovalenko@dekanat.local"}'::JSONB, NOW() - INTERVAL '44 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Blocked user', 'users', '5'::TEXT, '{"reason":"policy violation"}'::JSONB, NOW() - INTERVAL '43 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '6'::TEXT, '{"role":"student","email":"petrov@student.local"}'::JSONB, NOW() - INTERVAL '72 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '7'::TEXT, '{"role":"student","email":"kuznetsova@student.local"}'::JSONB, NOW() - INTERVAL '60 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '8'::TEXT, '{"role":"student","email":"orlov@student.local"}'::JSONB, NOW() - INTERVAL '48 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '9'::TEXT, '{"role":"student","email":"smirnova@student.local"}'::JSONB, NOW() - INTERVAL '36 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Registered new user', 'users', '10'::TEXT, '{"role":"student","email":"volkov@student.local"}'::JSONB, NOW() - INTERVAL '24 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Viewed audit logs', 'audit', NULL::TEXT, '{"page":1,"limit":50}'::JSONB, NOW() - INTERVAL '12 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'Updated ACL permissions', 'acl', '4'::TEXT, '{"role":"teacher","resource":"grades"}'::JSONB, NOW() - INTERVAL '6 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

INSERT INTO audit_log (user_id, action, resource, resource_id, metadata, created_at)
SELECT u.id, 'System backup created', 'audit', NULL::TEXT, '{"size":"2.4 GB","location":"backup_server"}'::JSONB, NOW() - INTERVAL '2 hours'
FROM users u WHERE u.email = 'admin@dekanat.local'
ON CONFLICT DO NOTHING;

COMMIT;
