-- +goose Up

-- Таблица crm_candidate
CREATE TABLE "crm_candidate" (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES "tenant"(id) ON DELETE RESTRICT,
    candidate_name                  TEXT NOT NULL,
    candidate_surname               TEXT NOT NULL,
    created_by                      UUID NULL REFERENCES "tenant_account"(id) ON DELETE SET NULL,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);
CREATE INDEX idx_crm_candidate_tenant_id ON "crm_candidate" (tenant_id);


-- +goose Down

DROP INDEX IF EXISTS idx_crm_candidate_tenant_id;
DROP TABLE IF EXISTS "crm_candidate";
