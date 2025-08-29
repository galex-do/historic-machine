-- Database initialization script
-- This ensures the database exists and has proper permissions

-- Create historical_events database if it doesn't exist
-- Note: This is primarily for documentation as PostgreSQL Docker image creates the DB from env vars

-- Grant permissions to postgres user
GRANT ALL PRIVILEGES ON DATABASE historical_events TO postgres;