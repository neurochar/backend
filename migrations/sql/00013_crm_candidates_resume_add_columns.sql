-- +goose Up
ALTER TABLE "crm_candidates_resume"
    ADD COLUMN "file_type" SMALLINT NOT NULL;


-- +goose Down
ALTER TABLE "crm_candidates_resume"
    DROP COLUMN IF EXISTS "file_type";