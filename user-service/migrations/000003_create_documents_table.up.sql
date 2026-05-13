CREATE TABLE IF NOT EXISTS documents (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_type VARCHAR(100) NOT NULL,
    status     VARCHAR(50)  NOT NULL DEFAULT 'pending',
    error      TEXT,
    image_key  TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_documents_user_id ON documents(user_id);
CREATE INDEX IF NOT EXISTS idx_documents_status  ON documents(status);
