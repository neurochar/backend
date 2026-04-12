-- +goose Up
ALTER TABLE "testing_room"
    ADD COLUMN "finished_at" TIMESTAMPTZ,
    ADD COLUMN "result_index" INT,
    ADD COLUMN "finished_ip" INET;


-- +goose Down
ALTER TABLE "testing_room"
    DROP COLUMN IF EXISTS "finished_at",
    DROP COLUMN IF EXISTS "result_index",
    DROP COLUMN IF EXISTS "finished_ip";