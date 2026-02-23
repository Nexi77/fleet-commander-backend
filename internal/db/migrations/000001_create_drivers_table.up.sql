-- Enable PostGIS extension (required for spatial data)
CREATE EXTENSION IF NOT EXISTS postgis;

-- Define a native PostgreSQL ENUM for the driver's status
CREATE TYPE driver_status AS ENUM ('AVAILABLE', 'ON_TRIP', 'OFFLINE');

-- Create the drivers table
CREATE TABLE drivers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status driver_status NOT NULL DEFAULT 'OFFLINE',
    
    -- GEOGRAPHY type is optimized for GPS coordinates (longitude/latitude on Earth's sphere).
    -- 4326 is the SRID for WGS84 (standard GPS coordinate system).
    current_location GEOGRAPHY(POINT, 4326),
    
    -- Always use TIMESTAMPTZ (with time zone) in distributed systems
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create a spatial index (GiST) for blazingly fast proximity queries 
-- (e.g., "Find drivers within 5km radius")
CREATE INDEX idx_drivers_location ON drivers USING GIST (current_location);