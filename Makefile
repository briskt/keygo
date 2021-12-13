app:
	docker-compose up -d app

migrate:
	docker-compose run --rm app migrate

new-migration:
	bash -c "cd application/migrations && goose create change_me sql"