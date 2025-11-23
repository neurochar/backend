-- +goose Up

ALTER TABLE tenant_account_code
ALTER COLUMN request_ip DROP NOT NULL;

-- +goose Down

ALTER TABLE tenant_account_code
ALTER COLUMN request_ip SET NOT NULL;