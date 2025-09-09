-- +goose Up
-- Add ownership tracking columns to events table
ALTER TABLE events 
ADD COLUMN created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN updated_by INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Create indexes for performance
CREATE INDEX idx_events_created_by ON events(created_by);
CREATE INDEX idx_events_updated_by ON events(updated_by);

-- Update the events_with_display_dates view to include ownership fields
DROP VIEW IF EXISTS events_with_display_dates;

CREATE OR REPLACE VIEW events_with_display_dates AS
SELECT 
    e.*,
    CASE 
        WHEN e.era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM e.event_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM e.event_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM e.event_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM e.event_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM e.event_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM e.event_date)::TEXT, ' AD'
            )
    END AS display_date,
    CASE 
        WHEN e.era = 'BC' THEN (EXTRACT(YEAR FROM e.event_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM e.event_date)
    END AS astronomical_year,
    -- Include aggregated tags as JSON
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'id', t.id,
                'name', t.name,
                'description', t.description,
                'color', t.color
            )
        ) FILTER (WHERE t.id IS NOT NULL),
        '[]'::json
    ) AS tags
FROM events e
LEFT JOIN event_tags et ON e.id = et.event_id
LEFT JOIN tags t ON et.tag_id = t.id
GROUP BY e.id, e.name, e.description, e.latitude, e.longitude, e.event_date, e.era, e.lens_type, e.dataset_id, e.created_at, e.updated_at, e.created_by, e.updated_by
ORDER BY astronomical_year ASC;

-- +goose Down
-- Remove the ownership tracking columns
DROP VIEW IF EXISTS events_with_display_dates;

-- Recreate the view without ownership fields
CREATE OR REPLACE VIEW events_with_display_dates AS
SELECT 
    e.*,
    CASE 
        WHEN e.era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM e.event_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM e.event_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM e.event_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM e.event_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM e.event_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM e.event_date)::TEXT, ' AD'
            )
    END AS display_date,
    CASE 
        WHEN e.era = 'BC' THEN (EXTRACT(YEAR FROM e.event_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM e.event_date)
    END AS astronomical_year,
    -- Include aggregated tags as JSON
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'id', t.id,
                'name', t.name,
                'description', t.description,
                'color', t.color
            )
        ) FILTER (WHERE t.id IS NOT NULL),
        '[]'::json
    ) AS tags
FROM events e
LEFT JOIN event_tags et ON e.id = et.event_id
LEFT JOIN tags t ON et.tag_id = t.id
GROUP BY e.id, e.name, e.description, e.latitude, e.longitude, e.event_date, e.era, e.lens_type, e.dataset_id, e.created_at, e.updated_at
ORDER BY astronomical_year ASC;

-- Drop indexes
DROP INDEX IF EXISTS idx_events_created_by;
DROP INDEX IF EXISTS idx_events_updated_by;

-- Remove columns
ALTER TABLE events 
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS created_by;