-- +goose Up
-- Add users and authentication system

-- Create access levels enum type
CREATE TYPE access_level AS ENUM ('guest', 'user', 'admin', 'super');

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    access_level access_level NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Create sessions table for JWT token management
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create indexes for users
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_access_level ON users(access_level);
CREATE INDEX idx_users_active ON users(is_active);

-- Create indexes for sessions
CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON user_sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_sessions_active ON user_sessions(is_active);

-- Add user tracking to events table
ALTER TABLE events ADD COLUMN created_by INTEGER REFERENCES users(id);
ALTER TABLE events ADD COLUMN updated_by INTEGER REFERENCES users(id);

-- Create index for event creators
CREATE INDEX idx_events_created_by ON events(created_by);
CREATE INDEX idx_events_updated_by ON events(updated_by);

-- Insert default admin user (password: "scriptor" - hashed)
-- Password hash for "scriptor" using bcrypt with cost 12
-- Note: This is a default hash - should be generated in production
INSERT INTO users (username, email, password_hash, access_level, is_active) VALUES 
('scriptor', 'admin@historymap.local', '$2a$12$uU18rg9ZVPfRQDO9SreSh.EF4n3t01sMgckgsAE35hmyYXEFq.HoC', 'super', true);

-- +goose Down
DROP INDEX IF EXISTS idx_events_updated_by;
DROP INDEX IF EXISTS idx_events_created_by;
ALTER TABLE events DROP COLUMN IF EXISTS updated_by;
ALTER TABLE events DROP COLUMN IF EXISTS created_by;

DROP INDEX IF EXISTS idx_sessions_active;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_token_hash;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP TABLE IF EXISTS user_sessions;

DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_access_level;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS access_level;