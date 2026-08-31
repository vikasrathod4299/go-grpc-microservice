CREATE TABLE IF NOT EXISTS drivers (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    phone VARCHAR(32) UNIQUE NOT NULL,
    vehicle_make VARCHAR(64) NOT NULL,
    vehicle_model VARCHAR(64) NOT NULL,
    license_plate VARCHAR(32) UNIQUE NOT NULL,
    availability VARCHAR(32) DEFAULT 'OFFLINE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_drivers_availability ON drivers(availability);
