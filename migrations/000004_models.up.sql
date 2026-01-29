CREATE TABLE if NOT EXISTS models (
    id uuid PRIMARY KEY NOT NULL,
    brand_id uuid NOT NULL REFERENCES brands(id),
    name VARCHAR(255) NOT NULL,
    generation VARCHAR(255) NOT NULL,
    body_type VARCHAR(128) NOT NULL, 
    transmission_type VARCHAR(128) NOT NULL,
    fuel_type VARCHAR(128) NOT NULL, 
    engine_displacement INTEGER NOT NULL,
    power_hp INTEGER NOT NULL,
    drive_type VARCHAR(128) NOT NULL,
    base_price INTEGER NOT NULL,
    technical_characteristics JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
)