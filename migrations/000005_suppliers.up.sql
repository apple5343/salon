CREATE TABLE IF NOT EXISTS suppliers (
    id uuid PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    country_code CHAR(2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)