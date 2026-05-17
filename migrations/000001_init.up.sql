-- migrations/000001_init.up.sql

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- TODO: Create refresh tokens table if needed
-- CREATE TABLE IF NOT EXISTS refresh_tokens (
--     id BIGSERIAL PRIMARY KEY,
--     user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     token_hash VARCHAR(255) NOT NULL,
--     expires_at TIMESTAMP NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

-- TODO: Create audit log table if needed
