DROP INDEX idx_car_models_brand       IF EXISTS;
DROP INDEX idx_car_models_fuel_type   IF EXISTS;
DROP INDEX idx_car_models_body_type   IF EXISTS;
DROP INDEX idx_car_models_class       IF EXISTS;
DROP INDEX idx_car_models_seats       IF EXISTS;

DROP TABLE IF EXISTS car_models;

DROP TYPE IF EXISTS car_class;
DROP TYPE IF EXISTS car_body_type;
DROP TYPE IF EXISTS car_transmission;
DROP TYPE IF EXISTS car_fuel_type;
