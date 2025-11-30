-- +goose Up

-- Таблица testing_room
CREATE TABLE "testing_room" (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES "tenant"(id) ON DELETE RESTRICT,
    "status"                        SMALLINT NOT NULL,
    candidate_id                    UUID REFERENCES "crm_candidate"(id) ON DELETE SET NULL,
    profile_id                      UUID REFERENCES "testing_profile"(id) ON DELETE SET NULL,
    personality_traits_map          JSONB NOT NULL,
    technique_data                  JSONB NOT NULL,
    raw_answer                      JSONB,
    candidate_answer_data           JSONB,
    result                          JSONB,
    created_by                      UUID NULL REFERENCES "tenant_account"(id) ON DELETE SET NULL,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);
CREATE INDEX idx_testing_room_tenant_id ON "testing_room" (tenant_id);


-- +goose Down

DROP INDEX IF EXISTS idx_testing_room_tenant_id;
DROP TABLE IF EXISTS "testing_room";
