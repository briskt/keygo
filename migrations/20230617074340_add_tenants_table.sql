-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tenants" (
  id text NOT NULL,
  name text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  deleted_at timestamp,
  PRIMARY KEY(id)
);
CREATE TABLE "tenant_users" (
  id text NOT NULL,
  tenant_id varchar(255) NOT NULL,
  user_id varchar(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  deleted_at timestamp,
  PRIMARY KEY(id),
  CONSTRAINT FK_tenant_users_tenants FOREIGN KEY (tenant_id) REFERENCES tenants(id),
  CONSTRAINT FK_tenant_users_user FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE (tenant_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tenant_users";
DROP TABLE "tenants";
-- +goose StatementEnd
