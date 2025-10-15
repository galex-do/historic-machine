-- +goose Up
-- Add internationalization support for event names and descriptions

-- Add new locale-specific columns
ALTER TABLE events ADD COLUMN name_en VARCHAR(255);
ALTER TABLE events ADD COLUMN name_ru VARCHAR(255);
ALTER TABLE events ADD COLUMN description_en TEXT;
ALTER TABLE events ADD COLUMN description_ru TEXT;

-- Migrate existing data to EN locale (copy current name/description to name_en/description_en)
UPDATE events SET name_en = name WHERE name_en IS NULL;
UPDATE events SET description_en = description WHERE description_en IS NULL;

-- Copy EN values to RU as well (default behavior - same content for all locales initially)
UPDATE events SET name_ru = name_en WHERE name_ru IS NULL;
UPDATE events SET description_ru = description_en WHERE description_ru IS NULL;

-- Make locale-specific fields NOT NULL now that they have data
ALTER TABLE events ALTER COLUMN name_en SET NOT NULL;
ALTER TABLE events ALTER COLUMN name_ru SET NOT NULL;

-- description can remain nullable since it was nullable before

-- Add indexes for locale-specific fields
CREATE INDEX idx_events_name_en ON events(name_en);
CREATE INDEX idx_events_name_ru ON events(name_ru);

-- Update the events view to include locale-specific data
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
GROUP BY e.id, e.name, e.description, e.latitude, e.longitude, e.event_date, e.era, e.lens_type, e.dataset_id, e.created_at, e.updated_at, e.created_by, e.updated_by, e.name_en, e.name_ru, e.description_en, e.description_ru
ORDER BY astronomical_year ASC;

-- +goose Down
-- Remove internationalization support

-- Drop the view first
DROP VIEW IF EXISTS events_with_display_dates;

-- Drop indexes
DROP INDEX IF EXISTS idx_events_name_en;
DROP INDEX IF EXISTS idx_events_name_ru;

-- Remove locale-specific columns
ALTER TABLE events DROP COLUMN IF EXISTS name_en;
ALTER TABLE events DROP COLUMN IF EXISTS name_ru;
ALTER TABLE events DROP COLUMN IF EXISTS description_en;
ALTER TABLE events DROP COLUMN IF EXISTS description_ru;

-- Recreate original view
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