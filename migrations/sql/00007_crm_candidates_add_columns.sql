-- +goose Up
ALTER TABLE "crm_candidate"
    ADD COLUMN candidate_birthday DATE,
    ADD COLUMN candidate_gender SMALLINT NOT NULL DEFAULT 0;


-- +goose Down
ALTER TABLE "crm_candidate"
    DROP COLUMN IF EXISTS candidate_gender,
    DROP COLUMN IF EXISTS candidate_birthday;