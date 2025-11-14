-- +goose Up
-- Add anonymous_sessions table for tracking unauthenticated users
CREATE TABLE IF NOT EXISTS anonymous_sessions (
    session_id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP
);

-- Add index for efficient active session queries
CREATE INDEX IF NOT EXISTS idx_anonymous_sessions_last_seen ON anonymous_sessions(last_seen_at);

-- +goose Down
DROP INDEX IF EXISTS idx_anonymous_sessions_last_seen;
DROP TABLE IF EXISTS anonymous_sessions;
