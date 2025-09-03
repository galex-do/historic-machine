-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create datasets table
CREATE TABLE event_datasets (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    description TEXT,
    event_count INTEGER DEFAULT 0,
    uploaded_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add dataset foreign key to events table
ALTER TABLE events ADD COLUMN dataset_id INTEGER REFERENCES event_datasets(id) ON DELETE CASCADE;

-- Create indexes for performance
CREATE INDEX idx_events_dataset_id ON events(dataset_id);
CREATE INDEX idx_datasets_uploaded_by ON event_datasets(uploaded_by);
CREATE INDEX idx_datasets_created_at ON event_datasets(created_at);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove indexes
DROP INDEX IF EXISTS idx_datasets_created_at;
DROP INDEX IF EXISTS idx_datasets_uploaded_by;
DROP INDEX IF EXISTS idx_events_dataset_id;

-- Remove foreign key from events table
ALTER TABLE events DROP COLUMN IF EXISTS dataset_id;

-- Drop datasets table
DROP TABLE IF EXISTS event_datasets;