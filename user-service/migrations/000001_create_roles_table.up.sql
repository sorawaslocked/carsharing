CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (id, name)
VALUES (1, 'user'),
       (2, 'admin'),
       (3, 'tech_support'),
       (4, 'finance_manager'),
       (5, 'maintenance_specialist');

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);