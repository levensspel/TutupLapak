include .env

run:
	go run .

create_migration:
	@if [ -z "$(title)" ]; then \
		echo "error: please specify a title first"; \
		exit 1; \
	fi
	migrate create -ext=sql -dir=src/database/migrations -seq $(title)

migrate_up:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

.PHONY: run create_migration migrate_up migrate_down
