-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tokens" (
     id text NOT NULL,
     user_id text NOT NULL,
     auth_id text NOT NULL,
     hash text NOT NULL,
     last_used_at timestamp,
     expires_at timestamp NOT NULL,
     created_at timestamp NOT NULL,
     updated_at timestamp NOT NULL,
     deleted timestamp,
     PRIMARY KEY(id),
     FOREIGN KEY(user_id) REFERENCES "users" (id) ON DELETE CASCADE ON UPDATE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tokens";
-- +goose StatementEnd
