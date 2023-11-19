-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tenants" (
  id text NOT NULL,
  name text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  deleted timestamp,
  PRIMARY KEY(id),
  UNIQUE (name, deleted)
);
ALTER TABLE "users" ADD "tenant_id" text NULL;
ALTER TABLE "users" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON DELETE SET NULL ON UPDATE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP CONSTRAINT "users_tenant_id_fkey";
ALTER TABLE "users" DROP "tenant_id";
DROP TABLE "tenants";
-- +goose StatementEnd
