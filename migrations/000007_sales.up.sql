CREATE TABLE IF NOT EXISTS sales (
    id uuid PRIMARY KEY NOT NULL,
    car_id uuid NOT NULL REFERENCES cars(id),
    client_id uuid NOT NULL REFERENCES clients(id),
    employee_id uuid NOT NULL REFERENCES employees(id),
    sale_date DATE NOT NULL DEFAULT CURRENT_DATE,
    original_price   NUMERIC(14,2) NOT NULL,
    discount_amount  NUMERIC(14,2) DEFAULT 0 CHECK (discount_amount >= 0),
    discount_percent NUMERIC(5,2) DEFAULT 0 CHECK (discount_percent BETWEEN 0 AND 100),
    final_price      NUMERIC(14,2) GENERATED ALWAYS AS (original_price - discount_amount) STORED,
    payment_type VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
)