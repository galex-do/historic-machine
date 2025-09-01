-- Complete database schema for Historical Events Mapping

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

-- Create indexes
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_lens_type ON events(lens_type);
CREATE INDEX idx_events_era ON events(era);
CREATE INDEX idx_events_location ON events(latitude, longitude);
CREATE INDEX idx_date_templates_group_id ON date_templates(group_id);
CREATE INDEX idx_date_templates_dates ON date_templates(start_date, end_date);
CREATE INDEX idx_date_template_groups_order ON date_template_groups(display_order);
CREATE INDEX idx_date_templates_order ON date_templates(display_order);

-- Insert sample historical events
INSERT INTO events (name, description, latitude, longitude, event_date, era, lens_type) VALUES 
('Foundation of Rome', 'Legendary founding of Rome by Romulus', 41.9028, 12.4964, '0753-04-21', 'BC', 'historic'),
('Battle of Marathon', 'Greeks defeat Persian invasion', 38.1462, 23.9608, '0490-09-12', 'BC', 'military'),
('Birth of Plato', 'Great philosopher born in Athens', 37.9755, 23.7348, '0428-05-21', 'BC', 'cultural'),
('Battle of Gaugamela', 'Alexander defeats Darius III', 36.3402, 43.1503, '0331-10-01', 'BC', 'military'),
('Assassination of Julius Caesar', 'Caesar killed in the Roman Senate', 41.8955, 12.4823, '0044-03-15', 'BC', 'political'),
('Fall of Constantinople', 'End of the Byzantine Empire', 41.0082, 28.9784, '1453-05-29', 'AD', 'military'),
('Discovery of America', 'Columbus reaches the New World', 25.0343, -77.3963, '1492-10-12', 'AD', 'historic');

-- Insert date template groups
INSERT INTO date_template_groups (name, description, display_order) VALUES 
('Roman Empire', 'Major periods and events of the Roman Empire', 1),
('Medieval Period', 'Key events and periods of the Middle Ages', 2),
('Ancient Greece', 'Classical period and major Greek events', 3),
('Renaissance', 'European Renaissance period and events', 4),
('Modern Era', 'Major events from 1500 onwards', 5);

-- Insert sample date templates for Roman Empire
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(1, 'Roman Republic', 'From founding to Julius Caesar', '0509-01-01', 'BC', '0049-01-01', 'BC', 1),
(1, 'Rule of Julius Caesar', 'Dictatorship of Julius Caesar', '0049-01-01', 'BC', '0044-03-15', 'BC', 2),
(1, 'Rule of Augustus', 'First Roman Emperor Octavian Augustus', '0031-01-01', 'BC', '0014-08-19', 'AD', 3),
(1, 'Julio-Claudian Dynasty', 'Augustus to Nero', '0027-01-01', 'BC', '0068-06-09', 'AD', 4),
(1, 'Flavian Dynasty', 'Vespasian, Titus, and Domitian', '0069-01-01', 'AD', '0096-09-18', 'AD', 5),
(1, 'Antonine Dynasty', 'Height of Roman Empire', '0096-01-01', 'AD', '0192-12-31', 'AD', 6),
(1, 'Crisis of Third Century', 'Period of military anarchy', '0235-01-01', 'AD', '0284-01-01', 'AD', 7),
(1, 'Late Roman Empire', 'From Diocletian to fall', '0284-01-01', 'AD', '0476-09-04', 'AD', 8);

-- Insert sample date templates for Medieval Period  
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(2, 'Early Middle Ages', 'Fall of Rome to Charlemagne', '0476-09-04', 'AD', '0800-12-25', 'AD', 1),
(2, 'Carolingian Empire', 'Charlemagne and successors', '0800-12-25', 'AD', '0987-01-01', 'AD', 2),
(2, 'High Middle Ages', 'Crusades and medieval prosperity', '1000-01-01', 'AD', '1300-01-01', 'AD', 3),
(2, 'First Crusade', 'Recapture of Jerusalem', '1095-11-27', 'AD', '1099-07-15', 'AD', 4),
(2, 'Third Crusade', 'Richard Lionheart vs Saladin', '1189-01-01', 'AD', '1192-09-02', 'AD', 5),
(2, 'Late Middle Ages', 'Black Death to Renaissance', '1300-01-01', 'AD', '1453-05-29', 'AD', 6),
(2, 'Hundred Years War', 'England vs France', '1337-05-01', 'AD', '1453-10-19', 'AD', 7),
(2, 'Black Death', 'Bubonic plague pandemic', '1347-01-01', 'AD', '1351-12-31', 'AD', 8);

-- Insert sample date templates for Ancient Greece
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(3, 'Archaic Period', 'Early Greek civilization', '0800-01-01', 'BC', '0480-01-01', 'BC', 1),
(3, 'Classical Period', 'Golden age of Athens', '0480-01-01', 'BC', '0323-06-10', 'BC', 2),
(3, 'Persian Wars', 'Greeks vs Persian Empire', '0499-01-01', 'BC', '0449-01-01', 'BC', 3),
(3, 'Peloponnesian War', 'Athens vs Sparta', '0431-01-01', 'BC', '0404-01-01', 'BC', 4),
(3, 'Age of Alexander', 'Alexander the Great conquests', '0336-01-01', 'BC', '0323-06-10', 'BC', 5),
(3, 'Hellenistic Period', 'Greek culture spread', '0323-06-10', 'BC', '0031-09-02', 'BC', 6);

-- Insert sample date templates for Renaissance
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(4, 'Italian Renaissance', 'Birth of Renaissance in Italy', '1350-01-01', 'AD', '1550-01-01', 'AD', 1),
(4, 'Northern Renaissance', 'Renaissance in Northern Europe', '1450-01-01', 'AD', '1600-01-01', 'AD', 2),
(4, 'Age of Exploration', 'Discovery of New World', '1415-01-01', 'AD', '1600-01-01', 'AD', 3),
(4, 'Protestant Reformation', 'Luther to Peace of Westphalia', '1517-10-31', 'AD', '1648-10-24', 'AD', 4);

-- Insert sample date templates for Modern Era
INSERT INTO date_templates (group_id, name, description, start_date, start_era, end_date, end_era, display_order) VALUES 
(5, 'Age of Enlightenment', 'Reason and scientific revolution', '1650-01-01', 'AD', '1800-01-01', 'AD', 1),
(5, 'French Revolution', 'Revolution and Napoleonic Wars', '1789-07-14', 'AD', '1815-06-18', 'AD', 2),
(5, 'Industrial Revolution', 'Steam age and modernization', '1760-01-01', 'AD', '1840-01-01', 'AD', 3),
(5, 'World War I', 'The Great War', '1914-07-28', 'AD', '1918-11-11', 'AD', 4),
(5, 'World War II', 'Global conflict', '1939-09-01', 'AD', '1945-09-02', 'AD', 5);

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

-- Create view for events with display dates  
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