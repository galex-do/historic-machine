-- +goose Up
-- Enable PostGIS extension for geospatial data
CREATE EXTENSION IF NOT EXISTS postgis;

-- Add location column as GEOGRAPHY for accurate global calculations
ALTER TABLE events ADD COLUMN location GEOGRAPHY(POINT, 4326);

-- Update existing records with location data from lat/lng
UPDATE events 
SET location = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography 
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Create spatial index for efficient geographic queries
CREATE INDEX idx_events_location ON events USING GIST (location);

-- Add constraints for data integrity
ALTER TABLE events 
ADD CONSTRAINT valid_location 
CHECK (location IS NOT NULL);

-- Create index for bounding box queries (faster for map viewport)
CREATE INDEX idx_events_geom_bbox ON events USING GIST (location::geometry);

-- +goose Down
-- Remove PostGIS enhancements
DROP INDEX IF EXISTS idx_events_location;
DROP INDEX IF EXISTS idx_events_geom_bbox;
ALTER TABLE events DROP CONSTRAINT IF EXISTS valid_location;
ALTER TABLE events DROP COLUMN IF EXISTS location;
DROP EXTENSION IF EXISTS postgis;