-- +goose Up
-- Add internationalization support for date template groups and templates

-- Add locale-specific columns to date_template_groups
ALTER TABLE date_template_groups ADD COLUMN name_en VARCHAR(255);
ALTER TABLE date_template_groups ADD COLUMN name_ru VARCHAR(255);
ALTER TABLE date_template_groups ADD COLUMN description_en TEXT;
ALTER TABLE date_template_groups ADD COLUMN description_ru TEXT;

-- Migrate existing data to EN locale (copy current name/description to name_en/description_en)
UPDATE date_template_groups SET name_en = name WHERE name_en IS NULL;
UPDATE date_template_groups SET description_en = description WHERE description_en IS NULL;

-- Set RU values to same as EN initially (will be updated later with actual translations)
UPDATE date_template_groups SET name_ru = name_en WHERE name_ru IS NULL;
UPDATE date_template_groups SET description_ru = description_en WHERE description_ru IS NULL;

-- Make locale-specific fields NOT NULL for date_template_groups
ALTER TABLE date_template_groups ALTER COLUMN name_en SET NOT NULL;
ALTER TABLE date_template_groups ALTER COLUMN name_ru SET NOT NULL;

-- Add locale-specific columns to date_templates
ALTER TABLE date_templates ADD COLUMN name_en VARCHAR(255);
ALTER TABLE date_templates ADD COLUMN name_ru VARCHAR(255);
ALTER TABLE date_templates ADD COLUMN description_en TEXT;
ALTER TABLE date_templates ADD COLUMN description_ru TEXT;

-- Migrate existing data to EN locale
UPDATE date_templates SET name_en = name WHERE name_en IS NULL;
UPDATE date_templates SET description_en = description WHERE description_en IS NULL;

-- Set RU values to same as EN initially
UPDATE date_templates SET name_ru = name_en WHERE name_ru IS NULL;
UPDATE date_templates SET description_ru = description_en WHERE description_ru IS NULL;

-- Make locale-specific fields NOT NULL for date_templates
ALTER TABLE date_templates ALTER COLUMN name_en SET NOT NULL;
ALTER TABLE date_templates ALTER COLUMN name_ru SET NOT NULL;

-- Add indexes for locale-specific fields
CREATE INDEX idx_template_groups_name_en ON date_template_groups(name_en);
CREATE INDEX idx_template_groups_name_ru ON date_template_groups(name_ru);
CREATE INDEX idx_templates_name_en ON date_templates(name_en);
CREATE INDEX idx_templates_name_ru ON date_templates(name_ru);

-- Update the date_templates_with_display view to include locale-specific data
DROP VIEW IF EXISTS date_templates_with_display;
CREATE OR REPLACE VIEW date_templates_with_display AS
SELECT 
    dt.id,
    dt.group_id,
    dtg.name as group_name,
    dt.name,
    dt.description,
    dt.name_en,
    dt.name_ru,
    dt.description_en,
    dt.description_ru,
    dtg.name_en as group_name_en,
    dtg.name_ru as group_name_ru,
    dt.start_date,
    dt.start_era,
    dt.end_date,
    dt.end_era,
    dt.display_order,
    -- Format display dates  
    CASE 
        WHEN dt.start_era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.start_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.start_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.start_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.start_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.start_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.start_date)::TEXT, ' AD'
            )
    END AS start_display_date,
    CASE 
        WHEN dt.end_era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.end_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.end_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.end_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.end_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.end_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.end_date)::TEXT, ' AD'
            )
    END AS end_display_date,
    -- Calculate astronomical years for sorting
    CASE 
        WHEN dt.start_era = 'BC' THEN (EXTRACT(YEAR FROM dt.start_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM dt.start_date)
    END AS start_astronomical_year,
    CASE 
        WHEN dt.end_era = 'BC' THEN (EXTRACT(YEAR FROM dt.end_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM dt.end_date)
    END AS end_astronomical_year
FROM date_templates dt
JOIN date_template_groups dtg ON dt.group_id = dtg.id
ORDER BY dtg.display_order, dt.display_order;

-- +goose Down
-- Remove internationalization support

-- Drop the view first
DROP VIEW IF EXISTS date_templates_with_display;

-- Drop indexes
DROP INDEX IF EXISTS idx_template_groups_name_en;
DROP INDEX IF EXISTS idx_template_groups_name_ru;
DROP INDEX IF EXISTS idx_templates_name_en;
DROP INDEX IF EXISTS idx_templates_name_ru;

-- Remove locale-specific columns from date_template_groups
ALTER TABLE date_template_groups DROP COLUMN IF EXISTS name_en;
ALTER TABLE date_template_groups DROP COLUMN IF EXISTS name_ru;
ALTER TABLE date_template_groups DROP COLUMN IF EXISTS description_en;
ALTER TABLE date_template_groups DROP COLUMN IF EXISTS description_ru;

-- Remove locale-specific columns from date_templates
ALTER TABLE date_templates DROP COLUMN IF EXISTS name_en;
ALTER TABLE date_templates DROP COLUMN IF EXISTS name_ru;
ALTER TABLE date_templates DROP COLUMN IF EXISTS description_en;
ALTER TABLE date_templates DROP COLUMN IF EXISTS description_ru;

-- Recreate original view
CREATE OR REPLACE VIEW date_templates_with_display AS
SELECT 
    dt.id,
    dt.group_id,
    dtg.name as group_name,
    dt.name,
    dt.description,
    dt.start_date,
    dt.start_era,
    dt.end_date,
    dt.end_era,
    dt.display_order,
    -- Format display dates  
    CASE 
        WHEN dt.start_era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.start_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.start_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.start_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.start_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.start_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.start_date)::TEXT, ' AD'
            )
    END AS start_display_date,
    CASE 
        WHEN dt.end_era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.end_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.end_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.end_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM dt.end_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM dt.end_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM dt.end_date)::TEXT, ' AD'
            )
    END AS end_display_date,
    -- Calculate astronomical years for sorting
    CASE 
        WHEN dt.start_era = 'BC' THEN (EXTRACT(YEAR FROM dt.start_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM dt.start_date)
    END AS start_astronomical_year,
    CASE 
        WHEN dt.end_era = 'BC' THEN (EXTRACT(YEAR FROM dt.end_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM dt.end_date)
    END AS end_astronomical_year
FROM date_templates dt
JOIN date_template_groups dtg ON dt.group_id = dtg.id
ORDER BY dtg.display_order, dt.display_order;
