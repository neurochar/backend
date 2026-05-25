-- +goose Up
ALTER TABLE "crm_candidates_resume"
    ADD COLUMN "error_text" TEXT;


-- +goose Down
ALTER TABLE "crm_candidates_resume"
    DROP COLUMN IF EXISTS "error_text";