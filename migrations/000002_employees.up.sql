CREATE TABLE IF NOT EXISTS employees (
    id uuid PRIMARY KEY NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    email VARCHAR(64) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    passport_series CHAR(4) NOT NULL,
    passport_number CHAR(6) NOT NULL,
    passport_issued_by VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL,
    hire_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT unique_employee_email UNIQUE (email),
    CONSTRAINT unique_employee_phone UNIQUE (phone),
    CONSTRAINT unique_employee_passport UNIQUE (passport_series, passport_number)
)