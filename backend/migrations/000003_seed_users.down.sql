-- =============================================================
-- ROLLBACK: Remove seed data
-- =============================================================

BEGIN;

-- Remove audit log entries (except those we added)
DELETE FROM audit_log 
WHERE user_id IN (SELECT id FROM users WHERE email IN (
    'admin@dekanat.local', 'dean@dekanat.local', 'ivanova@dekanat.local',
    'sokolov@dekanat.local', 'kovalenko@dekanat.local', 'petrov@student.local',
    'kuznetsova@student.local', 'orlov@student.local', 'smirnova@student.local',
    'volkov@student.local'
)) AND metadata != '{}'::JSONB;

-- Remove user roles (keep composite keys)
DELETE FROM user_roles 
WHERE user_id IN (SELECT id FROM users WHERE email IN (
    'dean@dekanat.local', 'ivanova@dekanat.local', 'sokolov@dekanat.local',
    'kovalenko@dekanat.local', 'petrov@student.local', 'kuznetsova@student.local',
    'orlov@student.local', 'smirnova@student.local', 'volkov@student.local'
));

-- Remove users (keep admin)
DELETE FROM users WHERE email IN (
    'dean@dekanat.local', 'ivanova@dekanat.local', 'sokolov@dekanat.local',
    'kovalenko@dekanat.local', 'petrov@student.local', 'kuznetsova@student.local',
    'orlov@student.local', 'smirnova@student.local', 'volkov@student.local'
);

-- Remove ACL entries for non-admin roles
DELETE FROM acl_entries 
WHERE role_id IN (SELECT id FROM roles WHERE name IN ('dean', 'teacher', 'student'));

-- Remove roles (keep only admin if needed)
DELETE FROM roles WHERE name IN ('dean', 'teacher', 'student');

-- Remove resources
DELETE FROM resources WHERE name IN (
    'users', 'students', 'teachers', 'grades', 'attendance', 
    'exam_schedule', 'reports', 'audit', 'acl', 'profile'
);

COMMIT;
