-- +goose Up
-- Add source field to events table for storing HTTP/HTTPS links to source information

ALTER TABLE events ADD COLUMN source TEXT;

-- +goose Down
-- Remove source field from events table

ALTER TABLE events DROP COLUMN source;