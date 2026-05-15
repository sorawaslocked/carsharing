CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id                   UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    email                citext       NOT NULL UNIQUE,
    phone_number         VARCHAR(20)  UNIQUE,
    first_name           VARCHAR(100) NOT NULL,
    last_name            VARCHAR(100) NOT NULL,
    birth_date           DATE         NOT NULL,
    password_hash        bytea        NOT NULL,
    profile_image_key    TEXT,
    is_document_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    is_email_verified    BOOLEAN      NOT NULL DEFAULT FALSE,
    is_suspended         BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_name VARCHAR(50) NOT NULL REFERENCES roles(name) ON DELETE RESTRICT,
    PRIMARY KEY (user_id, role_name)
);

CREATE INDEX IF NOT EXISTS idx_users_email             ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_email_verified ON users(is_email_verified);
CREATE INDEX IF NOT EXISTS idx_users_is_suspended      ON users(is_suspended);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id      ON user_roles(user_id);
