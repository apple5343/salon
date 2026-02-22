CREATE TABLE IF NOT EXISTS brands (
    id uuid PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    description VARCHAR(2048),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)