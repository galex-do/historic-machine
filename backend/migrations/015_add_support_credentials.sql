-- +goose Up
-- Create support_credentials table for funding/donation options

CREATE TABLE IF NOT EXISTS support_credentials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    is_url BOOLEAN DEFAULT false,
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_support_credentials_active ON support_credentials(is_active, display_order);

-- +goose Down
DROP INDEX IF EXISTS idx_support_credentials_active;
DROP TABLE IF EXISTS support_credentials;
