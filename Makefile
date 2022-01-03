app:
	docker-compose up -d app

bounce:
	docker-compose kill app && docker-compose up -d app

migrate:
	docker-compose run --rm app migrate

migratedown:
	docker-compose run --rm app bash -c "goose -dir migrations postgres postgres://keygo:keygo@db:5432/keygo?sslmode=disable down"

new-migration:
	bash -c "cd application/migrations && goose create change_me sql"

adminer:
	docker-compose up -d adminer

.PHONY: app, bounce, migrate, migratedown, new-migration, adminer