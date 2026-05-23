CREATE TYPE maintenance_record_status AS ENUM ('pending', 'in_progress', 'completed');

CREATE TABLE car_maintenance_templates (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(200) NOT NULL,
    km_interval  INT,
    day_interval INT,
    is_mandatory BOOLEAN      NOT NULL DEFAULT false,
    warn_pct     NUMERIC(5,2) NOT NULL DEFAULT 80,
    pull_pct     NUMERIC(5,2) NOT NULL DEFAULT 100,
    created_at   TIMESTAMPTZ  NOT NULL,
    updated_at   TIMESTAMPTZ  NOT NULL
);

CREATE TABLE car_maintenance_records (
    id                 UUID                      PRIMARY KEY DEFAULT gen_random_uuid(),
    car_id             UUID                      NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    template_id        UUID                      NOT NULL REFERENCES car_maintenance_templates (id),
    status             maintenance_record_status NOT NULL DEFAULT 'pending',
    odometer_at        INT                       NOT NULL,
    due_by             TIMESTAMPTZ,
    completed_km       INT,
    cost_tenge         INT,
    assigned_to        VARCHAR(100),
    completed_at       TIMESTAMPTZ,
    notes              TEXT,
    receipt_image_keys TEXT[]                    NOT NULL DEFAULT '{}',
    created_at         TIMESTAMPTZ               NOT NULL,
    updated_at         TIMESTAMPTZ               NOT NULL
);

CREATE TABLE car_service_states (
    car_id        UUID NOT NULL REFERENCES cars (id) ON DELETE CASCADE,
    template_id   UUID NOT NULL REFERENCES car_maintenance_templates (id) ON DELETE CASCADE,
    last_km       INT  NOT NULL DEFAULT 0,
    last_date     TIMESTAMPTZ,
    next_due_km   INT,
    next_due_date TIMESTAMPTZ,
    PRIMARY KEY (car_id, template_id)
);

CREATE INDEX idx_car_maintenance_records_car_id      ON car_maintenance_records (car_id);
CREATE INDEX idx_car_maintenance_records_template_id ON car_maintenance_records (template_id);
CREATE INDEX idx_car_maintenance_records_status      ON car_maintenance_records (status);
