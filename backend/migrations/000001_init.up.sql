-- migrations/000001_init.up.sql

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create user_roles table (user to role mapping)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id),
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create resources table
CREATE TABLE IF NOT EXISTS resources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create ACL entries table
CREATE TABLE IF NOT EXISTS acl_entries (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    resource_id INT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    can_read BOOLEAN NOT NULL DEFAULT false,
    can_write BOOLEAN NOT NULL DEFAULT false,
    can_delete BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, resource_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_acl_entries_role_id ON acl_entries(role_id);
CREATE INDEX IF NOT EXISTS idx_acl_entries_resource_id ON acl_entries(resource_id);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator with full access'),
    ('teacher', 'Teacher role'),
    ('student', 'Student role'),
    ('dean', 'Dean role')
ON CONFLICT (name) DO NOTHING;

-- Insert default resources
INSERT INTO resources (name, description) VALUES
    ('admin', 'Admin panel and ACL management'),
    ('students', 'Student records'),
    ('grades', 'Grade management'),
    ('schedule', 'Schedule management'),
    ('teachers', 'Teacher records'),
    ('profile', 'User profile'),
    ('reports', 'System reports')
ON CONFLICT (name) DO NOTHING;

-- Insert default ACL rules
-- Admin: full access to all resources
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete) 
SELECT r.id, res.id, true, true, true 
FROM roles r, resources res 
WHERE r.name = 'admin'
ON CONFLICT (role_id, resource_id) DO NOTHING;

-- Student: can read own profile, grades, schedule; can write to own profile
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete) 
SELECT r.id, res.id, true, false, false 
FROM roles r, resources res 
WHERE r.name = 'student' AND res.name IN ('profile', 'grades', 'schedule')
ON CONFLICT (role_id, resource_id) DO NOTHING;

INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete) 
SELECT r.id, res.id, true, true, false 
FROM roles r, resources res 
WHERE r.name = 'student' AND res.name = 'profile'
ON CONFLICT (role_id, resource_id) DO NOTHING;

-- Teacher: can read students, grades, schedule; can write grades
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete) 
SELECT r.id, res.id, true, false, false 
FROM roles r, resources res 
WHERE r.name = 'teacher' AND res.name IN ('students', 'schedule')
ON CONFLICT (role_id, resource_id) DO NOTHING;

INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete) 
SELECT r.id, res.id, true, true, false 
FROM roles r, resources res 
WHERE r.name = 'teacher' AND res.name = 'grades'
ON CONFLICT (role_id, resource_id) DO NOTHING;

-- Dean: can read and write most resources
INSERT INTO acl_entries (role_id, resource_id, can_read, can_write, can_delete) 
SELECT r.id, res.id, true, true, false 
FROM roles r, resources res 
WHERE r.name = 'dean'
ON CONFLICT (role_id, resource_id) DO NOTHING;
