-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Таблица файлов
CREATE TABLE "file" (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id                UUID NOT NULL,
    file_target             TEXT NOT NULL,
    assigned_to_target      BOOLEAN NOT NULL DEFAULT FALSE,
    storage_file_key        TEXT NULL,
    original_file_name      TEXT NOT NULL,
    uploaded_to_storage     BOOLEAN NOT NULL DEFAULT FALSE,
    to_delete_from_storage  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at              TIMESTAMPTZ NULL
);
CREATE INDEX idx_file_storage_file_key ON "file" (storage_file_key);
CREATE INDEX idx_file_to_delete_from_storage ON "file" (to_delete_from_storage) WHERE to_delete_from_storage = TRUE;

-- Таблица ролей
CREATE TABLE "role" (
    id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        TEXT NOT NULL,
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    is_super    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
INSERT INTO "role" (name, is_system, is_super) VALUES ('пользователь', TRUE, FALSE), ('администратор', TRUE, TRUE);


-- Таблица связи ролей и прав
CREATE TABLE "role_to_right" (
    role_id             BIGINT NOT NULL REFERENCES "role"(id) ON DELETE CASCADE,
    role_right_id       BIGINT NOT NULL,
    value               INT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, role_right_id)
);

-- Таблица аккаунтов
CREATE TABLE "account" (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id             INT NOT NULL REFERENCES "role"(id) ON DELETE RESTRICT,
    email               TEXT NOT NULL UNIQUE,
    password_hash       TEXT NOT NULL,
    is_email_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    is_blocked          BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at       TIMESTAMPTZ,
    last_request_at     TIMESTAMPTZ,
    last_request_ip     INET,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX idx_account_role_id ON "account" (role_id) WHERE deleted_at IS NULL;

INSERT INTO "account" (id, role_id, email, password_hash, is_email_verified) VALUES ('8646d9c0-daee-4e89-81bc-b00851fd926a', 2, 'edkirill@yandex.ru', '$2a$04$hBpYvS1On29wYe6yPLdBL.EoWH0WoD7JGgQaIxV4mA.TTXEsAXfdO', TRUE);

-- Таблица кодов к аккаунтам
CREATE TABLE "account_code" (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id          UUID NOT NULL REFERENCES "account"(id) ON DELETE RESTRICT,
    code_type           SMALLINT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT FALSE,
    code                TEXT NOT NULL,
    request_ip          INET NOT NULL,
    attempts            INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_account_code_account_id ON "account_code" (account_id);

-- Таблица с сессиями панели управления
CREATE TABLE "auth_admin_session" (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id          UUID NOT NULL REFERENCES "account"(id) ON DELETE RESTRICT,
    last_request_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_request_ip     INET,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX idx_auth_admin_session_account_id ON "auth_admin_session" (account_id) WHERE deleted_at IS NULL;

-- Таблица профилей
CREATE TABLE "profile" (
    id                          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    account_id                  UUID NOT NULL REFERENCES "account"(id) ON DELETE RESTRICT,
    name                        TEXT NOT NULL,
    surname                     TEXT NOT NULL,
    photo_100x100_file_id       UUID  NULL REFERENCES "file"(id) ON DELETE RESTRICT,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ
);
CREATE INDEX idx_profile_account_id ON "profile" (account_id) WHERE deleted_at IS NULL;

INSERT INTO "profile" (account_id, name, surname) VALUES ('8646d9c0-daee-4e89-81bc-b00851fd926a', 'Кирилл', 'Администратор');

-- +goose Down
DROP INDEX IF EXISTS idx_profile_account_id;
DROP TABLE IF EXISTS "profile";

DROP INDEX IF EXISTS idx_auth_admin_session_account_id;
DROP TABLE IF EXISTS "auth_admin_session";

DROP INDEX IF EXISTS idx_account_code_account_id;
DROP TABLE IF EXISTS "account_code";

DROP INDEX IF EXISTS idx_account_role_id;
DROP TABLE IF EXISTS "account";

DROP TABLE IF EXISTS "role_to_right";

DROP TABLE IF EXISTS "role";

DROP INDEX IF EXISTS idx_file_to_delete_from_storage;
DROP INDEX IF EXISTS idx_file_storage_file_key;
DROP TABLE IF EXISTS "file";
