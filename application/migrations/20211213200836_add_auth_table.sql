-- +goose Up
-- +goose StatementBegin

CREATE TABLE "auths" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    provider text NOT NULL,
    provider_id text NOT NULL,
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
