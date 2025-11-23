-- +goose Up

-- Таблица tenant
CREATE TABLE "tenant" (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text_id                 TEXT NOT NULL,
    is_demo                 BOOLEAN NOT NULL DEFAULT FALSE,
    is_active               BOOLEAN NOT NULL,
    "name"                  TEXT NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at              TIMESTAMPTZ
);
CREATE UNIQUE INDEX tenant_text_id_uniq ON tenant(text_id) WHERE deleted_at IS NULL;

INSERT INTO "tenant" (id, text_id, is_demo, is_active, "name") VALUES ('00cb0816-08ae-430c-a9b2-69227b17677f', 'demo', TRUE, TRUE, 'Demo workspace');

-- Таблица аккаунтов в tenant
CREATE TABLE "tenant_account" (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES "tenant"(id) ON DELETE RESTRICT,
    role_id                         INT NOT NULL,
    email                           TEXT NOT NULL,
    password_hash                   TEXT NOT NULL,
    is_confirmed                    BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_verified               BOOLEAN NOT NULL DEFAULT FALSE,
    is_blocked                      BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at                   TIMESTAMPTZ,
    last_request_at                 TIMESTAMPTZ,
    last_request_ip                 INET,
    profile_name                    TEXT NOT NULL,
    profile_surname                 TEXT NOT NULL,
    profile_photo_100x100_file_id   UUID  NULL REFERENCES "file"(id) ON DELETE RESTRICT,
    profile_photo_original_file_id  UUID  NULL REFERENCES "file"(id) ON DELETE RESTRICT,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);
CREATE INDEX idx_tenant_account_tenant_id ON "tenant_account" (tenant_id);
CREATE UNIQUE INDEX idx_tenant_account_tenant_id_email_uniq ON tenant_account(tenant_id, email) WHERE deleted_at IS NULL;

INSERT INTO "tenant_account" (id, tenant_id, role_id, email, password_hash, is_confirmed, is_email_verified, profile_name, profile_surname) VALUES
(
    '9c73cf85-8039-4216-a8d0-cff97ec502b3',
    '00cb0816-08ae-430c-a9b2-69227b17677f',
    1,
    'edkirill@yandex.ru',
    '$2a$04$hBpYvS1On29wYe6yPLdBL.EoWH0WoD7JGgQaIxV4mA.TTXEsAXfdO',
    TRUE,
    TRUE,
    'Кирилл',
    'Администратор'
);

-- Таблица кодов к аккаунтам
CREATE TABLE "tenant_account_code" (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id          UUID NOT NULL REFERENCES "tenant_account"(id) ON DELETE CASCADE,
    code_type           SMALLINT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT FALSE,
    code                TEXT NOT NULL,
    request_ip          INET NOT NULL,
    attempts            INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tenant_account_code_account_id ON "tenant_account_code" (account_id);

-- Таблица с сессиями
CREATE TABLE "tenant_auth_session" (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id                  UUID NOT NULL REFERENCES "tenant_account"(id) ON DELETE CASCADE,
    refresh_token               UUID NOT NULL,
    refresh_version             INT NOT NULL DEFAULT 1,
    refresh_token_issued_at     TIMESTAMPTZ NOT NULL,
    refresh_token_expires_at    TIMESTAMPTZ NOT NULL,
    refresh_token_request_ip    INET,
    create_request_ip           INET,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ
);
CREATE INDEX idx_tenant_auth_session_account_id ON "tenant_auth_session" (account_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_tenant_auth_session_account_id;
DROP TABLE IF EXISTS "tenant_auth_session";

DROP INDEX IF EXISTS idx_tenant_account_code_account_id;
DROP TABLE IF EXISTS "tenant_account_code";

DROP INDEX IF EXISTS idx_tenant_account_tenant_id_email_uniq;
DROP INDEX IF EXISTS idx_tenant_account_tenant_id;
DROP TABLE IF EXISTS "tenant_account";

DROP INDEX IF EXISTS tenant_text_id_uniq;
DROP TABLE IF EXISTS "tenant";
