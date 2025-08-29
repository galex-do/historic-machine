-- +goose Up
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    event_date TIMESTAMP NOT NULL,
    lens_type VARCHAR(100) NOT NULL DEFAULT 'historic',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for geo-queries
CREATE INDEX idx_events_coordinates ON events(latitude, longitude);
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_lens_type ON events(lens_type);

-- Insert sample data
INSERT INTO events (name, description, latitude, longitude, event_date, lens_type) VALUES 
('Fall of Constantinople', 'Ottoman Empire conquered Byzantine Empire, marking the end of the Byzantine Empire', 41.0082, 28.9784, '1453-05-29 00:00:00', 'historic'),
('Discovery of America', 'Christopher Columbus reached the Americas', 25.7617, -80.1918, '1492-10-12 00:00:00', 'historic'),
('French Revolution Begins', 'Storming of the Bastille marked the beginning of the French Revolution', 48.8534, 2.3488, '1789-07-14 00:00:00', 'historic');

-- +goose Down
DROP TABLE IF EXISTS events;