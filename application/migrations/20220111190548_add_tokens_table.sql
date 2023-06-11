-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tokens" (
     id text NOT NULL,
     auth_id text NOT NULL,
     hash text NOT NULL,
     last_login_at timestamp,
     expires_at timestamp,
     created_at timestamp NOT NULL,
     updated_at timestamp NOT NULL,
     deleted_at timestamp,
     PRIMARY KEY(id),
     FOREIGN KEY(auth_id) REFERENCES "auths" (id) ON DELETE CASCADE ON UPDATE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tokens";
-- +goose StatementEnd
