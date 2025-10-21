-- +goose Up
-- Fix astronomical_year calculation to properly handle month/day ordering within BC years
-- For BC dates, earlier months (Jan) happened AFTER later months (Dec) chronologically

DROP VIEW IF EXISTS events_with_display_dates;

CREATE VIEW events_with_display_dates AS
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
        WHEN e.era = 'BC' THEN 
            -- For BC: negate year, then subtract month/day contribution
            -- (Earlier months like January happened AFTER later months like December)
            (EXTRACT(YEAR FROM e.event_date) * -1) + 1 
            - (EXTRACT(MONTH FROM e.event_date) / 12.0) 
            - (EXTRACT(DAY FROM e.event_date) / 365.0)
        ELSE 
            -- For AD: positive year plus month/day contribution
            EXTRACT(YEAR FROM e.event_date) 
            + (EXTRACT(MONTH FROM e.event_date) / 12.0) 
            + (EXTRACT(DAY FROM e.event_date) / 365.0)
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
GROUP BY e.id, e.name, e.description, e.latitude, e.longitude, e.event_date, e.era, e.lens_type, e.source, e.dataset_id, e.created_at, e.updated_at, e.created_by, e.updated_by, e.name_en, e.name_ru, e.description_en, e.description_ru
ORDER BY astronomical_year ASC;

-- +goose Down
-- Revert to the previous view definition without month/day precision

DROP VIEW IF EXISTS events_with_display_dates;

CREATE VIEW events_with_display_dates AS
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
GROUP BY e.id, e.name, e.description, e.latitude, e.longitude, e.event_date, e.era, e.lens_type, e.source, e.dataset_id, e.created_at, e.updated_at, e.created_by, e.updated_by;
