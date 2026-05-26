-- +goose Up
ALTER TABLE "testing_room"
    ADD COLUMN "is_processed" BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN "process_tries" INT NOT NULL DEFAULT 0,
    ADD COLUMN "process_error" TEXT,
    ADD COLUMN "need_process_at" TIMESTAMPTZ;


-- +goose Down
ALTER TABLE "testing_room"
    DROP COLUMN IF EXISTS "is_processed",
    DROP COLUMN IF EXISTS "process_tries",
    DROP COLUMN IF EXISTS "process_error",
    DROP COLUMN IF EXISTS "need_process_at";