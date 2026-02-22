CREATE TABLE IF NOT EXISTS cars (
    id uuid PRIMARY KEY NOT NULL,
    model_id uuid NOT NULL REFERENCES models(id),
    supplier_id uuid NOT NULL REFERENCES suppliers(id),
    vin VARCHAR(255) NOT NULL UNIQUE CHECK (length(vin) = 17),
    year INTEGER NOT NULL,
    color VARCHAR(64) NOT NULL,
    interior_color VARCHAR(64) NOT NULL,
    mileage BIGINT NOT NULL CHECK (mileage >= 0),
    price NUMERIC(14,2) NOT NULL CHECK (price >= 0),
    status VARCHAR(32) NOT NULL,
    options JSONB,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)