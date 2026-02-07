-- +goose Up
ALTER TABLE tags ADD COLUMN weight INTEGER NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE tags DROP COLUMN weight;
