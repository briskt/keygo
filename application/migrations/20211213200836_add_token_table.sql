-- +goose Up
-- +goose StatementBegin

CREATE TABLE "auths" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    avatar_url text NOT NULL,
    provider text NOT NULL,
    provider_id text NOT NULL,
    access_token text NOT NULL,
    refresh_token text NOT NULL,
    expiry timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES "users" (id) ON DELETE CASCADE ON UPDATE RESTRICT,
    UNIQUE(user_id, provider),
    UNIQUE(provider, provider_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "auths";
-- +goose StatementEnd
