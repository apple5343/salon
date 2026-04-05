DROP INDEX IF EXISTS idx_employees_hire_date;
DROP INDEX IF EXISTS idx_employees_role;
DROP INDEX IF EXISTS idx_employees_status;

DROP INDEX IF EXISTS idx_events_event_type;
DROP INDEX IF EXISTS idx_events_created_at;
DROP INDEX IF EXISTS idx_events_actor_id;
DROP INDEX IF EXISTS idx_events_entity_type_id;

DROP INDEX IF EXISTS idx_models_brand_id;

DROP INDEX IF EXISTS idx_sales_employee_date;
DROP INDEX IF EXISTS idx_sales_created_at;
DROP INDEX IF EXISTS idx_sales_sale_date;
DROP INDEX IF EXISTS idx_sales_status;
DROP INDEX IF EXISTS idx_sales_employee_id;
DROP INDEX IF EXISTS idx_sales_client_id;
DROP INDEX IF EXISTS idx_sales_car_id;

DROP INDEX IF EXISTS idx_cars_status_model_id;
DROP INDEX IF EXISTS idx_cars_created_at;
DROP INDEX IF EXISTS idx_cars_year;
DROP INDEX IF EXISTS idx_cars_price;
DROP INDEX IF EXISTS idx_cars_status;
DROP INDEX IF EXISTS idx_cars_supplier_id;
DROP INDEX IF EXISTS idx_cars_model_id;