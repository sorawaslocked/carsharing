CREATE TYPE insurance_type   AS ENUM ('osago', 'kasko');
CREATE TYPE insurance_status AS ENUM ('active', 'expired', 'cancelled');

CREATE TABLE car_insurances (
    id            UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id        UUID             NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    type          insurance_type   NOT NULL,
    status        insurance_status NOT NULL DEFAULT 'active',
    provider      VARCHAR(200)     NOT NULL,
    policy_num    VARCHAR(100)     NOT NULL,
    starts_at     TIMESTAMPTZ      NOT NULL,
    expires_at    TIMESTAMPTZ      NOT NULL,
    cost_tenge    INT              NOT NULL DEFAULT 0,
    image_keys    TEXT[]           NOT NULL DEFAULT '{}',
    notes         TEXT,
    created_at    TIMESTAMPTZ      NOT NULL,
    updated_at    TIMESTAMPTZ      NOT NULL
);

CREATE INDEX idx_car_insurances_car_id     ON car_insurances (car_id);
CREATE INDEX idx_car_insurances_expires_at ON car_insurances (expires_at);
