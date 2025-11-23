-- +goose Up

-- Таблица registration
CREATE TABLE "registration" (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                   TEXT NOT NULL,
    tariff                  INT NOT NULL DEFAULT 0,
    is_finished             BOOLEAN NOT NULL,
    tenant_id               UUID NULL REFERENCES "tenant"(id) ON DELETE CASCADE,
    request_ip              INET,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX registration_email ON registration(email);


-- +goose Down

DROP INDEX IF EXISTS registration_email;
DROP TABLE IF EXISTS "registration";
