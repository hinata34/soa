-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    login text NOT NULL UNIQUE,
    password text NOT NULL,
    name text NOT NULL DEFAULT '',
    surname text NOT NULL DEFAULT '',
    birthday timestamp NOT NULL DEFAULT NOW(),
    email text NOT NULL UNIQUE,
    mobile_number text NOT NULL DEFAULT '',
    created timestamp NOT NULL DEFAULT NOW(),
    updated timestamp NOT NULL DEFAULT NOW(),
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
