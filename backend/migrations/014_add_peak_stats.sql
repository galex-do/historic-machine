-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS peak_stats (
    id INTEGER PRIMARY KEY DEFAULT 1,
    peak_concurrent_sessions INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT single_row CHECK (id = 1)
);

-- Insert initial row
INSERT INTO peak_stats (id, peak_concurrent_sessions, updated_at)
VALUES (1, 0, NOW())
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS peak_stats;
-- +goose StatementEnd
