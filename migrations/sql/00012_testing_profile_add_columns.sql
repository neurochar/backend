-- +goose Up
ALTER TABLE "testing_profile"
    ADD COLUMN IF NOT EXISTS "description" TEXT NOT NULL DEFAULT '';


-- +goose Down
ALTER TABLE "testing_profile"
    DROP COLUMN IF EXISTS "description";