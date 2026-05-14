DROP INDEX idx_cars_model_id IF EXISTS;
DROP INDEX idx_cars_status   IF EXISTS;
DROP INDEX idx_cars_location IF EXISTS;

DROP TABLE IF EXISTS cars;

DROP TYPE IF EXISTS car_status;
