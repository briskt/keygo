app: db migrate adminer
	docker-compose up -d app

test:
	docker-compose run --rm test

bounce:
	docker-compose kill app && docker-compose up -d app

migrate: db
	docker-compose run --rm app bash -c "goose -dir migrations postgres postgres://keygo:keygo@db:5432/keygo?sslmode=disable up"

migratedown: db
	docker-compose run --rm app bash -c "goose -dir migrations postgres postgres://keygo:keygo@db:5432/keygo?sslmode=disable down"

new-migration: db
	docker-compose run --rm app bash -c "cd migrations && goose create change_me sql"

db:
	docker-compose up -d db

fresh:
	docker-compose kill db
	docker-compose rm -f db
	make migrate

adminer:
	docker-compose up -d adminer

.PHONY: app bounce migrate migratedown new-migration adminer
