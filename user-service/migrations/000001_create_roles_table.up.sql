CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (id, name)
VALUES (0, 'user'),
       (1, 'admin'),
       (2, 'tech_support'),
       (3, 'finance_manager'),
       (4, 'maintenance_specialist');

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);