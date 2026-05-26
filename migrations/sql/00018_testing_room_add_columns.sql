-- +goose Up
ALTER TABLE "testing_room"
    ADD COLUMN "started_at" TIMESTAMPTZ,
    ADD COLUMN "traits_statuses" JSONB;


-- +goose Down
ALTER TABLE "testing_room"
    DROP COLUMN IF EXISTS "started_at",
    DROP COLUMN IF EXISTS "traits_statuses";