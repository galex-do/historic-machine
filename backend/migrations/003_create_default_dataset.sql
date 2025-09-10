-- +goose Up
-- Create default dataset and bind first 7 events

-- Insert default dataset for the first 7 historical events
INSERT INTO event_datasets (filename, description, event_count, uploaded_by, created_at, updated_at)
VALUES (
    'default_ancient_events.json',
    'Default dataset containing the first 7 foundational ancient civilization events',
    7,
    (SELECT id FROM users WHERE access_level = 'super' LIMIT 1),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Get the dataset ID for binding events
-- Update the first 7 events to be associated with the default dataset
UPDATE events 
SET dataset_id = (SELECT id FROM event_datasets WHERE filename = 'default_ancient_events.json')
WHERE id IN (9, 11, 12, 13, 14, 15, 17);

-- +goose Down
-- Remove default dataset binding and delete the dataset

-- Remove dataset association from the first 7 events
UPDATE events 
SET dataset_id = NULL
WHERE id IN (9, 11, 12, 13, 14, 15, 17);

-- Delete the default dataset
DELETE FROM event_datasets WHERE filename = 'default_ancient_events.json';