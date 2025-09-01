-- +goose Up
-- Add support for BC dates using astronomical year numbering
-- PostgreSQL natively supports negative years in dates

-- Add era field to distinguish BC/AD for display purposes
ALTER TABLE events ADD COLUMN era VARCHAR(2) DEFAULT 'AD';

-- Update existing events to have proper era
UPDATE events SET era = 'AD' WHERE event_date >= '0001-01-01';

-- Create index for era-based queries
CREATE INDEX idx_events_era ON events(era);

-- Insert sample BC events for testing
INSERT INTO events (name, description, latitude, longitude, event_date, era, lens_type) VALUES 
('Assassination of Julius Caesar', 'Julius Caesar was assassinated on the Ides of March', 41.8919, 12.5113, '0044-03-15 BC', 'BC', 'political'),
('Battle of Gaugamela', 'Alexander the Great defeated Darius III of Persia', 36.3500, 43.1500, '0331-10-01 BC', 'BC', 'military'),
('Foundation of Rome', 'Traditional founding date of Rome by Romulus', 41.9028, 12.4964, '0753-04-21 BC', 'BC', 'historic'),
('Birth of Plato', 'The famous Greek philosopher was born', 37.9755, 23.7348, '0428-01-01 BC', 'BC', 'cultural'),
('Battle of Marathon', 'Greeks defeated the Persian invasion', 38.1462, 23.9703, '0490-09-12 BC', 'BC', 'military');

-- Create a function to calculate years between dates across BC/AD boundary
CREATE OR REPLACE FUNCTION calculate_historical_years(start_date DATE, start_era VARCHAR, end_date DATE, end_era VARCHAR)
RETURNS INTEGER AS $$
DECLARE
    start_astronomical INTEGER;
    end_astronomical INTEGER;
    year_diff INTEGER;
BEGIN
    -- Convert to astronomical year numbering
    IF start_era = 'BC' THEN
        start_astronomical = (EXTRACT(YEAR FROM start_date) * -1) + 1;
    ELSE
        start_astronomical = EXTRACT(YEAR FROM start_date);
    END IF;
    
    IF end_era = 'BC' THEN
        end_astronomical = (EXTRACT(YEAR FROM end_date) * -1) + 1;
    ELSE
        end_astronomical = EXTRACT(YEAR FROM end_date);
    END IF;
    
    year_diff = end_astronomical - start_astronomical;
    RETURN year_diff;
END;
$$ LANGUAGE plpgsql;

-- Create a view for easy BC/AD date handling
CREATE OR REPLACE VIEW events_with_display_dates AS
SELECT 
    id,
    name,
    description,
    latitude,
    longitude,
    event_date,
    era,
    lens_type,
    created_at,
    CASE 
        WHEN era = 'BC' THEN 
            CONCAT(
                LPAD(EXTRACT(DAY FROM event_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM event_date)::TEXT, ' BC'
            )
        ELSE 
            CONCAT(
                LPAD(EXTRACT(DAY FROM event_date)::TEXT, 2, '0'), '.',
                LPAD(EXTRACT(MONTH FROM event_date)::TEXT, 2, '0'), '.',
                EXTRACT(YEAR FROM event_date)::TEXT, ' AD'
            )
    END AS display_date,
    -- Calculate astronomical year for sorting
    CASE 
        WHEN era = 'BC' THEN (EXTRACT(YEAR FROM event_date) * -1) + 1
        ELSE EXTRACT(YEAR FROM event_date)
    END AS astronomical_year
FROM events
ORDER BY astronomical_year ASC;

-- +goose Down
DROP VIEW IF EXISTS events_with_display_dates;
DROP FUNCTION IF EXISTS calculate_historical_years(DATE, VARCHAR, DATE, VARCHAR);
DROP INDEX IF EXISTS idx_events_era;
ALTER TABLE events DROP COLUMN IF EXISTS era;