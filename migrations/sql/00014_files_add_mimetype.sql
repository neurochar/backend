-- +goose Up
ALTER TABLE "file"
    ADD COLUMN "file_mimetype" TEXT,
    ADD COLUMN "file_hash" TEXT;


-- +goose Down
ALTER TABLE "file"
    DROP COLUMN IF EXISTS "file_mimetype",
    DROP COLUMN IF EXISTS "file_hash";