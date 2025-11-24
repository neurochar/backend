-- +goose Up

-- Таблица testing_profile
CREATE TABLE "testing_profile" (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES "tenant"(id) ON DELETE RESTRICT,
    "name"                          TEXT NOT NULL,
    personality_traits_map          JSONB NOT NULL,
    created_by                      UUID NULL REFERENCES "tenant_account"(id) ON DELETE SET NULL,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);
CREATE INDEX idx_testing_profile_tenant_id ON "testing_profile" (tenant_id);


-- +goose Down

DROP INDEX IF EXISTS idx_testing_profile_tenant_id;
DROP TABLE IF EXISTS "testing_profile";
