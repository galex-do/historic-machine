-- +goose Up
ALTER TABLE event_datasets ADD COLUMN IF NOT EXISTS modified BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE event_datasets DROP COLUMN IF EXISTS modified;
