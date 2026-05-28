-- migrations/000001_init.down.sql

-- Drop in reverse order of dependencies
DROP TABLE IF EXISTS acl_entries CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS resources CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
