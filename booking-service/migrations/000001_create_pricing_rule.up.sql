CREATE TABLE IF NOT EXISTS pricing_rules (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id            UUID,
    class               TEXT,
    type                TEXT        NOT NULL,
    rate_tenge          INT         NOT NULL,
    rate_per_km_tenge   INT,
    free_minutes        INT,
    min_charge_tenge    INT,
    overtime_policy     TEXT,
    overtime_rate_tenge INT,
    is_active           BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
