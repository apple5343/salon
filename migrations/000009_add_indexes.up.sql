CREATE INDEX idx_cars_model_id ON cars(model_id);
CREATE INDEX idx_cars_supplier_id ON cars(supplier_id);
CREATE INDEX idx_cars_status ON cars(status);
CREATE INDEX idx_cars_price ON cars(price);
CREATE INDEX idx_cars_year ON cars(year);
CREATE INDEX idx_cars_created_at ON cars(created_at);

CREATE INDEX idx_cars_status_model_id ON cars(status, model_id);

CREATE INDEX idx_sales_car_id ON sales(car_id);
CREATE INDEX idx_sales_client_id ON sales(client_id);
CREATE INDEX idx_sales_employee_id ON sales(employee_id);
CREATE INDEX idx_sales_status ON sales(status);
CREATE INDEX idx_sales_sale_date ON sales(sale_date);
CREATE INDEX idx_sales_created_at ON sales(created_at);

CREATE INDEX idx_sales_employee_date ON sales(employee_id, sale_date);

CREATE INDEX idx_models_brand_id ON models(brand_id);

CREATE INDEX idx_events_entity_type_id ON events(entity_type, entity_id);
CREATE INDEX idx_events_actor_id ON events(actor_id);
CREATE INDEX idx_events_created_at ON events(created_at DESC);
CREATE INDEX idx_events_event_type ON events(event_type);

CREATE INDEX idx_employees_status ON employees(status);
CREATE INDEX idx_employees_role ON employees(role);
CREATE INDEX idx_employees_hire_date ON employees(hire_date);