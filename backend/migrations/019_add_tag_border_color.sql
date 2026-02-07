-- +goose Up
ALTER TABLE tags ADD COLUMN IF NOT EXISTS border_color VARCHAR(7) DEFAULT NULL;

-- Update the events view to include border_color in the tags JSON
CREATE OR REPLACE VIEW events_with_display_dates AS
SELECT 
  id, name, description, latitude, longitude, event_date, era, lens_type,
  created_at, updated_at, created_by, updated_by, dataset_id, source,
  name_en, name_ru, description_en, description_ru,
  CASE 
    WHEN era = 'BC' THEN 
      CONCAT(LPAD(EXTRACT(DAY FROM event_date)::TEXT, 2, '0'), '.', 
             LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.', 
             EXTRACT(YEAR FROM event_date)::TEXT, ' BC')
    ELSE 
      CONCAT(LPAD(EXTRACT(DAY FROM event_date)::TEXT, 2, '0'), '.', 
             LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.', 
             EXTRACT(YEAR FROM event_date)::TEXT, ' AD')
  END AS display_date,
  CASE 
    WHEN era = 'BC' THEN 
      EXTRACT(YEAR FROM event_date) * -1 + 1 - EXTRACT(MONTH FROM event_date) / 12.0 - EXTRACT(DAY FROM event_date) / 365.0
    ELSE 
      EXTRACT(YEAR FROM event_date) + EXTRACT(MONTH FROM event_date) / 12.0 + EXTRACT(DAY FROM event_date) / 365.0
  END AS astronomical_year,
  COALESCE(
    (SELECT json_agg(
      json_build_object(
        'id', t.id, 
        'name', t.name, 
        'description', t.description, 
        'color', t.color, 
        'border_color', t.border_color,
        'weight', t.weight
      ) ORDER BY t.weight DESC, t.name
    )
    FROM event_tags et2
    JOIN tags t ON et2.tag_id = t.id
    WHERE et2.event_id = e.id), 
    '[]'::json
  ) AS tags
FROM events e
GROUP BY id, name, description, latitude, longitude, event_date, era, lens_type, source, dataset_id, created_at, updated_at, created_by, updated_by, name_en, name_ru, description_en, description_ru
ORDER BY astronomical_year;

-- +goose Down
ALTER TABLE tags DROP COLUMN IF EXISTS border_color;

-- Restore the view without border_color
CREATE OR REPLACE VIEW events_with_display_dates AS
SELECT 
  id, name, description, latitude, longitude, event_date, era, lens_type,
  created_at, updated_at, created_by, updated_by, dataset_id, source,
  name_en, name_ru, description_en, description_ru,
  CASE 
    WHEN era = 'BC' THEN 
      CONCAT(LPAD(EXTRACT(DAY FROM event_date)::TEXT, 2, '0'), '.', 
             LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.', 
             EXTRACT(YEAR FROM event_date)::TEXT, ' BC')
    ELSE 
      CONCAT(LPAD(EXTRACT(DAY FROM event_date)::TEXT, 2, '0'), '.', 
             LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.', 
             EXTRACT(YEAR FROM event_date)::TEXT, ' AD')
  END AS display_date,
  CASE 
    WHEN era = 'BC' THEN 
      EXTRACT(YEAR FROM event_date) * -1 + 1 - EXTRACT(MONTH FROM event_date) / 12.0 - EXTRACT(DAY FROM event_date) / 365.0
    ELSE 
      EXTRACT(YEAR FROM event_date) + EXTRACT(MONTH FROM event_date) / 12.0 + EXTRACT(DAY FROM event_date) / 365.0
  END AS astronomical_year,
  COALESCE(
    (SELECT json_agg(
      json_build_object(
        'id', t.id, 
        'name', t.name, 
        'description', t.description, 
        'color', t.color, 
        'weight', t.weight
      ) ORDER BY t.weight DESC, t.name
    )
    FROM event_tags et2
    JOIN tags t ON et2.tag_id = t.id
    WHERE et2.event_id = e.id), 
    '[]'::json
  ) AS tags
FROM events e
GROUP BY id, name, description, latitude, longitude, event_date, era, lens_type, source, dataset_id, created_at, updated_at, created_by, updated_by, name_en, name_ru, description_en, description_ru
ORDER BY astronomical_year;
