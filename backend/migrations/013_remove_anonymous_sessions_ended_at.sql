-- +goose Up
-- Remove ended_at field from anonymous_sessions as it's never set
-- Use last_seen_at instead to calculate session duration
ALTER TABLE anonymous_sessions DROP COLUMN IF EXISTS ended_at;

-- +goose Down
-- Restore ended_at field if migration is rolled back
ALTER TABLE anonymous_sessions ADD COLUMN ended_at TIMESTAMP;
