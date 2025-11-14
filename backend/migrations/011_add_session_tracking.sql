-- +goose Up
-- Add session tracking columns to user_sessions table for statistics

ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS ended_at TIMESTAMP;

-- Create index for efficient active user queries
CREATE INDEX IF NOT EXISTS idx_user_sessions_last_seen ON user_sessions(last_seen_at) WHERE is_active = true;

-- +goose Down
DROP INDEX IF EXISTS idx_user_sessions_last_seen;
ALTER TABLE user_sessions DROP COLUMN IF EXISTS ended_at;
ALTER TABLE user_sessions DROP COLUMN IF EXISTS last_seen_at;
