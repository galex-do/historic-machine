-- +goose Up
ALTER TABLE tags ADD COLUMN IF NOT EXISTS border_color VARCHAR(7) DEFAULT NULL;

-- +goose Down
ALTER TABLE tags DROP COLUMN IF EXISTS border_color;
