-- +goose Up
CREATE TABLE auths (
    auth_id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    email VARCHAR(4) NOT NULL DEFAULT '',   -- placeholder only
    email_cipher TEXT,
    email_nonce TEXT,
    email_hash VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    verify_email SMALLINT DEFAULT 0
);

CREATE INDEX idx_auths_deleted_at ON auths (deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_auths_deleted_at;
DROP TABLE IF EXISTS auths;
