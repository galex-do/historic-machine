-- +goose Up
-- Complete database schema for Historical Events Mapping

-- Create users table for authentication
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    access_level VARCHAR(20) DEFAULT 'guest' CHECK (access_level IN ('guest', 'user', 'editor', 'super')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Create user sessions table for authentication
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Create tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7) DEFAULT '#3B82F6',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create event_datasets table
CREATE TABLE event_datasets (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    description TEXT,
    event_count INTEGER DEFAULT 0,
    uploaded_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create events table with all necessary columns
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    event_date DATE NOT NULL,
    era VARCHAR(2) DEFAULT 'AD' CHECK (era IN ('BC', 'AD')),
    lens_type VARCHAR(50) NOT NULL,
    dataset_id INTEGER REFERENCES event_datasets(id) ON DELETE SET NULL,
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    updated_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create date template groups
CREATE TABLE date_template_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create date templates
CREATE TABLE date_templates (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES date_template_groups(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    start_era VARCHAR(2) NOT NULL DEFAULT 'AD',
    end_date DATE NOT NULL,
    end_era VARCHAR(2) NOT NULL DEFAULT 'AD',
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create event_tags junction table
CREATE TABLE event_tags (
    event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, tag_id)
);

-- Create indexes
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_lens_type ON events(lens_type);
CREATE INDEX idx_events_era ON events(era);
CREATE INDEX idx_events_location ON events(latitude, longitude);
CREATE INDEX idx_events_dataset ON events(dataset_id);
CREATE INDEX idx_events_created_by ON events(created_by);
CREATE INDEX idx_events_updated_by ON events(updated_by);
CREATE INDEX idx_date_templates_group_id ON date_templates(group_id);
CREATE INDEX idx_date_templates_dates ON date_templates(start_date, end_date);
CREATE INDEX idx_date_template_groups_order ON date_template_groups(display_order);
CREATE INDEX idx_date_templates_order ON date_templates(display_order);
CREATE INDEX idx_event_tags_event ON event_tags(event_id);
CREATE INDEX idx_event_tags_tag ON event_tags(tag_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(token_hash);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);
CREATE INDEX idx_tags_name ON tags(name);


-- Insert default dataset for sample events
INSERT INTO event_datasets (filename, description, event_count, uploaded_by, created_at, updated_at)
VALUES (
    'default_sample_events.json',
    'Default sample historical events seeded during initial database setup',
    7,
    (SELECT id FROM users WHERE access_level = 'super' LIMIT 1),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Insert sample historical events and associate with default dataset
INSERT INTO events (name, description, latitude, longitude, event_date, era, lens_type, dataset_id) VALUES 
('Foundation of Rome', 'Legendary founding of Rome by Romulus', 41.9028, 12.4964, '0753-04-21', 'BC', 'historic', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json')),
('Battle of Marathon', 'Greeks defeat Persian invasion', 38.1462, 23.9608, '0490-09-12', 'BC', 'military', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json')),
('Birth of Plato', 'Great philosopher born in Athens', 37.9755, 23.7348, '0428-05-21', 'BC', 'cultural', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json')),
('Battle of Gaugamela', 'Alexander defeats Darius III', 36.3402, 43.1503, '0331-10-01', 'BC', 'military', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json')),
('Assassination of Julius Caesar', 'Caesar killed in the Roman Senate', 41.8955, 12.4823, '0044-03-15', 'BC', 'political', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json')),
('Fall of Constantinople', 'End of the Byzantine Empire', 41.0082, 28.9784, '1453-05-29', 'AD', 'military', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json')),
('Discovery of America', 'Columbus reaches the New World', 25.0343, -77.3963, '1492-10-12', 'AD', 'historic', (SELECT id FROM event_datasets WHERE filename = 'default_sample_events.json'));

-- Insert tags for events
INSERT INTO tags (name, description, color) VALUES
-- Existing base tags
('Ancient', 'Auto-generated tag for Ancient', '#3B82F6'),
('Cultural', 'Auto-generated tag for Cultural', '#10B981'),
('Greek', 'Auto-generated tag for Greek', '#059669'),
('Military', 'Auto-generated tag for Military', '#DC2626'),
('Political', 'Auto-generated tag for Political', '#7C3AED'),
('Roman', 'Auto-generated tag for Roman', '#EA580C'),
('Renaissance', 'Auto-generated tag for Renaissance', '#DB2777'),
('Philosophy', 'Auto-generated tag for Philosophy', '#9333EA'),
('Medieval', 'Auto-generated tag for Medieval', '#84CC16'),
('Byzantine', 'Auto-generated tag for Byzantine', '#F59E0B'),
('Conquest', 'Auto-generated tag for Conquest', '#EF4444'),
-- Nations and Civilizations
('Romans', 'Related to the Roman Empire and Roman civilization', '#DC143C'),
('Greeks', 'Related to ancient Greek civilization and culture', '#4169E1'),
('Persians', 'Related to the Persian Empire and Persian civilization', '#FFD700'),
('Byzantines', 'Related to the Byzantine Empire', '#800080'),
('Ottomans', 'Related to the Ottoman Empire', '#008000'),
('Spanish', 'Related to Spanish exploration and empire', '#FF4500'),
-- Historical Personalities
('Alexander the Great', 'The famous Macedonian conqueror', '#FF6B35'),
('Julius Caesar', 'The famous Roman general and dictator', '#8B0000'),
('Columbus', 'Christopher Columbus, the explorer', '#1E90FF'),
('Romulus', 'Legendary founder of Rome', '#B22222'),
('Plato', 'The great Greek philosopher', '#4682B4'),
-- Places and Concepts
('Rome', 'The eternal city and center of Roman power', '#8B4513'),
('Athens', 'The birthplace of democracy and philosophy', '#6A5ACD'),
('Senate', 'Related to Roman Senate and political institutions', '#696969'),
('Constantinople', 'The great city and Byzantine capital', '#DAA520'),
('Exploration', 'Related to exploration and discovery', '#20B2AA'),
('Founding', 'Related to the establishment of cities or nations', '#CD853F');

-- Insert event-tag relationships
INSERT INTO event_tags (event_id, tag_id) VALUES
-- Foundation of Rome (id=1)
(1, (SELECT id FROM tags WHERE name = 'Ancient')),
(1, (SELECT id FROM tags WHERE name = 'Political')),
(1, (SELECT id FROM tags WHERE name = 'Roman')),
(1, (SELECT id FROM tags WHERE name = 'Romans')),
(1, (SELECT id FROM tags WHERE name = 'Rome')),
(1, (SELECT id FROM tags WHERE name = 'Romulus')),
(1, (SELECT id FROM tags WHERE name = 'Founding')),

-- Battle of Marathon (id=2)
(2, (SELECT id FROM tags WHERE name = 'Ancient')),
(2, (SELECT id FROM tags WHERE name = 'Greek')),
(2, (SELECT id FROM tags WHERE name = 'Military')),
(2, (SELECT id FROM tags WHERE name = 'Greeks')),
(2, (SELECT id FROM tags WHERE name = 'Persians')),

-- Birth of Plato (id=3)
(3, (SELECT id FROM tags WHERE name = 'Cultural')),
(3, (SELECT id FROM tags WHERE name = 'Greek')),
(3, (SELECT id FROM tags WHERE name = 'Philosophy')),
(3, (SELECT id FROM tags WHERE name = 'Greeks')),
(3, (SELECT id FROM tags WHERE name = 'Athens')),
(3, (SELECT id FROM tags WHERE name = 'Plato')),

-- Battle of Gaugamela (id=4)
(4, (SELECT id FROM tags WHERE name = 'Ancient')),
(4, (SELECT id FROM tags WHERE name = 'Conquest')),
(4, (SELECT id FROM tags WHERE name = 'Military')),
(4, (SELECT id FROM tags WHERE name = 'Greeks')),
(4, (SELECT id FROM tags WHERE name = 'Persians')),
(4, (SELECT id FROM tags WHERE name = 'Alexander the Great')),

-- Assassination of Julius Caesar (id=5)
(5, (SELECT id FROM tags WHERE name = 'Ancient')),
(5, (SELECT id FROM tags WHERE name = 'Political')),
(5, (SELECT id FROM tags WHERE name = 'Roman')),
(5, (SELECT id FROM tags WHERE name = 'Romans')),
(5, (SELECT id FROM tags WHERE name = 'Julius Caesar')),
(5, (SELECT id FROM tags WHERE name = 'Senate')),

-- Fall of Constantinople (id=6)
(6, (SELECT id FROM tags WHERE name = 'Byzantine')),
(6, (SELECT id FROM tags WHERE name = 'Medieval')),
(6, (SELECT id FROM tags WHERE name = 'Military')),
(6, (SELECT id FROM tags WHERE name = 'Byzantines')),
(6, (SELECT id FROM tags WHERE name = 'Ottomans')),
(6, (SELECT id FROM tags WHERE name = 'Constantinople')),

-- Discovery of America (id=7)
(7, (SELECT id FROM tags WHERE name = 'Cultural')),
(7, (SELECT id FROM tags WHERE name = 'Renaissance')),
(7, (SELECT id FROM tags WHERE name = 'Spanish')),
(7, (SELECT id FROM tags WHERE name = 'Columbus')),
(7, (SELECT id FROM tags WHERE name = 'Exploration'));

-- Insert date template groups
INSERT INTO date_template_groups (name, description, display_order) VALUES 
-- Ancient Civilizations (3500 BC - 1000 BC)
('Dawn of Civilization', 'First cities, writing, and the birth of complex society', 0),
('Early Dynastic Period', 'Egyptian dynasties and Sumerian city-states', 1),
('First Empires', 'Akkadian Empire and Egyptian Old Kingdom', 2),
('Bronze Age Powers', 'Babylonian Empire and Egyptian Middle Kingdom', 3),
('Imperial Age', 'Egyptian New Kingdom and Hittite Empire', 4),
('Bronze Age Collapse', 'Sea Peoples and civilizational collapse', 5),
-- Classical and Later Periods (1000 BC onwards)
('Ancient Greece', 'Classical period and major Greek events', 6),
('Roman Empire', 'Major periods and events of the Roman Empire', 7),
('Medieval Period', 'Key events and periods of the Middle Ages', 8),
('Renaissance', 'European Renaissance period and events', 9),
('Modern Era', 'Major events from 1500 onwards', 10);

-- Insert date templates for ancient civilizations

-- Dawn of Civilization (Group ID: 1)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(1, 'First Cities', 'Earliest urban settlements in Mesopotamia', '3500-01-01', 'BC', '3200-01-01', 'BC', 1),
(1, 'Invention of Writing', 'Development of cuneiform and hieroglyphs', '3300-01-01', 'BC', '3000-01-01', 'BC', 2),
(1, 'Egyptian Unification', 'Narmer unites Upper and Lower Egypt', '3100-01-01', 'BC', '3000-01-01', 'BC', 3);

-- Early Dynastic Period (Group ID: 2)  
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(2, 'Egyptian Early Dynasties', 'First through Second Dynasties', '3000-01-01', 'BC', '2686-01-01', 'BC', 1),
(2, 'Sumerian City-States', 'Independent Mesopotamian cities', '2900-01-01', 'BC', '2334-01-01', 'BC', 2),
(2, 'Age of Pyramids', 'Egyptian Third and Fourth Dynasties', '2686-01-01', 'BC', '2500-01-01', 'BC', 3),
(2, 'Great Pyramid Era', 'Construction of Giza pyramids', '2580-01-01', 'BC', '2510-01-01', 'BC', 4);

-- First Empires (Group ID: 3)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(3, 'Akkadian Empire', 'Sargon conquers Mesopotamia', '2334-01-01', 'BC', '2154-01-01', 'BC', 1),
(3, 'Egyptian Old Kingdom', 'Height of pyramid building', '2686-01-01', 'BC', '2181-01-01', 'BC', 2),
(3, 'Third Dynasty of Ur', 'Neo-Sumerian revival', '2112-01-01', 'BC', '2004-01-01', 'BC', 3),
(3, 'Climate Crisis', 'Drought and empire collapse', '2200-01-01', 'BC', '2100-01-01', 'BC', 4);

-- Bronze Age Powers (Group ID: 4)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(4, 'Old Babylonian Period', 'Rise of Babylon under Hammurabi', '2000-01-01', 'BC', '1600-01-01', 'BC', 1),
(4, 'Egyptian Middle Kingdom', 'Egypt reunified and prosperous', '2055-01-01', 'BC', '1650-01-01', 'BC', 2),
(4, 'Code of Hammurabi', 'First comprehensive law code', '1792-01-01', 'BC', '1750-01-01', 'BC', 3),
(4, 'Hyksos Period', 'Foreign rulers control northern Egypt', '1720-01-01', 'BC', '1550-01-01', 'BC', 4);

-- Imperial Age (Group ID: 5)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(5, 'Egyptian New Kingdom', 'Egypt as imperial superpower', '1550-01-01', 'BC', '1077-01-01', 'BC', 1),
(5, 'Hittite Empire', 'Anatolian power challenges Egypt', '1650-01-01', 'BC', '1180-01-01', 'BC', 2),
(5, 'Reign of Hatshepsut', 'Female pharaoh and prosperity', '1479-01-01', 'BC', '1458-01-01', 'BC', 3),
(5, 'Thutmose III Campaigns', 'Egyptian military expansion', '1458-01-01', 'BC', '1425-01-01', 'BC', 4),
(5, 'Ramesses II Era', 'The Great Builder pharaoh', '1279-01-01', 'BC', '1213-01-01', 'BC', 5),
(5, 'Battle of Kadesh', 'Largest chariot battle in history', '1274-01-01', 'BC', '1274-12-31', 'BC', 6);

-- Bronze Age Collapse (Group ID: 6)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(6, 'Sea Peoples Invasions', 'Mysterious raiders devastate civilizations', '1230-01-01', 'BC', '1150-01-01', 'BC', 1),
(6, 'Fall of Hittite Empire', 'End of Anatolian superpower', '1200-01-01', 'BC', '1180-01-01', 'BC', 2),
(6, 'Bronze Age Collapse', 'Civilizational crisis across Eastern Mediterranean', '1200-01-01', 'BC', '1150-01-01', 'BC', 3),
(6, 'End of New Kingdom', 'Egypt loses imperial power', '1200-01-01', 'BC', '1077-01-01', 'BC', 4),
(6, 'Rise of Iron Age', 'Transition to new technologies and powers', '1150-01-01', 'BC', '1000-01-01', 'BC', 5);

-- Insert sample date templates for Ancient Greece (Group ID: 7)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(7, 'Archaic Period', 'Early Greek civilization', '0800-01-01', 'BC', '0480-01-01', 'BC', 1),
(7, 'Classical Period', 'Golden age of Athens', '0480-01-01', 'BC', '0323-06-10', 'BC', 2),
(7, 'Persian Wars', 'Greeks vs Persian Empire', '0499-01-01', 'BC', '0449-01-01', 'BC', 3),
(7, 'Peloponnesian War', 'Athens vs Sparta', '0431-01-01', 'BC', '0404-01-01', 'BC', 4),
(7, 'Age of Alexander', 'Alexander the Great conquests', '0336-01-01', 'BC', '0323-06-10', 'BC', 5),
(7, 'Hellenistic Period', 'Greek culture spread', '0323-06-10', 'BC', '0031-09-02', 'BC', 6);

-- Insert sample date templates for Roman Empire (Group ID: 8)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(8, 'Roman Republic', 'From founding to Julius Caesar', '0509-01-01', 'BC', '0049-01-01', 'BC', 1),
(8, 'Rule of Julius Caesar', 'Dictatorship of Julius Caesar', '0049-01-01', 'BC', '0044-03-15', 'BC', 2),
(8, 'Rule of Augustus', 'First Roman Emperor Octavian Augustus', '0031-01-01', 'BC', '0014-08-19', 'AD', 3),
(8, 'Julio-Claudian Dynasty', 'Augustus to Nero', '0027-01-01', 'BC', '0068-06-09', 'AD', 4),
(8, 'Flavian Dynasty', 'Vespasian, Titus, and Domitian', '0069-01-01', 'AD', '0096-09-18', 'AD', 5),
(8, 'Antonine Dynasty', 'Height of Roman Empire', '0096-01-01', 'AD', '0192-12-31', 'AD', 6),
(8, 'Crisis of Third Century', 'Period of military anarchy', '0235-01-01', 'AD', '0284-01-01', 'AD', 7),
(8, 'Late Roman Empire', 'From Diocletian to fall', '0284-01-01', 'AD', '0476-09-04', 'AD', 8);

-- Insert sample date templates for Medieval Period (Group ID: 9)  
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(9, 'Early Middle Ages', 'Fall of Rome to Charlemagne', '0476-09-04', 'AD', '0800-12-25', 'AD', 1),
(9, 'Carolingian Empire', 'Charlemagne and successors', '0800-12-25', 'AD', '0987-01-01', 'AD', 2),
(9, 'High Middle Ages', 'Crusades and medieval prosperity', '1000-01-01', 'AD', '1300-01-01', 'AD', 3),
(9, 'First Crusade', 'Recapture of Jerusalem', '1095-11-27', 'AD', '1099-07-15', 'AD', 4),
(9, 'Third Crusade', 'Richard Lionheart vs Saladin', '1189-01-01', 'AD', '1192-09-02', 'AD', 5),
(9, 'Late Middle Ages', 'Black Death to Renaissance', '1300-01-01', 'AD', '1453-05-29', 'AD', 6),
(9, 'Hundred Years War', 'England vs France', '1337-05-01', 'AD', '1453-10-19', 'AD', 7),
(9, 'Black Death', 'Bubonic plague pandemic', '1347-01-01', 'AD', '1351-12-31', 'AD', 8);

-- Insert sample date templates for Renaissance (Group ID: 10)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(10, 'Italian Renaissance', 'Birth of Renaissance in Italy', '1350-01-01', 'AD', '1550-01-01', 'AD', 1),
(10, 'Northern Renaissance', 'Renaissance in Northern Europe', '1450-01-01', 'AD', '1600-01-01', 'AD', 2),
(10, 'Age of Exploration', 'Discovery of New World', '1415-01-01', 'AD', '1600-01-01', 'AD', 3),
(10, 'Protestant Reformation', 'Luther to Peace of Westphalia', '1517-10-31', 'AD', '1648-10-24', 'AD', 4);

-- Insert sample date templates for Modern Era (Group ID: 11)
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(11, 'Age of Enlightenment', 'Reason and scientific revolution', '1650-01-01', 'AD', '1800-01-01', 'AD', 1),
(11, 'French Revolution', 'Revolution and Napoleonic Wars', '1789-07-14', 'AD', '1815-06-18', 'AD', 2),
(11, 'Industrial Revolution', 'Steam age and modernization', '1760-01-01', 'AD', '1840-01-01', 'AD', 3),
(11, 'World War I', 'The Great War', '1914-07-28', 'AD', '1918-11-11', 'AD', 4),
(11, 'World War II', 'Global conflict', '1939-09-01', 'AD', '1945-09-02', 'AD', 5);

-- Create view for easy template querying with display dates
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

-- Create view for events with display dates and tags
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
DROP VIEW IF EXISTS events_with_display_dates;
DROP VIEW IF EXISTS date_templates_with_display;
DROP TABLE IF EXISTS event_tags;
DROP TABLE IF EXISTS date_templates;
DROP TABLE IF EXISTS date_template_groups;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS event_datasets;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS users;