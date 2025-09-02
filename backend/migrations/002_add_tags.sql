-- +goose Up
-- Add tags functionality to Historical Events Mapping

-- Create tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7) DEFAULT '#3b82f6', -- Hex color for tag display
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create junction table for many-to-many relationship between events and tags
CREATE TABLE event_tags (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, tag_id) -- Prevent duplicate tag assignments
);

-- Create indexes for performance
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_event_tags_event_id ON event_tags(event_id);
CREATE INDEX idx_event_tags_tag_id ON event_tags(tag_id);

-- Insert sample tags for historical categorization
INSERT INTO tags (name, description, color) VALUES 
('Roman', 'Events related to Roman civilization and empire', '#dc2626'),
('Greek', 'Events related to Ancient Greece and Hellenistic period', '#2563eb'),
('Military', 'Wars, battles, and military campaigns', '#7c2d12'),
('Political', 'Government changes, assassinations, political events', '#9333ea'),
('Cultural', 'Philosophy, arts, literature, and cultural developments', '#059669'),
('Byzantine', 'Events related to the Byzantine Empire', '#c2410c'),
('Medieval', 'Events from the Middle Ages period', '#4338ca'),
('Renaissance', 'Events from the Renaissance period', '#ea580c'),
('Ancient', 'Very old historical events from antiquity', '#0f172a'),
('Empire', 'Events related to imperial powers and empires', '#991b1b'),
('Conquest', 'Military conquests and territorial expansion', '#7f1d1d'),
('Philosophy', 'Events related to philosophical movements and thinkers', '#065f46');

-- Add sample tags to existing events
-- Foundation of Rome: Roman, Ancient, Political
INSERT INTO event_tags (event_id, tag_id) VALUES 
(1, (SELECT id FROM tags WHERE name = 'Roman')),
(1, (SELECT id FROM tags WHERE name = 'Ancient')),
(1, (SELECT id FROM tags WHERE name = 'Political'));

-- Battle of Marathon: Greek, Military, Ancient
INSERT INTO event_tags (event_id, tag_id) VALUES 
(2, (SELECT id FROM tags WHERE name = 'Greek')),
(2, (SELECT id FROM tags WHERE name = 'Military')),
(2, (SELECT id FROM tags WHERE name = 'Ancient'));

-- Birth of Plato: Greek, Cultural, Philosophy
INSERT INTO event_tags (event_id, tag_id) VALUES 
(3, (SELECT id FROM tags WHERE name = 'Greek')),
(3, (SELECT id FROM tags WHERE name = 'Cultural')),
(3, (SELECT id FROM tags WHERE name = 'Philosophy'));

-- Battle of Gaugamela: Military, Conquest, Ancient
INSERT INTO event_tags (event_id, tag_id) VALUES 
(4, (SELECT id FROM tags WHERE name = 'Military')),
(4, (SELECT id FROM tags WHERE name = 'Conquest')),
(4, (SELECT id FROM tags WHERE name = 'Ancient'));

-- Assassination of Julius Caesar: Roman, Political, Ancient
INSERT INTO event_tags (event_id, tag_id) VALUES 
(5, (SELECT id FROM tags WHERE name = 'Roman')),
(5, (SELECT id FROM tags WHERE name = 'Political')),
(5, (SELECT id FROM tags WHERE name = 'Ancient'));

-- Fall of Constantinople: Byzantine, Military, Medieval
INSERT INTO event_tags (event_id, tag_id) VALUES 
(6, (SELECT id FROM tags WHERE name = 'Byzantine')),
(6, (SELECT id FROM tags WHERE name = 'Military')),
(6, (SELECT id FROM tags WHERE name = 'Medieval'));

-- Discovery of America: Renaissance, Cultural
INSERT INTO event_tags (event_id, tag_id) VALUES 
(7, (SELECT id FROM tags WHERE name = 'Renaissance')),
(7, (SELECT id FROM tags WHERE name = 'Cultural'));

-- Update events view to include tags
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
    -- Aggregate tags as JSON array
    COALESCE(
        (SELECT JSON_AGG(
            JSON_BUILD_OBJECT(
                'id', t.id,
                'name', t.name,
                'description', t.description,
                'color', t.color
            )
        )
        FROM event_tags et
        JOIN tags t ON et.tag_id = t.id
        WHERE et.event_id = e.id),
        '[]'::json
    ) AS tags
FROM events e
ORDER BY astronomical_year ASC;

-- +goose Down
DROP VIEW IF EXISTS events_with_display_dates;
DROP TABLE IF EXISTS event_tags;
DROP TABLE IF EXISTS tags;

-- Recreate original view without tags
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
    END AS astronomical_year
FROM events e
ORDER BY astronomical_year ASC;