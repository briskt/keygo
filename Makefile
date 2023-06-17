app: db migrate adminer ui-app
	docker-compose up -d app

ui-app:
	docker-compose up -d ui-app

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

install-js-deps:
	docker-compose run --rm ui-app npm install

proxy:
	docker-compose up -d proxy

.PHONY: app ui-app test bounce migrate migratedown new-migration db fresh adminer install-js-deps proxy
