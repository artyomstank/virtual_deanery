-- migrations/000001_init.down.sql

-- Drop users table
DROP TABLE IF EXISTS users CASCADE;

-- TODO: Drop other tables if added
-- DROP TABLE IF EXISTS refresh_tokens CASCADE;
