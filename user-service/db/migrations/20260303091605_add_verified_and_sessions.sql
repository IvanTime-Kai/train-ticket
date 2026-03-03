-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN is_verified   TINYINT(1) NOT NULL DEFAULT 0 AFTER is_active,
    ADD COLUMN last_login_at DATETIME   NULL AFTER is_verified;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE user_sessions (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    user_id       CHAR(36)     NOT NULL,
    device        VARCHAR(255) NULL,
    ip_address    VARCHAR(45)  NULL,
    logged_in_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logged_out_at DATETIME     NULL,
    CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_sessions;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN is_verified,
    DROP COLUMN last_login_at;
-- +goose StatementEnd