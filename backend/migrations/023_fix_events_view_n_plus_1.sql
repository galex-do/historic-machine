-- +goose Up
-- Replace the correlated subquery (N+1) in events_with_display_dates with a
-- single-pass LEFT JOIN + json_agg. Tag JSON shape and column order are
-- identical to the previous version so no Go or frontend changes are needed.

DROP VIEW IF EXISTS events_with_display_dates;

CREATE VIEW events_with_display_dates AS
SELECT
  e.id,
  e.name,
  e.description,
  e.latitude,
  e.longitude,
  e.event_date,
  e.era,
  e.lens_type,
  e.created_at,
  e.updated_at,
  e.created_by,
  e.updated_by,
  e.dataset_id,
  e.source,
  e.name_en,
  e.name_ru,
  e.description_en,
  e.description_ru,
  CASE
    WHEN e.era = 'BC' THEN
      CONCAT(LPAD(EXTRACT(DAY   FROM e.event_date)::TEXT, 2, '0'), '.',
             LPAD(EXTRACT(MONTH FROM e.event_date)::TEXT, 2, '0'), '.',
             EXTRACT(YEAR FROM e.event_date)::TEXT, ' BC')
    ELSE
      CONCAT(LPAD(EXTRACT(DAY   FROM e.event_date)::TEXT, 2, '0'), '.',
             LPAD(EXTRACT(MONTH FROM e.event_date)::TEXT, 2, '0'), '.',
             EXTRACT(YEAR FROM e.event_date)::TEXT, ' AD')
  END AS display_date,
  CASE
    WHEN e.era = 'BC' THEN
      EXTRACT(YEAR FROM e.event_date) * -1 + 1
        - EXTRACT(MONTH FROM e.event_date) / 12.0
        - EXTRACT(DAY   FROM e.event_date) / 365.0
    ELSE
      EXTRACT(YEAR FROM e.event_date)
        + EXTRACT(MONTH FROM e.event_date) / 12.0
        + EXTRACT(DAY   FROM e.event_date) / 365.0
  END AS astronomical_year,
  -- Single-pass aggregation replaces the correlated subquery that ran once
  -- per event row (N+1). FILTER (WHERE t.id IS NOT NULL) handles events that
  -- have no tags and would otherwise produce a [null] array.
  COALESCE(
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'id',           t.id,
        'name',         t.name,
        'description',  t.description,
        'color',        t.color,
        'border_color', t.border_color,
        'key_color',    t.key_color,
        'emoji',        t.emoji,
        'weight',       t.weight
      ) ORDER BY t.weight DESC, t.name
    ) FILTER (WHERE t.id IS NOT NULL),
    '[]'::json
  ) AS tags
FROM events e
LEFT JOIN event_tags et ON et.event_id = e.id
LEFT JOIN tags       t  ON t.id = et.tag_id
GROUP BY
  e.id, e.name, e.description, e.latitude, e.longitude,
  e.event_date, e.era, e.lens_type, e.source, e.dataset_id,
  e.created_at, e.updated_at, e.created_by, e.updated_by,
  e.name_en, e.name_ru, e.description_en, e.description_ru
ORDER BY astronomical_year;

-- +goose Down
-- Restore the 022 correlated-subquery version

DROP VIEW IF EXISTS events_with_display_dates;

CREATE VIEW events_with_display_dates AS
SELECT
  id, name, description, latitude, longitude, event_date, era, lens_type,
  created_at, updated_at, created_by, updated_by, dataset_id, source,
  name_en, name_ru, description_en, description_ru,
  CASE
    WHEN era = 'BC' THEN
      CONCAT(LPAD(EXTRACT(DAY   FROM event_date)::TEXT, 2, '0'), '.',
             LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.',
             EXTRACT(YEAR FROM event_date)::TEXT, ' BC')
    ELSE
      CONCAT(LPAD(EXTRACT(DAY   FROM event_date)::TEXT, 2, '0'), '.',
             LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.',
             EXTRACT(YEAR FROM event_date)::TEXT, ' AD')
  END AS display_date,
  CASE
    WHEN era = 'BC' THEN
      EXTRACT(YEAR FROM event_date) * -1 + 1
        - EXTRACT(MONTH FROM event_date) / 12.0
        - EXTRACT(DAY   FROM event_date) / 365.0
    ELSE
      EXTRACT(YEAR FROM event_date)
        + EXTRACT(MONTH FROM event_date) / 12.0
        + EXTRACT(DAY   FROM event_date) / 365.0
  END AS astronomical_year,
  COALESCE(
    (SELECT json_agg(
      json_build_object(
        'id',           t.id,
        'name',         t.name,
        'description',  t.description,
        'color',        t.color,
        'border_color', t.border_color,
        'key_color',    t.key_color,
        'emoji',        t.emoji,
        'weight',       t.weight
      ) ORDER BY t.weight DESC, t.name
    )
    FROM event_tags et2
    JOIN tags t ON et2.tag_id = t.id
    WHERE et2.event_id = e.id),
    '[]'::json
  ) AS tags
FROM events e
GROUP BY
  id, name, description, latitude, longitude, event_date, era, lens_type,
  source, dataset_id, created_at, updated_at, created_by, updated_by,
  name_en, name_ru, description_en, description_ru
ORDER BY astronomical_year;
