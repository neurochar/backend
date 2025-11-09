-- +goose Up

-- Таблица emailing
CREATE TABLE "emailing" (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_data            JSONB NOT NULL,
    sent_at                 TIMESTAMPTZ NULL,
    request_ip              INET NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_emailing_sent_at ON "emailing" (sent_at);

-- +goose Down

DROP INDEX IF EXISTS idx_emailing_sent_at;
DROP TABLE IF EXISTS "emailing";
