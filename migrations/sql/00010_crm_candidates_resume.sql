-- +goose Up

-- Таблица crm_candidates_resume
CREATE TABLE "crm_candidates_resume" (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES "tenant"(id) ON DELETE RESTRICT,
    "status"                        SMALLINT NOT NULL,
    candidate_id                    UUID REFERENCES "crm_candidate"(id) ON DELETE RESTRICT,
    file_id                         UUID NOT NULL REFERENCES "file"(id) ON DELETE RESTRICT,
    file_hash                       TEXT NOT NULL,
    analyze_data                    JSONB,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);
CREATE INDEX idx_crm_candidates_resume_tenant_id ON "crm_candidates_resume" (tenant_id);
CREATE INDEX idx_crm_candidates_resume_file_hash ON "crm_candidates_resume" (file_hash);
CREATE INDEX idx_crm_candidates_resume_status ON "crm_candidates_resume" ("status");

-- +goose Down

DROP INDEX IF EXISTS idx_crm_candidates_resume_status;
DROP INDEX IF EXISTS idx_crm_candidates_resume_file_hash;
DROP INDEX IF EXISTS idx_crm_candidates_resume_tenant_id;
DROP TABLE IF EXISTS "crm_candidates_resume";
