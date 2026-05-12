CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(50) PRIMARY KEY
);

INSERT INTO roles (name) VALUES
    ('user'),
    ('admin'),
    ('fleet_manager'),
    ('user_manager'),
    ('booking_manager');
