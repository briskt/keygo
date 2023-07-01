-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    id text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    role text NOT NULL,
    email text NOT NULL UNIQUE,
    avatar_url text NOT NULL,
    last_login_at timestamp,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
